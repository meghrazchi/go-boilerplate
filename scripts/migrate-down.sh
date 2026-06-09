#!/usr/bin/env sh
set -eu

: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_USER:=postgres}"
: "${DB_PASSWORD:=postgres}"
: "${DB_NAME:=go_boilerplate}"
: "${DB_SSL_MODE:=disable}"

DBSTRING="host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}"
goose -dir migrations postgres "${DBSTRING}" down
