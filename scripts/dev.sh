#!/bin/bash
# FileENIAC - Development Script

set -e

MODE="${1:-dev}"

case "$MODE" in
  dev)
    echo "Starting backend in dev mode..."
    cd "$(dirname "$0")/../backend"
    go run . --dev
    ;;
  build)
    echo "Building backend..."
    cd "$(dirname "$0")/.."
    make build
    ;;
  desktop)
    echo "Starting desktop app..."
    cd "$(dirname "$0")/../apps/desktop"
    npm run tauri dev
    ;;
  test)
    echo "Running tests..."
    cd "$(dirname "$0")/../backend"
    go test ./... -v
    ;;
  *)
    echo "Usage: $0 {dev|build|desktop|test}"
    exit 1
    ;;
esac
