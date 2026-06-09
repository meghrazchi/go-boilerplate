package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return chimiddleware.Timeout(timeout)
}
