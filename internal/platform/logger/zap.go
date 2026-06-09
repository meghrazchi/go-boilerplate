package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/your-org/go-ddd-boilerplate/internal/config"
)

func New(cfg *config.Config) (*zap.Logger, error) {
	var zapCfg zap.Config
	if cfg.IsProduction() {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	level := zapcore.DebugLevel
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.LogLevel))); err != nil {
		level = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	return zapCfg.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
