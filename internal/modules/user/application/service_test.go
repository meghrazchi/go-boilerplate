package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application/commands"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application/queries"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
)

type inMemoryRepository struct {
	usersByID    map[uuid.UUID]*domain.User
	usersByEmail map[string]*domain.User
}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		usersByID:    map[uuid.UUID]*domain.User{},
		usersByEmail: map[string]*domain.User{},
	}
}

func (r *inMemoryRepository) Create(_ context.Context, user *domain.User) error {
	r.usersByID[user.ID()] = user
	r.usersByEmail[user.Email().String()] = user
	return nil
}

func (r *inMemoryRepository) FindByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryRepository) FindByEmail(_ context.Context, email domain.Email) (*domain.User, error) {
	user, ok := r.usersByEmail[email.String()]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *inMemoryRepository) List(_ context.Context, params domain.ListParams) ([]*domain.User, int64, error) {
	users := make([]*domain.User, 0, len(r.usersByID))
	for _, user := range r.usersByID {
		users = append(users, user)
	}
	return users, int64(len(users)), nil
}

func (r *inMemoryRepository) Update(_ context.Context, user *domain.User) error {
	old, ok := r.usersByID[user.ID()]
	if !ok {
		return domain.ErrUserNotFound
	}
	delete(r.usersByEmail, old.Email().String())
	r.usersByID[user.ID()] = user
	r.usersByEmail[user.Email().String()] = user
	return nil
}

func (r *inMemoryRepository) Delete(_ context.Context, id uuid.UUID) error {
	user, ok := r.usersByID[id]
	if !ok {
		return domain.ErrUserNotFound
	}
	delete(r.usersByID, id)
	delete(r.usersByEmail, user.Email().String())
	return nil
}

func TestCreateUser(t *testing.T) {
	service := application.NewService(newInMemoryRepository())

	user, err := service.CreateUser(context.Background(), commands.CreateUserCommand{Name: "Ada Lovelace", Email: "ADA@example.com"})
	if err != nil {
		t.Fatalf("expected user to be created, got %v", err)
	}
	if user.Email().String() != "ada@example.com" {
		t.Fatalf("expected normalized email, got %q", user.Email().String())
	}
}

func TestCreateUserRejectsDuplicateEmail(t *testing.T) {
	service := application.NewService(newInMemoryRepository())

	_, err := service.CreateUser(context.Background(), commands.CreateUserCommand{Name: "Ada Lovelace", Email: "ada@example.com"})
	if err != nil {
		t.Fatalf("create first user: %v", err)
	}

	_, err = service.CreateUser(context.Background(), commands.CreateUserCommand{Name: "Ada Byron", Email: "ada@example.com"})
	if err != domain.ErrUserAlreadyExists {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestListUsersNormalizesPagination(t *testing.T) {
	repository := newInMemoryRepository()
	service := application.NewService(repository)

	_, _, err := service.ListUsers(context.Background(), queries.ListUsersQuery{Page: -5, Limit: 999})
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
}
