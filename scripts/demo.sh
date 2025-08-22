#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
HOST="${HOST:-localhost}"
PORT="${PORT:-8080}"
BASE_URL="http://${HOST}:${PORT}"
DEMO_DIR="${ROOT_DIR}/.tmp"
DB_PATH="${DEMO_DIR}/app.db"
DEBUG_LOG="${DEBUG_LOG:-false}"

# Parse command line arguments
RUN_TESTS=true
RUN_API=true

case "${1:-}" in
  "test"|"tests")
    RUN_TESTS=true
    RUN_API=false
    ;;
  "api")
    RUN_TESTS=false
    RUN_API=true
    ;;
  ""|"demo"|"all")
    RUN_TESTS=true
    RUN_API=true
    ;;
  *)
    echo "Usage: $0 [test|api|demo]"
    echo "  test - run only unit tests"
    echo "  api  - run only API demo"
    echo "  demo - run tests + API demo (default)"
    exit 1
    ;;
esac

pretty() {
  if command -v jq >/dev/null 2>&1; then
    jq . || cat
  else
    cat
  fi
}

# Run tests if requested
if [[ "$RUN_TESTS" == "true" ]]; then
  echo "==> Running unit tests"
  if command -v richgo >/dev/null 2>&1; then
    make -C "$ROOT_DIR" test-full-log | sed 's/^/    /'
  else
    make -C "$ROOT_DIR" test | sed 's/^/    /'
  fi
fi

echo
echo "==> Starting server on ${HOST}:${PORT} in background with DB at ${DB_PATH}"
mkdir -p "$DEMO_DIR"
rm -f "$DB_PATH"

start_server() {
  local p="$1"
  (
    cd "$ROOT_DIR"
    HOST="$HOST" PORT="$p" DATABASE_URL="$DB_PATH" go run cmd/server/main.go
  ) >"${DEMO_LOG:-/dev/null}" 2>&1 &
  SERVER_PID=$!
}

is_pid_alive() {
  kill -0 "$1" >/dev/null 2>&1
}

# Check if a TCP port is already in use
port_in_use() {
  local p="$1"
  if command -v lsof >/dev/null 2>&1; then
    lsof -nP -iTCP:"${p}" -sTCP:LISTEN >/dev/null 2>&1
    return $?
  elif command -v nc >/dev/null 2>&1; then
    nc -z -w 1 "${HOST}" "${p}" >/dev/null 2>&1
    return $?
  else
    return 1
  fi
}

# Try to start on PORT, fallback to next ports if busy
ATTEMPTS=10
for i in $(seq 0 $((ATTEMPTS-1))); do
  TRY_PORT=$((PORT + i))
  if port_in_use "$TRY_PORT"; then
    echo "    Port ${TRY_PORT} is busy, trying next port..."
    continue
  fi
  start_server "$TRY_PORT"

  # Wait for server to become healthy and process to stay alive
  for j in {1..40}; do
    if is_pid_alive "$SERVER_PID" && curl -sSf "http://${HOST}:${TRY_PORT}/health" >/dev/null; then
      BASE_URL="http://${HOST}:${TRY_PORT}"
      echo "    Server is up on ${HOST}:${TRY_PORT} (pid ${SERVER_PID})"
      break 2
    fi
    # If process died, break and try next port
    if ! is_pid_alive "$SERVER_PID"; then
      echo "    Server process died on port ${TRY_PORT}, trying next port..."
      break
    fi
    sleep 0.25
  done

  # Process failed to start; print last log lines and try next port
  echo "    Failed to start on ${HOST}:${TRY_PORT}, trying next port..."
  if [[ "$DEBUG_LOG" == "true" ]]; then
    tail -n 5 "$DEMO_LOG" || true
  fi
  sleep 0.2
done

if ! is_pid_alive "$SERVER_PID"; then
  echo "    Server failed to start on any port ${PORT}-$((PORT+ATTEMPTS-1))."
  if [[ "$DEBUG_LOG" == "true" ]]; then
    echo "    See ${DEMO_LOG} for details"
  fi
  echo "    You may need to kill hanging processes: lsof -i :${PORT}-$((PORT+ATTEMPTS-1))"
  exit 1
fi

cleanup() {
  printf "\n==> Stopping server (pid: %s)\n" "${SERVER_PID:-}"
  if [[ -n "${SERVER_PID:-}" ]] && is_pid_alive "$SERVER_PID"; then
    kill "$SERVER_PID" >/dev/null 2>&1 || true
  fi
  
  # Kill any remaining processes on demo ports to prevent hanging
  echo "==> Cleaning up any remaining processes on demo ports..."
  for i in $(seq 0 $((ATTEMPTS-1))); do
    TRY_PORT=$((PORT + i))
    if port_in_use "$TRY_PORT"; then
      echo "    Killing processes on port ${TRY_PORT}..."
      if command -v lsof >/dev/null 2>&1; then
        lsof -ti:"${TRY_PORT}" | xargs kill -9 >/dev/null 2>&1 || true
      fi
    fi
  done
  
  # Best-effort small delay to let sqlite close file handles
  sleep 0.2 || true
  if [[ -d "$DEMO_DIR" ]]; then
    echo "==> Cleaning up $DEMO_DIR"
    rm -rf "$DEMO_DIR" || true
  fi
}
trap cleanup EXIT INT TERM

echo "==> Waiting for server to become healthy..."
curl -sSf "${BASE_URL}/health" >/dev/null || { echo "    Server not healthy"; exit 1; }

echo
echo "==> Demo API calls"

echo "-- Create todo A"
curl -sS -X POST "${BASE_URL}/api/v1/todos/" \
  -H 'Content-Type: application/json' \
  -d '{"title":"A","description":"first"}' | pretty

echo "-- Create duplicate A (expect error)"
curl -sS -X POST "${BASE_URL}/api/v1/todos/" \
  -H 'Content-Type: application/json' \
  -d '{"title":"A","description":"dup"}' | pretty

echo "-- List todos"
curl -sS "${BASE_URL}/api/v1/todos/" | pretty

echo "-- Update todo id=1"
curl -sS -X PUT "${BASE_URL}/api/v1/todos/1" \
  -H 'Content-Type: application/json' \
  -d '{"title":"A","description":"updated","completed":true}' | pretty

echo "-- Get todo id=1"
curl -sS "${BASE_URL}/api/v1/todos/1" | pretty

echo "-- Delete todo id=1"
curl -sS -X DELETE "${BASE_URL}/api/v1/todos/1" | pretty

echo "-- List todos after delete"
curl -sS "${BASE_URL}/api/v1/todos/" | pretty

echo "-- Show all todos including soft-deleted (direct DB check)"
echo "    Note: This shows the actual database state including soft-deleted records"
sqlite3 "${DB_PATH}" "SELECT id, title, description, completed, deleted_at FROM todos;" 2>/dev/null || echo "    SQLite3 not available, but records are stored with soft delete"

echo "\n==> Demo complete"


