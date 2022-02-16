package scheduler

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
)

type SchedulerServer struct {
	repo *schedulerRepository
	*schedulerrpc.UnimplementedScheduledJobServiceServer
}

func Run(ctx context.Context, cfg *Config) error {
	var level zerolog.Level
	if cfg.Debug {
		level = zerolog.DebugLevel
	} else {
		level = zerolog.InfoLevel
	}
	logger := log.Level(level)
	repo, err := newSchedulerRepository(cfg.DSN(), logger)
	if err != nil {
		return err
	}
	if err := repo.migrate(ctx); err != nil {
		return err
	}
	schedulerServer := SchedulerServer{repo: repo}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(apiAuthWrapper(cfg)),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_auth.StreamServerInterceptor(apiAuthWrapper(cfg)),
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	)
	schedulerrpc.RegisterScheduledJobServiceServer(server, &schedulerServer)
	lis, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return err
	}
	return server.Serve(lis)
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

func (o *SchedulerServer) CreateScheduledJob(ctx context.Context, req *schedulerrpc.CreateScheduledJobRequest) (*schedulerrpc.CreateScheduledJobResponse, error) {
	id, err := o.repo.insert(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := schedulerrpc.CreateScheduledJobResponse{Id: id}
	return &resp, nil
}
