#!/usr/bin/env bash
apt-get update; apt-get install -y curl

export LIB_MIGRATE_VERSION=v4.10.0
export LIB_PLATFORM_VERSION=linux
export SSL_MODE=?sslmode=disable

curl -L https://github.com/golang-migrate/migrate/releases/download/$LIB_MIGRATE_VERSION/migrate."$LIB_PLATFORM_VERSION"-amd64.tar.gz | tar xvz && mv migrate."$LIB_PLATFORM_VERSION"-amd64 migrate
chmod +x migrate

./migrate -source file://internal/postgres/migrations -database "$POSTGRES_ITEMS_CONNECTION$SSL_MODE" up

echo "Executing migrations in $POSTGRES_ITEMS_CONNECTION"

exec "$@"