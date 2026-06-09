package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/your-org/go-ddd-boilerplate/internal/config"
)

func Serve(cfg *config.Config, logger *zap.Logger, handler http.Handler) error {
	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       cfg.HTTPReadTimeoutDuration(),
		WriteTimeout:      cfg.HTTPWriteTimeoutDuration(),
		IdleTimeout:       cfg.HTTPIdleTimeoutDuration(),
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("http server starting", zap.String("address", cfg.Address()))
		errCh <- server.ListenAndServe()
	}()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-shutdownCh:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeoutDuration())
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
		return err
	}

	logger.Info("http server stopped")
	return http.ErrServerClosed
}
