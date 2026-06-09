package http

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handler.ListUsers)
		r.Post("/", handler.CreateUser)
		r.Get("/{id}", handler.GetUser)
		r.Put("/{id}", handler.UpdateUser)
		r.Delete("/{id}", handler.DeleteUser)
	})
}
