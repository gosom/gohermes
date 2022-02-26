package scheduler

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
	"github.com/rs/zerolog"
)

type Notify struct {
	Table  string `json:"table"`
	Action string `json:"action"`
}

type MonitorScheduledJobs struct {
	log       zerolog.Logger
	repo      *schedulerRepository
	executors int
	buffSize  int
	resetJobs chan struct{}
	//toExecute chan schedulerrpc.ScheduledJob
}

func NewMonitorScheduledJobs(logger zerolog.Logger, repo *schedulerRepository,
	executors int) (*MonitorScheduledJobs, error) {
	ans := MonitorScheduledJobs{
		log:       logger,
		repo:      repo,
		executors: executors,
		resetJobs: make(chan struct{}),
	}
	ans.buffSize = executors * 10
	return &ans, nil
}

func (o *MonitorScheduledJobs) Run(ctx context.Context) error {
	var execErrors []<-chan error
	go o.postgresListen(ctx)
	jobsC, errC := o.scheduledJobProvider(ctx)
	execErrors = append(execErrors, errC)
	for i := 0; i < o.executors; i++ {
		errc := o.execute(ctx, i+1, jobsC)
		execErrors = append(execErrors, errc)
	}
	return waitErrors(execErrors...)
}

func (o *MonitorScheduledJobs) postgresListen(ctx context.Context) error {
	events, errs := o.repo.waitNotification(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-errs:
			return err
		case ev := <-events:
			o.log.Debug().Msgf("got notification %+v", ev)
			if ev.Table == "jobs" {
				o.resetJobs <- struct{}{}
			}
		}
	}
	return nil
}

func (o *MonitorScheduledJobs) scheduledJobProvider(ctx context.Context) (<-chan schedulerrpc.ScheduledJob, <-chan error) {
	out := make(chan schedulerrpc.ScheduledJob, o.buffSize)
	errc := make(chan error, 1)
	go func() {
		logger := o.log.With().Str("component", "provider").Logger()
		wg := sync.WaitGroup{}
		defer func() {
			logger.Warn().Msg("job provider is exiting")
		}()
		defer close(out)
		defer close(errc)
		defaultWaitDuration := 5 * time.Minute
		timer := time.NewTimer(defaultWaitDuration)
		defer timer.Stop()
		for {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			jobs, nextJob, err := o.repo.findScheduledJobsForExecution(ctx, o.buffSize)
			if err != nil {
				errc <- err // TODO probably we need to reconnect to db and continue
				return
			}
			logger.Debug().Msgf("selected %d jobs for execution", len(jobs))
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := range jobs {
					out <- jobs[i]
				}
				logger.Debug().Msgf("sent %d jobs for execution", len(jobs))
			}()
			now := time.Now().UTC()
			logger.Debug().Msgf("next is scheduled for %s", nextJob.ScheduledAt.AsTime())
			if waitTime := nextJob.ScheduledAt.AsTime().Sub(now); waitTime > 0 {
				logger.Debug().Msgf("going to sleep for %s", waitTime)
				timer.Reset(waitTime)
			} else {
				logger.Debug().Msgf("going to sleep for %s", defaultWaitDuration)
				timer.Reset(defaultWaitDuration)
			}
			select {
			case <-o.resetJobs:
				logger.Debug().Msgf("need to reset timer and check for new jobs")
			case <-timer.C:
				logger.Debug().Msgf("waking up")
			case <-ctx.Done():
				logger.Warn().Msg("received context done")
				wg.Wait()
				return
			}
		}
		wg.Wait()
	}()
	return out, errc
}

func (o *MonitorScheduledJobs) execute(ctx context.Context, num int, jobs <-chan schedulerrpc.ScheduledJob) <-chan error {
	errc := make(chan error, 1)
	logger := o.log.With().Str("component", "executor").Logger()
	go func() {
		defer func() {
			logger.Warn().Msgf("exiting executor %d", num)
		}()
		httpExecutor := NewHttpExecutor(logger)
		logger.Debug().Msgf("starting executor %d", num)
		for j := range jobs {
			logger.Info().Msgf("executing job id=%s", j.Id)
			var executor Executor
			if strings.HasPrefix(j.Endpoint, "http") {
				executor = httpExecutor
			} else {
				errc <- fmt.Errorf("not supported %s", j.Endpoint)
				return
			}
			resp, _ := executor.Execute(ctx, j)
			if err := o.repo.SaveExecutionResponse(ctx, j.Id, resp); err != nil {
				errc <- err
				return
			}
		}
	}()
	return errc
}
