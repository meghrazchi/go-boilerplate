package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	if r.statusCode == 0 {
		r.statusCode = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(data)
	r.bytes += n
	return n, err
}

func Logger(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &responseRecorder{ResponseWriter: w}

			next.ServeHTTP(recorder, r)

			statusCode := recorder.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			logger.Info("http request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", statusCode),
				zap.Int("bytes", recorder.bytes),
				zap.Duration("duration", time.Since(start)),
				zap.String("request_id", chimiddleware.GetReqID(r.Context())),
				zap.String("remote_addr", r.RemoteAddr),
			)
		})
	}
}
