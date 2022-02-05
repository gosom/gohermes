package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-stack/stack"
	"github.com/google/uuid"
	"github.com/gosom/gohermes/pkg/utils"
	"github.com/rs/zerolog"
)

func RequestId(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), utils.RidKey, rid)
		rw.Header().Add("X-Request-ID", rid)

		next.ServeHTTP(rw, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func Logger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			defer func() {
				logger.Info().
					Str("request-id", utils.GetReqID(r.Context())).
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("query", r.URL.RawQuery).
					Str("ip", r.RemoteAddr).
					Str("user-agent", r.UserAgent()).
					Dur("latency", time.Since(start)).
					Msg("request completed")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Recover(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					err, ok := p.(error)
					if !ok {
						err = fmt.Errorf("%v", p)
					}
					var stackTrace stack.CallStack
					traces := stack.Trace().TrimRuntime()
					for i := 0; i < len(traces); i++ {
						t := traces[i]
						tFunc := t.Frame().Function
						if tFunc == "runtime.gopanic" {
							continue
						}
						if tFunc == "net/http.HandlerFunc.ServeHTTP" {
							break
						}
						stackTrace = append(stackTrace, t)
					}

					logger.WithLevel(zerolog.PanicLevel).
						Err(err).
						Str("request-id", utils.GetReqID(r.Context())).
						Str("stack", fmt.Sprintf("%+v", stackTrace)).
						Msg(err.Error())

					http.Error(rw, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, r)
		}
		return http.HandlerFunc(fn)
	}
}
