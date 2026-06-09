//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	userpersistence "github.com/your-org/go-ddd-boilerplate/internal/modules/user/infrastructure/persistence"
)

func setupPostgres(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("go_boilerplate_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	db, err := gorm.Open(gormpostgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm connection: %v", err)
	}

	if err := db.AutoMigrate(&userpersistence.UserModel{}); err != nil {
		t.Fatalf("migrate test database: %v", err)
	}

	cleanup := func() {
		_ = container.Terminate(ctx)
	}
	return db, cleanup
}
