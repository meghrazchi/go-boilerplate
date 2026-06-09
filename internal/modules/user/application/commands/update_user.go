package commands

import "github.com/google/uuid"

type UpdateUserCommand struct {
	ID    uuid.UUID
	Name  string
	Email string
}
