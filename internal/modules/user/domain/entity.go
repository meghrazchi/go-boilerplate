package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id        uuid.UUID
	name      string
	email     Email
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(name string, email Email) (*User, error) {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return nil, ErrInvalidUserName
	}

	now := time.Now().UTC()
	return &User{
		id:        uuid.New(),
		name:      name,
		email:     email,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func RehydrateUser(id uuid.UUID, name string, email Email, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) Update(name string, email Email) error {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return ErrInvalidUserName
	}
	u.name = name
	u.email = email
	u.updatedAt = time.Now().UTC()
	return nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}
