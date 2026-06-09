package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/your-org/go-ddd-boilerplate/internal/platform/response"
)

func Recovery(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered", zap.Any("panic", recovered))
					response.Error(w, http.StatusInternalServerError, "Internal server error", nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
