package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application/commands"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application/queries"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/interfaces/http/dto"
	"github.com/your-org/go-ddd-boilerplate/internal/platform/response"
	platformvalidator "github.com/your-org/go-ddd-boilerplate/internal/platform/validator"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

type Handler struct {
	service   *application.Service
	validator *platformvalidator.Validator
	logger    *zap.Logger
}

func NewHandler(service *application.Service, validator *platformvalidator.Validator, logger *zap.Logger) *Handler {
	return &Handler{service: service, validator: validator, logger: logger}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateUserRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	if validationErrors := h.validator.ValidateStruct(request); validationErrors != nil {
		response.Error(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	user, err := h.service.CreateUser(r.Context(), commands.CreateUserCommand{Name: request.Name, Email: request.Email})
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "User created successfully", dto.FromUser(user))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	user, err := h.service.GetUser(r.Context(), queries.GetUserQuery{ID: id})
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User fetched successfully", dto.FromUser(user))
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := normalizePage(parsePositiveInt(query.Get("page"), defaultPage))
	limit := normalizeLimit(parsePositiveInt(query.Get("limit"), defaultLimit))

	users, total, err := h.service.ListUsers(r.Context(), queries.ListUsersQuery{
		Page:   page,
		Limit:  limit,
		Search: query.Get("search"),
		Sort:   query.Get("sort"),
		Order:  query.Get("order"),
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Paginated(
		w,
		http.StatusOK,
		"Users fetched successfully",
		dto.FromUsers(users),
		response.NewPaginationMeta(page, limit, total),
	)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	var request dto.UpdateUserRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	if validationErrors := h.validator.ValidateStruct(request); validationErrors != nil {
		response.Error(w, http.StatusBadRequest, "Validation failed", validationErrors)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), commands.UpdateUserCommand{ID: id, Name: request.Name, Email: request.Email})
	if err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User updated successfully", dto.FromUser(user))
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseUUIDParam(w, r, "id")
	if !ok {
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, "User deleted successfully", nil)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		response.Error(w, http.StatusNotFound, "User not found", nil)
	case errors.Is(err, domain.ErrUserAlreadyExists):
		response.Error(w, http.StatusConflict, "User already exists", nil)
	case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrInvalidUserName):
		response.Error(w, http.StatusBadRequest, "Invalid user payload", nil)
	default:
		h.logger.Error("unhandled user error", zap.Error(err))
		response.Error(w, http.StatusInternalServerError, "Internal server error", nil)
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		if errors.Is(err, io.EOF) {
			response.Error(w, http.StatusBadRequest, "Request body must not be empty", nil)
			return false
		}
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload", map[string]string{"body": err.Error()})
		return false
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload", map[string]string{"body": "request body must contain a single JSON object"})
		return false
	}

	return true
}

func parseUUIDParam(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid UUID parameter", map[string]string{name: "must be a valid UUID"})
		return uuid.Nil, false
	}
	return id, true
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func normalizePage(page int) int {
	if page < 1 {
		return defaultPage
	}
	return page
}

func normalizeLimit(limit int) int {
	if limit < 1 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}
