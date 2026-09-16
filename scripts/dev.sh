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
