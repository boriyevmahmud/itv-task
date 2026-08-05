#!/usr/bin/env bash
# Per-boot PostgreSQL reconciliation: start the cluster, ensure the role
# password and the application database exist, then wait for readiness.
# Safe to run repeatedly.
set -euo pipefail

PG_VERSION="$(ls /etc/postgresql 2>/dev/null | sort -V | tail -1 || true)"
PG_VERSION="${PG_VERSION:-16}"
DB_NAME="movies"
DB_PASSWORD="postgres"

echo "==> Starting PostgreSQL cluster ${PG_VERSION}/main"
sudo pg_ctlcluster "${PG_VERSION}" main start 2>/dev/null || true

echo "==> Waiting for PostgreSQL to accept connections"
for i in $(seq 1 30); do
  if sudo -u postgres pg_isready -q; then
    break
  fi
  sleep 1
done
sudo -u postgres pg_isready

echo "==> Ensuring 'postgres' role password"
sudo -u postgres psql -v ON_ERROR_STOP=1 -c "ALTER USER postgres WITH PASSWORD '${DB_PASSWORD}';"

echo "==> Ensuring '${DB_NAME}' database exists"
if ! sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}'" | grep -q 1; then
  sudo -u postgres createdb "${DB_NAME}"
  echo "    created database ${DB_NAME}"
else
  echo "    database ${DB_NAME} already exists"
fi

echo "==> PostgreSQL is ready"
