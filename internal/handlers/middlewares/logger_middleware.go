package middlewares

import (
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func LoggerMiddleware(log *zap.SugaredLogger) mux.MiddlewareFunc {
	log = log.With(zap.String("comp", "http middleware"))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLog := log.
				With("method", r.Method).
				With("path", r.RequestURI).
				With("request ID", uuid.New().String())

			requestLog.Info("request started")

			startTime := time.Now()
			lrw := negroni.NewResponseWriter(w)
			next.ServeHTTP(lrw, r)
			duration := time.Since(startTime)

			status := lrw.Status()
			if status >= 500 {
				requestLog.
					With("status", status).
					Warnf("request failed with duration %f seconds", duration.Seconds())
			} else {
				requestLog.
					With("status", status).
					Infof("request finished with duration %f seconds", duration.Seconds())
			}
		})
	}
}
