package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"

	"github.com/gosom/gohermes/pkg/container"
)

type Server struct {
	di               *container.ServiceContainer
	log              zerolog.Logger
	srv              *http.Server
	router           *chi.Mux
	extraMiddlewares []func(http.Handler) http.Handler
}

func New(di *container.ServiceContainer) (*Server, error) {
	srv := &http.Server{
		Addr: di.Cfg.ServerAddress,
		// TODO
		//WriteTimeout: cfg.ServerWriteTimeout,
		//ReadTimeout:  cfg.ServerReadTimeout,
		//IdleTimeout:  cfg.ServerIdleTimeout,
	}
	ans := Server{
		di:  di,
		log: di.Logger.With().Str("component", "server").Logger(),
		srv: srv,
	}
	ans.router = chi.NewRouter()
	ans.setupMiddleware()
	ans.srv.Handler = ans.router
	return &ans, nil
}

func (o *Server) AddMidlewares(middlewares ...func(http.Handler) http.Handler) {
	o.extraMiddlewares = append(o.extraMiddlewares, middlewares...)
}

func (o *Server) Run() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	go func() {
		<-signalChan
		o.log.Warn().Msg("catched interrupt or kill signal")
		cancel()
	}()

	srvErrC := make(chan error, 1)
	go func() {
		defer close(srvErrC)
		o.log.Info().Msg(fmt.Sprintf("Starting listening on %s", o.srv.Addr))
		err = o.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			srvErrC <- err
			return
		} else {
			srvErrC <- nil
			return
		}
	}()

	<-ctx.Done()
	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	err = o.srv.Shutdown(gracefullCtx)
	if err != nil {
		o.log.Warn().Msg(err.Error())
		return
	}
	err = <-srvErrC
	if err == nil {
		o.log.Warn().Msg("server gracefully exited")
	} else {
		o.log.Warn().Msgf("server exited with error: %s", err.Error())
	}
	return
}

func (o *Server) setupMiddleware() {
	o.router.Use(middleware.RealIP)
	o.router.Use(RequestId)
	for i := range o.extraMiddlewares {
		o.router.Use(o.extraMiddlewares[i])
	}
	o.router.Use(Logger(o.log))
	o.router.Use(Recover(o.log))
}

func (o *Server) GetRouter() *chi.Mux {
	return o.router
}
