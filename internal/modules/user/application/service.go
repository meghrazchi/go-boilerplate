package application

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application/commands"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/application/queries"
	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateUser(ctx context.Context, command commands.CreateUserCommand) (*domain.User, error) {
	email, err := domain.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	_, err = s.repository.FindByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	user, err := domain.NewUser(command.Name, email)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, query queries.GetUserQuery) (*domain.User, error) {
	return s.repository.FindByID(ctx, query.ID)
}

func (s *Service) ListUsers(ctx context.Context, query queries.ListUsersQuery) ([]*domain.User, int64, error) {
	params := domain.ListParams{
		Page:   normalizePage(query.Page),
		Limit:  normalizeLimit(query.Limit),
		Search: strings.TrimSpace(query.Search),
		Sort:   normalizeSort(query.Sort),
		Order:  normalizeOrder(query.Order),
	}
	return s.repository.List(ctx, params)
}

func (s *Service) UpdateUser(ctx context.Context, command commands.UpdateUserCommand) (*domain.User, error) {
	email, err := domain.NewEmail(command.Email)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.FindByID(ctx, command.ID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repository.FindByEmail(ctx, email)
	if err == nil && existing.ID() != command.ID {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	if err := user.Update(command.Name, email); err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repository.Delete(ctx, id)
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeSort(sort string) string {
	switch sort {
	case "name", "email", "created_at", "updated_at":
		return sort
	default:
		return "created_at"
	}
}

func normalizeOrder(order string) string {
	if strings.EqualFold(order, "asc") {
		return "asc"
	}
	return "desc"
}
