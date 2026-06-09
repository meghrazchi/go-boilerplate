package app

import (
	"database/sql"
	"net/http"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/your-org/go-ddd-boilerplate/internal/config"
	healthhttp "github.com/your-org/go-ddd-boilerplate/internal/modules/health/interfaces/http"
	userapp "github.com/your-org/go-ddd-boilerplate/internal/modules/user/application"
	userpersistence "github.com/your-org/go-ddd-boilerplate/internal/modules/user/infrastructure/persistence"
	userhttp "github.com/your-org/go-ddd-boilerplate/internal/modules/user/interfaces/http"
	platformvalidator "github.com/your-org/go-ddd-boilerplate/internal/platform/validator"
	"github.com/your-org/go-ddd-boilerplate/internal/routes"
)

func Bootstrap(cfg *config.Config, logger *zap.Logger, db *gorm.DB, sqlDB *sql.DB) http.Handler {
	validator := platformvalidator.New()

	userRepository := userpersistence.NewGormRepository(db)
	userService := userapp.NewService(userRepository)
	userHandler := userhttp.NewHandler(userService, validator, logger)

	healthHandler := healthhttp.NewHandler(sqlDB)

	return routes.NewRouter(routes.Dependencies{
		Config:        cfg,
		Logger:        logger,
		UserHandler:   userHandler,
		HealthHandler: healthHandler,
	})
}
