package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
	"github.com/rs/zerolog"
)

type ExecutorResponse struct {
	Data interface{} `json:"data"`
	Err  string      `json:"error"`
}

type Executor interface {
	Execute(ctx context.Context, job schedulerrpc.ScheduledJob) (ExecutorResponse, error)
}

type HttpExecutor struct {
	log    zerolog.Logger
	client *http.Client
	pool   sync.Pool
}

func NewHttpExecutor(logger zerolog.Logger) *HttpExecutor {
	ans := HttpExecutor{
		log: logger,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		pool: sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}
	return &ans
}

func (o *HttpExecutor) Execute(ctx context.Context, job schedulerrpc.ScheduledJob) (ExecutorResponse, error) {
	var ans ExecutorResponse
	buf := o.pool.Get().(*bytes.Buffer)
	buf.Reset()
	defer o.pool.Put(buf)
	buf.Write(job.Data)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, job.Endpoint, buf)
	if err != nil {
		o.log.Error().Msgf("error executing %s -> %s : %s", job.Id, job.Endpoint, err.Error())
		ans.Err = err.Error()
		return ans, err
	}
	resp, err := o.client.Do(req)
	if err != nil {
		o.log.Error().Msgf("error executing %s -> %s : %s", job.Id, job.Endpoint, err.Error())
		ans.Err = err.Error()
		return ans, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("invalid status code: %d", resp.StatusCode)
		o.log.Error().Msgf("error executing %s -> %s : %s", job.Id, job.Endpoint, err.Error())
		ans.Err = err.Error()
		return ans, err

	}
	if err := json.NewDecoder(resp.Body).Decode(ans.Data); err != nil {
		o.log.Error().Msgf("error executing %s -> %s : %s", job.Id, job.Endpoint, err.Error())
		ans.Err = err.Error()
		return ans, err
	}
	return ans, nil
}
