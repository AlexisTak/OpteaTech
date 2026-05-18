#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_PORT="${API_PORT:-3001}"
WEB_PORT="${WEB_PORT:-3000}"
ADMIN_PORT="${ADMIN_PORT:-3002}"

echo ""
echo "Starting OpteaTech services..."
echo "- Go API:    http://127.0.0.1:$API_PORT"
echo "- Frontend:  http://127.0.0.1:$WEB_PORT"
echo "- Admin:     http://127.0.0.1:$ADMIN_PORT"
echo ""
echo "Press Ctrl+C to stop all services."
echo ""

PIDS=()

cleanup() {
  echo ""
  echo "Stopping all services..."
  for pid in "${PIDS[@]:-}"; do
    kill "$pid" >/dev/null 2>&1 || true
  done
  wait >/dev/null 2>&1 || true
}

trap cleanup EXIT INT TERM

start_service() {
  local name="$1"
  local command="$2"
  (
    bash -lc "$command" 2>&1 | sed -u "s/^/[$name] /"
  ) &
  PIDS+=("$!")
}

start_service "api" "cd '$ROOT_DIR/api' && PORT='$API_PORT' go run ./cmd/server"
start_service "web" "cd '$ROOT_DIR/apps/web' && npm run dev -- --port '$WEB_PORT'"
start_service "admin" "cd '$ROOT_DIR/apps/admin' && npm run dev -- --port '$ADMIN_PORT'"

wait
