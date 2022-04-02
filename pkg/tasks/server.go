package tasks

import (
	"github.com/gosom/gohermes/pkg/container"
	"github.com/hibiken/asynq"
)

type AsyncServer struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func (o AsyncServer) Run() error {
	return o.srv.Run(o.mux)
}

type Task struct {
	Pattern string
	Handler asynq.Handler
}

func NewDefaultAsyncWorker(di *container.ServiceContainer, bgtasks ...Task) AsyncServer {
	ans := AsyncServer{}
	ans.srv = asynq.NewServer(
		asynq.RedisClientOpt{Addr: di.Cfg.RedisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: di.Cfg.WorkerConcurrency,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)
	ans.mux = asynq.NewServeMux()
	for _, t := range bgtasks {
		ans.mux.Handle(t.Pattern, t.Handler)
	}
	return ans
}
