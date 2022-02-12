package server

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
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
	var bufferPool = sync.Pool{
		New: func() interface{} { return new(bytes.Buffer) },
	}
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			buf := bufferPool.Get().(*bytes.Buffer)
			buf.Reset()
			defer bufferPool.Put(buf)
			ww.Tee(buf)
			defer func() {
				statusCode := ww.Status()
				var ev *zerolog.Event
				msg := http.StatusText(statusCode)
				if statusCode >= 200 && statusCode < 400 {
					ev = logger.Info()
				} else if statusCode >= 400 && statusCode < 500 {
					ev = logger.Warn()
				} else {
					ev = logger.Error()
				}
				ev = ev.
					Str("request-id", utils.GetReqID(r.Context())).
					Int("status", statusCode).
					Int("bytes", ww.BytesWritten()).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("query", r.URL.RawQuery).
					Str("ip", r.RemoteAddr).
					Str("user-agent", r.UserAgent()).
					Dur("latency", time.Since(start))
				if statusCode < 200 || statusCode >= 400 {
					ev = ev.RawJSON("body", buf.Bytes())
				}
				ev.Msg(msg)
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
