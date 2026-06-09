//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/your-org/go-ddd-boilerplate/internal/modules/user/domain"
	userpersistence "github.com/your-org/go-ddd-boilerplate/internal/modules/user/infrastructure/persistence"
)

func TestUserRepositoryCreateAndFind(t *testing.T) {
	db, cleanup := setupPostgres(t)
	defer cleanup()

	repository := userpersistence.NewGormRepository(db)

	email, err := domain.NewEmail("ada@example.com")
	if err != nil {
		t.Fatalf("create email: %v", err)
	}

	user, err := domain.NewUser("Ada Lovelace", email)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := repository.Create(context.Background(), user); err != nil {
		t.Fatalf("save user: %v", err)
	}

	found, err := repository.FindByID(context.Background(), user.ID())
	if err != nil {
		t.Fatalf("find user: %v", err)
	}

	if found.Email().String() != "ada@example.com" {
		t.Fatalf("expected email ada@example.com, got %s", found.Email().String())
	}
}
