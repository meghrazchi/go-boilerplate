package http_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
	userhttp "github.com/your-org/go-ddd-boilerplate/internal/modules/user/interfaces/http"
	platformvalidator "github.com/your-org/go-ddd-boilerplate/internal/platform/validator"
)

type testRepository struct {
	usersByID    map[uuid.UUID]*domain.User
	usersByEmail map[string]*domain.User
}

func newTestRepository() *testRepository {
	return &testRepository{usersByID: map[uuid.UUID]*domain.User{}, usersByEmail: map[string]*domain.User{}}
}

func (r *testRepository) Create(_ context.Context, user *domain.User) error {
	r.usersByID[user.ID()] = user
	r.usersByEmail[user.Email().String()] = user
	return nil
}

func (r *testRepository) FindByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *testRepository) FindByEmail(_ context.Context, email domain.Email) (*domain.User, error) {
	user, ok := r.usersByEmail[email.String()]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *testRepository) List(_ context.Context, _ domain.ListParams) ([]*domain.User, int64, error) {
	users := make([]*domain.User, 0, len(r.usersByID))
	for _, user := range r.usersByID {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

func (r *testRepository) Update(_ context.Context, user *domain.User) error {
	r.usersByID[user.ID()] = user
	r.usersByEmail[user.Email().String()] = user
	return nil
}

func (r *testRepository) Delete(_ context.Context, id uuid.UUID) error {
	delete(r.usersByID, id)
	return nil
}

func TestCreateUserHandlerValidationError(t *testing.T) {
	router := chi.NewRouter()
	service := application.NewService(newTestRepository())
	handler := userhttp.NewHandler(service, platformvalidator.New(), zap.NewNop())
	userhttp.RegisterRoutes(router, handler)

	request := httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBufferString(`{"name":"A","email":"bad"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestCreateUserHandlerSuccess(t *testing.T) {
	router := chi.NewRouter()
	service := application.NewService(newTestRepository())
	handler := userhttp.NewHandler(service, platformvalidator.New(), zap.NewNop())
	userhttp.RegisterRoutes(router, handler)

	request := httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBufferString(`{"name":"Ada Lovelace","email":"ada@example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}
}
