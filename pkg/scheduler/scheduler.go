package scheduler

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/gosom/gohermes/pkg/container"
	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
)

type SchedulerServer struct {
	repo *schedulerRepository
	*schedulerrpc.UnimplementedScheduledJobServiceServer
}

func Run(ctx context.Context, di *container.ServiceContainer) error {
	repo, err := newSchedulerRepository(di.Cfg.DSN(), di.Logger)
	if err != nil {
		return err
	}
	if err := repo.migrate(ctx); err != nil {
		return err
	}
	schedulerServer := SchedulerServer{repo: repo}
	server := grpc.NewServer()
	schedulerrpc.RegisterScheduledJobServiceServer(server, &schedulerServer)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 1512))
	if err != nil {
		return err
	}
	return server.Serve(lis)
}

func (o *SchedulerServer) CreateScheduledJob(ctx context.Context, req *schedulerrpc.CreateScheduledJobRequest) (*schedulerrpc.CreateScheduledJobResponse, error) {
	return nil, nil
}
