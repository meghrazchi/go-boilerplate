package config

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppName  string `envconfig:"APP_NAME" default:"go-ddd-boilerplate"`
	AppEnv   string `envconfig:"APP_ENV" default:"development"`
	AppPort  string `envconfig:"APP_PORT" default:"8080"`
	AppDebug bool   `envconfig:"APP_DEBUG" default:"true"`

	DBHost         string `envconfig:"DB_HOST" default:"localhost"`
	DBPort         string `envconfig:"DB_PORT" default:"5432"`
	DBUser         string `envconfig:"DB_USER" default:"postgres"`
	DBPassword     string `envconfig:"DB_PASSWORD" default:"postgres"`
	DBName         string `envconfig:"DB_NAME" default:"go_boilerplate"`
	DBSSLMode      string `envconfig:"DB_SSL_MODE" default:"disable"`
	DBTimezone     string `envconfig:"DB_TIMEZONE" default:"UTC"`
	DBMaxIdleConns int    `envconfig:"DB_MAX_IDLE_CONNS" default:"10"`
	DBMaxOpenConns int    `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
	DBMaxLifetime  int    `envconfig:"DB_MAX_LIFETIME_SECONDS" default:"300"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`

	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"http://localhost:3000,http://localhost:5173"`
	CORSAllowedMethods []string `envconfig:"CORS_ALLOWED_METHODS" default:"GET,POST,PUT,PATCH,DELETE,OPTIONS"`
	CORSAllowedHeaders []string `envconfig:"CORS_ALLOWED_HEADERS" default:"Accept,Authorization,Content-Type,X-Request-ID"`

	RequestTimeoutSeconds   int   `envconfig:"REQUEST_TIMEOUT_SECONDS" default:"30"`
	ShutdownTimeoutSeconds  int   `envconfig:"SHUTDOWN_TIMEOUT_SECONDS" default:"10"`
	HTTPReadTimeoutSeconds  int   `envconfig:"HTTP_READ_TIMEOUT_SECONDS" default:"10"`
	HTTPWriteTimeoutSeconds int   `envconfig:"HTTP_WRITE_TIMEOUT_SECONDS" default:"15"`
	HTTPIdleTimeoutSeconds  int   `envconfig:"HTTP_IDLE_TIMEOUT_SECONDS" default:"60"`
	MaxBodyBytes            int64 `envconfig:"MAX_BODY_BYTES" default:"1048576"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.AppPort) == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if _, err := strconv.Atoi(c.AppPort); err != nil {
		return fmt.Errorf("APP_PORT must be numeric")
	}
	if strings.TrimSpace(c.DBHost) == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if strings.TrimSpace(c.DBUser) == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if strings.TrimSpace(c.DBName) == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DBMaxIdleConns < 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be negative")
	}
	if c.DBMaxOpenConns <= 0 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be greater than zero")
	}
	if c.DBMaxIdleConns > c.DBMaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS cannot be greater than DB_MAX_OPEN_CONNS")
	}
	if c.RequestTimeoutSeconds <= 0 {
		return fmt.Errorf("REQUEST_TIMEOUT_SECONDS must be greater than zero")
	}
	if c.ShutdownTimeoutSeconds <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT_SECONDS must be greater than zero")
	}
	if c.MaxBodyBytes <= 0 {
		return fmt.Errorf("MAX_BODY_BYTES must be greater than zero")
	}
	return nil
}

func (c Config) Address() string {
	return ":" + c.AppPort
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

func (c Config) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBSSLMode,
		c.DBTimezone,
	)
}

func (c Config) DatabaseURL() string {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   fmt.Sprintf("%s:%s", c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	query := dsn.Query()
	query.Set("sslmode", c.DBSSLMode)
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

func (c Config) DBMaxLifetimeDuration() time.Duration {
	return time.Duration(c.DBMaxLifetime) * time.Second
}

func (c Config) RequestTimeoutDuration() time.Duration {
	return time.Duration(c.RequestTimeoutSeconds) * time.Second
}

func (c Config) ShutdownTimeoutDuration() time.Duration {
	return time.Duration(c.ShutdownTimeoutSeconds) * time.Second
}

func (c Config) HTTPReadTimeoutDuration() time.Duration {
	return time.Duration(c.HTTPReadTimeoutSeconds) * time.Second
}

func (c Config) HTTPWriteTimeoutDuration() time.Duration {
	return time.Duration(c.HTTPWriteTimeoutSeconds) * time.Second
}

func (c Config) HTTPIdleTimeoutDuration() time.Duration {
	return time.Duration(c.HTTPIdleTimeoutSeconds) * time.Second
}
