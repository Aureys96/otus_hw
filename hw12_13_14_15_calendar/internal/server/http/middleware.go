package internalhttp

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func setContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func (rw *statusWriter) WriteHeader(code int) {
	rw.status = code
}

func logRequest(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			from := time.Now()

			proxiedRequest := &statusWriter{ResponseWriter: w}
			next.ServeHTTP(proxiedRequest, r)
			logger.Info("http request",
				zap.String("ip", r.RemoteAddr),
				zap.String("timeStamp", from.String()),
				zap.String("method", r.Method),
				zap.Int("code", proxiedRequest.status),
				zap.Duration("latency", time.Since(from)),
			)
		})
	}
}

func handlePanic(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Error("internal error")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
