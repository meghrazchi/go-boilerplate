package dto

import (
	"time"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromUser(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID().String(),
		Name:      user.Name(),
		Email:     user.Email().String(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func FromUsers(users []*domain.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, FromUser(user))
	}
	return responses
}
