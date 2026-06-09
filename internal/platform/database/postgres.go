package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/your-org/go-ddd-boilerplate/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*gorm.DB, *sql.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(cfg.DatabaseDSN()), &gorm.Config{
		Logger:         gormlogger.Default.LogMode(gormlogger.Silent),
		TranslateError: true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open postgres connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("extract sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetimeDuration())

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, nil, fmt.Errorf("ping postgres: %w", err)
	}

	logger.Info("postgres connected", zap.String("host", cfg.DBHost), zap.String("database", cfg.DBName))
	return gormDB, sqlDB, nil
}
