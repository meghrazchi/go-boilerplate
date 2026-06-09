package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func RequestID(next http.Handler) http.Handler {
	return chimiddleware.RequestID(next)
}
