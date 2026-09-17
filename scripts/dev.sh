#!/usr/bin/env bash

set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"

: "${JWT_SECRET:=local-development-secret}"
: "${BACKEND_ADDR:=:8080}"
: "${FRONTEND_ADDR:=3000}"
: "${NEXT_PUBLIC_API_URL:=http://localhost:8080}"

backend_pid=""
frontend_pid=""

cleanup() {
  trap - EXIT INT TERM
  [[ -z "$backend_pid" ]] || kill "$backend_pid" 2>/dev/null || true
  [[ -z "$frontend_pid" ]] || kill "$frontend_pid" 2>/dev/null || true
  wait "$backend_pid" "$frontend_pid" 2>/dev/null || true
}

trap cleanup EXIT INT TERM

if ! command -v go >/dev/null 2>&1; then
  echo "go is required" >&2
  exit 1
fi

if ! command -v npm >/dev/null 2>&1; then
  echo "npm is required" >&2
  exit 1
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "docker is required to start postgres" >&2
  exit 1
fi

if ! docker compose version >/dev/null 2>&1; then
  echo "docker compose is required to start postgres" >&2
  exit 1
fi

COMPOSE_FILE="$ROOT_DIR/docker-compose.yml"
echo "Starting postgres..."
docker compose -f "$COMPOSE_FILE" up -d postgres

echo "Waiting for postgres to become ready..."
postgres_ready=false
for _ in {1..30}; do
  if docker compose -f "$COMPOSE_FILE" exec -T postgres pg_isready -U postgres -d app >/dev/null 2>&1; then
    postgres_ready=true
    break
  fi
  sleep 1
done

if [[ "$postgres_ready" != true ]]; then
  echo "postgres did not become ready in time" >&2
  exit 1
fi

echo "Starting backend on ${BACKEND_ADDR}..."
(
  cd "$BACKEND_DIR"
  APP_ADDR="$BACKEND_ADDR" JWT_SECRET="$JWT_SECRET" go run ./cmd/server
) &
backend_pid=$!

echo "Starting frontend on http://localhost:${FRONTEND_ADDR}..."
(
  cd "$FRONTEND_DIR"
  NEXT_PUBLIC_API_URL="$NEXT_PUBLIC_API_URL" npm run dev -- --hostname 0.0.0.0 --port "$FRONTEND_ADDR"
) &
frontend_pid=$!

echo "Frontend: http://localhost:${FRONTEND_ADDR}"
echo "Backend:  http://localhost:8080"
echo "Press Ctrl+C to stop both services."

wait "$backend_pid" "$frontend_pid"
