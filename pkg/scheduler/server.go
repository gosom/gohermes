package scheduler

import (
	"context"

	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SchedulerServer struct {
	repo *schedulerRepository
	*schedulerrpc.UnimplementedScheduledJobServiceServer
}

func (o *SchedulerServer) CreateScheduledJob(ctx context.Context, req *schedulerrpc.CreateScheduledJobRequest) (*schedulerrpc.CreateScheduledJobResponse, error) {
	id, err := o.repo.insertScheduledJob(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := schedulerrpc.CreateScheduledJobResponse{Id: id}
	return &resp, nil
}

func recoveryFunc(p interface{}) (err error) {
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}

func apiAuthWrapper(cfg *Config) func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpc_auth.AuthFromMD(ctx, "x-api-key")
		if err != nil {
			return nil, err
		}
		if token != cfg.ApiKey {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}
		return ctx, nil
	}
}
