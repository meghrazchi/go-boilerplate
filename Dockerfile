# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25
ARG ALPINE_VERSION=3.20

FROM golang:${GO_VERSION}-alpine AS base
WORKDIR /src
RUN apk add --no-cache ca-certificates git tzdata
COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS builder
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM base AS goose-builder
ARG GOOSE_VERSION=v3.24.3
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@${GOOSE_VERSION}

FROM alpine:${ALPINE_VERSION} AS migrator
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H appuser
WORKDIR /app
COPY --from=goose-builder /go/bin/goose /usr/local/bin/goose
COPY --chown=appuser:appuser migrations /app/migrations
USER appuser
ENTRYPOINT ["goose", "-dir", "/app/migrations"]
CMD ["status"]

FROM alpine:${ALPINE_VERSION} AS runtime
RUN apk add --no-cache ca-certificates tzdata && adduser -D -H appuser
WORKDIR /app
COPY --from=builder --chown=appuser:appuser /out/api /app/api
COPY --chown=appuser:appuser docs /app/docs
COPY --chown=appuser:appuser migrations /app/migrations
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:${APP_PORT:-8080}/health >/dev/null || exit 1
ENTRYPOINT ["/app/api"]
