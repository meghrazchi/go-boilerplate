package domain

import (
	"context"

	"github.com/google/uuid"
)

type ListParams struct {
	Page   int
	Limit  int
	Search string
	Sort   string
	Order  string
}

func (p ListParams) Offset() int {
	if p.Page <= 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	List(ctx context.Context, params ListParams) ([]*User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
