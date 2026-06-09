# Go DDD Boilerplate

> **Go version:** This boilerplate targets Go 1.25 because the selected current dependency/toolchain set requires Go 1.25 or newer. Docker Compose builds with `golang:1.25-alpine` by default.


A production-minded Go REST API boilerplate using Chi, GORM, PostgreSQL, envconfig, Zap, goose, OpenAPI, Docker Compose, golangci-lint, pre-commit, GitHub Actions, httptest, and testcontainers.

## Architecture

```text
cmd/api                application entrypoint
internal/app           bootstrap and server lifecycle
internal/config        typed envconfig configuration
internal/platform      shared infrastructure utilities
internal/modules       DDD bounded contexts
migrations             goose SQL migrations
docs                   OpenAPI/Swagger assets
tests                  integration/e2e tests
```

Dependency direction:

```text
interfaces/http -> application -> domain
infrastructure  -> domain
```

The domain layer does not import Chi, GORM, Zap, validators, or PostgreSQL packages.

## Stack

```text
Language: Go
Router: Chi
ORM: GORM
Database: PostgreSQL
Validation: go-playground/validator
Config: envconfig + optional .env loading
Logger: zap
Migrations: goose
Docs: Swagger/OpenAPI
Lint: golangci-lint
Pre-commit: pre-commit
Docker: Docker + Docker Compose
CI: GitHub Actions
Testing: go test + httptest + testcontainers
```

## First run with Docker

```bash
cp .env.example .env
go mod tidy
make docker-init
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

`make docker-init` starts PostgreSQL, waits for the database health check, runs goose migrations in a one-shot migration container, then starts the API.

Useful Docker commands:

```bash
make docker-up              # start API + Postgres
make docker-migrate-up      # apply migrations
make docker-migrate-status  # show migration status
make docker-logs            # follow logs
make docker-down            # stop services
```

## First run locally

```bash
cp .env.example .env
go mod tidy
make install-tools
make docker-up              # starts Postgres and API; use Ctrl+C if foreground is preferred
make migrate-up             # local goose command against localhost Postgres
make run
```

## API routes

```text
GET     /health
GET     /ready
GET     /docs/openapi.yaml
GET     /docs/swagger
GET     /api/v1/users
POST    /api/v1/users
GET     /api/v1/users/{id}
PUT     /api/v1/users/{id}
DELETE  /api/v1/users/{id}
```

Example request:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Ada Lovelace","email":"ada@example.com"}'
```

## Verification

```bash
make verify
make test-integration
```

`make verify` runs formatting, module tidy, lint, unit/handler tests, and build. Integration tests require Docker because they use testcontainers.

## Production notes

Before deploying, replace default credentials, set `APP_ENV=production`, set a restrictive `CORS_ALLOWED_ORIGINS`, use managed secrets instead of committing `.env`, and run migrations as a controlled release step.

The included Docker image uses multi-stage builds, a non-root runtime user, health checks, read-only filesystem settings in Compose, dropped Linux capabilities, and a named PostgreSQL volume.

### Docker Go version note

This project requires Go 1.25. Docker Compose builds with `DOCKER_GO_VERSION=1.25` by default. If an older local `.env` contains `GO_VERSION=1.23`, update it to `GO_VERSION=1.25` or remove that line.
