-include .env

APP_NAME ?= go-ddd-boilerplate
APP_PORT ?= 8080
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= go_boilerplate
DB_SSL_MODE ?= disable
GO_VERSION ?= 1.25
DOCKER_GO_VERSION ?= 1.25
GOOSE_VERSION ?= v3.24.3
GOLANGCI_LINT_VERSION ?= v2.10.1

DATABASE_DSN := host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSL_MODE)
DOCKER_COMPOSE := DOCKER_GO_VERSION=$(DOCKER_GO_VERSION) docker compose

.PHONY: help run build test test-unit test-integration test-coverage lint fmt tidy verify docker-build docker-up docker-down docker-logs docker-ps docker-init docker-migrate-up docker-migrate-down docker-migrate-status migrate-up migrate-down migrate-status migrate-create install-tools precommit-install clean

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-24s %s\n", $$1, $$2}'

run: ## Run API locally
	go run ./cmd/api

build: ## Build API binary
	go build -trimpath -o bin/api ./cmd/api

test: test-unit ## Alias for unit/handler tests

test-unit: ## Run unit and handler tests
	go test ./...

test-integration: ## Run integration tests with testcontainers
	go test -tags=integration ./...

test-coverage: ## Run tests with coverage profile
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format Go files
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')
	test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './vendor/*'))"

tidy: ## Tidy Go modules
	go mod tidy

verify: fmt tidy lint test-unit build ## Run local verification checks

docker-build: ## Build Docker images
	$(DOCKER_COMPOSE) build

docker-up: ## Start API and Postgres
	$(DOCKER_COMPOSE) up --build -d postgres api

docker-down: ## Stop and remove Docker services
	$(DOCKER_COMPOSE) down

docker-logs: ## Follow API and Postgres logs
	$(DOCKER_COMPOSE) logs -f api postgres

docker-ps: ## Show Docker service status
	$(DOCKER_COMPOSE) ps

docker-init: ## Start Postgres, run migrations, then start API
	$(DOCKER_COMPOSE) up -d postgres
	$(DOCKER_COMPOSE) --profile tools run --rm migrate up
	$(DOCKER_COMPOSE) up --build -d api

docker-migrate-up: ## Apply migrations in Docker
	$(DOCKER_COMPOSE) --profile tools run --rm migrate up

docker-migrate-down: ## Roll back one migration in Docker
	$(DOCKER_COMPOSE) --profile tools run --rm migrate down

docker-migrate-status: ## Show Docker migration status
	$(DOCKER_COMPOSE) --profile tools run --rm migrate status

migrate-up: ## Apply goose migrations locally
	goose -dir migrations postgres "$(DATABASE_DSN)" up

migrate-down: ## Roll back one goose migration locally
	goose -dir migrations postgres "$(DATABASE_DSN)" down

migrate-status: ## Show goose migration status locally
	goose -dir migrations postgres "$(DATABASE_DSN)" status

migrate-create: ## Create migration: make migrate-create name=create_table
	@test -n "$(name)" || (echo "name is required. Example: make migrate-create name=create_users_table" && exit 1)
	goose -dir migrations create $(name) sql

install-tools: ## Install local dev tools
	go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

precommit-install: ## Install pre-commit hooks
	pre-commit install

clean: ## Remove local build outputs
	rm -rf bin coverage.out
