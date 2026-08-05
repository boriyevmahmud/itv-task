#!/usr/bin/env bash
# Idempotent repository bootstrap for the Movie API Cloud Agent environment.
# Installs PostgreSQL, prepares the .env file, and warms the Go module/build cache.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$REPO_ROOT"

echo "==> Installing system packages (PostgreSQL, Go toolchain deps)"
export DEBIAN_FRONTEND=noninteractive
sudo apt-get update -qq
sudo apt-get install -y -qq postgresql postgresql-contrib >/dev/null

echo "==> Ensuring .env exists"
if [ ! -f .env ]; then
  cat > .env <<'ENV'
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=movies
JWT_ACCESS_SECRET=dev_access_secret
JWT_REFRESH_SECRET=dev_refresh_secret

SERVICE_NAME=movies_service
ENV
  echo "    created .env"
else
  echo "    .env already present, leaving it untouched"
fi

echo "==> Downloading Go modules"
go mod download

echo "==> Building the application (warms the build cache)"
go build -o bin/movieapi ./cmd

echo "==> Install complete"
