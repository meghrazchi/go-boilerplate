package http

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/your-org/go-ddd-boilerplate/internal/platform/response"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, "Service is healthy", map[string]string{"status": "ok"})
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "Database is not ready", nil)
		return
	}

	response.Success(w, http.StatusOK, "Service is ready", map[string]string{"database": "ok"})
}
