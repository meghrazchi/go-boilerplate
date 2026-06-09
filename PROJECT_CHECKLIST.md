# Project checklist

## Included

- [x] DDD-oriented folder structure
- [x] Chi HTTP router
- [x] GORM PostgreSQL repository implementation
- [x] Domain entity/value object separation from GORM and HTTP
- [x] DTO request/response types
- [x] go-playground/validator integration
- [x] Central JSON response envelopes
- [x] Central domain-to-HTTP error mapping
- [x] Zap structured logging
- [x] Request ID, logging, recovery, CORS, timeout, security headers, and body-size middleware
- [x] envconfig typed configuration with `.env.example`
- [x] PostgreSQL connection pooling and ping checks
- [x] Goose SQL migration with PostgreSQL constraints and updated_at trigger
- [x] Health and readiness endpoints
- [x] OpenAPI document and Swagger UI endpoint
- [x] Multi-stage Dockerfile with non-root runtime
- [x] Docker Compose with Postgres healthcheck, named volume, bridge network, migration service, API healthcheck, read-only API filesystem, tmpfs, dropped capabilities, and no-new-privileges
- [x] Makefile commands for local and Docker workflows
- [x] golangci-lint config
- [x] pre-commit config
- [x] GitHub Actions CI
- [x] Unit, handler, and integration-test examples

## Local verification commands

```bash
cp .env.example .env
go mod tidy
make install-tools
make verify
make docker-init
make docker-migrate-status
make test-integration
```

## Known sandbox limitation

This archive was generated in an environment where external Go module downloads and Docker execution are unavailable. Because of that, the included `go.sum` is an empty placeholder and Docker containers cannot be started here. Run `go mod tidy` once locally or in CI to populate `go.sum`.

## Docker build troubleshooting

If you see `go.mod requires go >= 1.25.0` while Docker says it is running Go 1.23.x, rebuild with the updated Docker build argument/cache cleared:

```bash
docker compose build --no-cache
make docker-init
```

The expected Docker build image is `golang:1.25-alpine`.

### Docker Go version note

This project requires Go 1.25. Docker Compose builds with `DOCKER_GO_VERSION=1.25` by default. If an older local `.env` contains `GO_VERSION=1.23`, update it to `GO_VERSION=1.25` or remove that line.
