package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/your-org/go-ddd-boilerplate/internal/app"
	"github.com/your-org/go-ddd-boilerplate/internal/config"
	"github.com/your-org/go-ddd-boilerplate/internal/platform/database"
	platformlogger "github.com/your-org/go-ddd-boilerplate/internal/platform/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, err := platformlogger.New(cfg)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() { _ = logger.Sync() }()

	db, sqlDB, err := database.Connect(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("connect database", zap.Error(err))
	}
	defer func() { _ = sqlDB.Close() }()

	handler := app.Bootstrap(cfg, logger, db, sqlDB)
	if err := app.Serve(cfg, logger, handler); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
