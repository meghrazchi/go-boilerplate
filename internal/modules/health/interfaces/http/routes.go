package http

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Get("/health", handler.Health)
	r.Get("/ready", handler.Ready)
}
