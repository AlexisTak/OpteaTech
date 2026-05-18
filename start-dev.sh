#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_PORT="${API_PORT:-3001}"
BFF_PORT="${BFF_PORT:-4000}"
WEB_PORT="${WEB_PORT:-3000}"

for cmd in go npm sed; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Missing required command: $cmd"
    exit 1
  fi
done

check_port_free() {
  local port="$1"
  if ss -ltn | grep -qE ":${port}\\b"; then
    echo "Port ${port} is already in use."
    echo "Stop the existing process on this port or run with a different port value."
    echo "Example: API_PORT=3101 BFF_PORT=4100 WEB_PORT=3100 ./start-dev.sh"
    exit 1
  fi
}

check_port_free "$API_PORT"
check_port_free "$BFF_PORT"
check_port_free "$WEB_PORT"

PIDS=()

cleanup() {
  echo
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
start_service "bff" "cd '$ROOT_DIR/backend' && PORT='$BFF_PORT' GO_API_URL='http://127.0.0.1:$API_PORT' npm run start:dev"
start_service "web" "cd '$ROOT_DIR/frontend' && NEXT_PUBLIC_BFF_URL='http://127.0.0.1:$BFF_PORT' npm run dev -- --port '$WEB_PORT'"

echo ""
echo "Services starting..."
echo "- Go API:    http://127.0.0.1:$API_PORT"
echo "- Nest BFF:  http://127.0.0.1:$BFF_PORT/api/v1"
echo "- Swagger:   http://127.0.0.1:$BFF_PORT/docs"
echo "- Frontend:  http://127.0.0.1:$WEB_PORT"
echo "- Open this in browser for website: http://127.0.0.1:$WEB_PORT"
echo ""
echo "Press Ctrl+C to stop all services."

wait
