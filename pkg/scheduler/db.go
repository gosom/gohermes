package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	//_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type schedulerRepository struct {
	log  zerolog.Logger
	db   *pgxpool.Pool
	pool sync.Pool
}

func newSchedulerRepository(dsn string, logger zerolog.Logger) (*schedulerRepository, error) {
	ans := schedulerRepository{
		log: logger,
		pool: sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}
	var err error
	// TODO we should close the poll => pass db conn from the caller
	ans.db, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &ans, ans.db.Ping(context.Background())
}

func (o *schedulerRepository) migrate(ctx context.Context) error {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, dbschema)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (o *schedulerRepository) insertScheduledJob(ctx context.Context, req *schedulerrpc.CreateScheduledJobRequest) (string, error) {
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var jobId string
	err = tx.QueryRow(
		ctx, insertQ, req.Endpoint, req.Data, req.ScheduledAt.AsTime(),
	).Scan(&jobId)
	if err != nil {
		return jobId, err
	}
	return jobId, tx.Commit(ctx)
}

func (o *schedulerRepository) findScheduledJobsForExecution(ctx context.Context, limit int) ([]schedulerrpc.ScheduledJob, schedulerrpc.ScheduledJob, error) {
	var next schedulerrpc.ScheduledJob
	tx, err := o.db.Begin(ctx)
	if err != nil {
		return nil, next, err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Query(ctx, findQ, limit)
	if err != nil {
		return nil, next, err
	}
	defer rows.Close()
	var items []schedulerrpc.ScheduledJob
	for rows.Next() {
		var (
			job         schedulerrpc.ScheduledJob
			scheduledAt time.Time
			createdAt   time.Time
		)
		if err := rows.Scan(&job.Id, &job.Endpoint, &job.Data,
			&scheduledAt, &createdAt, &job.Selected); err != nil {
			return nil, next, err
		}
		job.CreatedAt = timestamppb.New(createdAt)
		job.ScheduledAt = timestamppb.New(scheduledAt)
		items = append(items, job)
	}
	if err := rows.Err(); err != nil {
		return nil, next, err
	}

	var lastTime time.Time
	if len(items) > 0 {
		ids := make([]interface{}, 0, len(items))
		where := make([]string, 0, len(items))
		for i := range items {
			ids = append(ids, items[i].Id)
			where = append(where, fmt.Sprintf("$%d", i+1))
		}
		q := fmt.Sprintf(updateJobsQ, strings.Join(where, ","))
		_, err := tx.Exec(ctx, q, ids...)
		if err != nil {
			return nil, next, err
		}
		lastTime = items[len(items)-1].ScheduledAt.AsTime()
	}
	if err := o.getNext(ctx, tx, &next, lastTime); err != nil {
		return items, next, err
	}
	return items, next, tx.Commit(ctx)
}
func (o *schedulerRepository) getNext(ctx context.Context, tx pgx.Tx, job *schedulerrpc.ScheduledJob, ts time.Time) error {
	var (
		scheduledAt time.Time
		createdAt   time.Time
	)
	err := tx.QueryRow(ctx, nextJobQ, ts).Scan(
		&job.Id, &job.Endpoint, &job.Data, &scheduledAt, &createdAt, &job.Selected,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	if err != nil {
		return err
	}
	job.CreatedAt = timestamppb.New(createdAt)
	job.ScheduledAt = timestamppb.New(scheduledAt)
	return nil
}

func (o *schedulerRepository) waitNotification(ctx context.Context) (<-chan Notify, <-chan error) {
	outc := make(chan Notify)
	errc := make(chan error, 1)
	go func() {
		defer close(outc)
		defer close(errc)
		conn, err := o.db.Acquire(context.Background())
		if err != nil {
			errc <- err
			return
		}
		defer conn.Release()
		_, err = conn.Exec(context.Background(), "listen events")
		if err != nil {
			errc <- err
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := conn.Conn().WaitForNotification(ctx)
				if err != nil {
					errc <- err
					return
				}
				var job Notify
				if err := json.Unmarshal([]byte(n.Payload), &job); err != nil {
					errc <- err
					return
				}
				outc <- job
			}
		}
	}()
	return outc, errc
}

func (o *schedulerRepository) SaveExecutionResponse(ctx context.Context, id string, resp ExecutorResponse) error {
	var success bool
	if len(resp.Err) == 0 {
		success = true
	}
	buf := o.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer o.pool.Put(buf)
	if err := json.NewEncoder(buf).Encode(resp); err != nil {
		return err
	}
	_, err := o.db.Exec(ctx, insertExecResp, id, time.Now().UTC(), success, buf)
	return err
}
