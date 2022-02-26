package scheduler

import (
	"context"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/gosom/gohermes/pkg/scheduler/schedulerrpc"
)

func Run(ctx context.Context, cfg *Config) error {
	logger := getLogger(cfg)
	repo, err := newSchedulerRepository(cfg.DSN(), logger)
	if err != nil {
		return err
	}
	if err := repo.migrate(ctx); err != nil {
		return err
	}

	monitor, err := NewMonitorScheduledJobs(
		logger, repo, cfg.Executors,
	)
	if err != nil {
		return err
	}

	server, err := getGrpcServer(cfg, logger, repo)
	if err != nil {
		return err
	}
	lis, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := monitor.Run(ctx); err != nil {
			logger.Error().Msgf("monitor exited: %s", err.Error())
		}
	}()
	// TODO return an error channel so we will return the error
	go func() {
		defer wg.Done()
		if err := server.Serve(lis); err != nil {
			logger.Error().Msgf("grpcServer: %s", err.Error())
		} else {
			logger.Info().Msgf("grcpServer: clean exit")
		}
	}()
	wg.Wait()
	return nil
}

func getGrpcServer(cfg *Config, logger zerolog.Logger, repo *schedulerRepository) (*grpc.Server, error) {
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
	return server, nil
}

func getLogger(cfg *Config) zerolog.Logger {
	var level zerolog.Level
	if cfg.Debug {
		level = zerolog.DebugLevel
	} else {
		level = zerolog.InfoLevel
	}
	logger := log.Level(level)
	return logger
}
