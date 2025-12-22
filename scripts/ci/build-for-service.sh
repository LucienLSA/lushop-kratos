#!/usr/bin/env bash
set -euo pipefail

if [ -z "${1:-}" ]; then
  echo "Usage: $0 <service_dir>"
  exit 1
fi

SERVICE_DIR="$1"
SERVICE_NAME="$(basename "$SERVICE_DIR")"

echo "Building service '$SERVICE_NAME' at '$SERVICE_DIR'"

cd "$SERVICE_DIR"

echo "Checking gofmt..."
gofmt -l .

echo "Running go vet..."
go vet ./...

if command -v golangci-lint >/dev/null 2>&1; then
  echo "Running golangci-lint..."
  golangci-lint run ./... || true
else
  echo "golangci-lint not found, skipping lint step."
fi

echo "Running tests..."
go test ./...

echo "Building binary..."
mkdir -p ../build
go build -o "../build/${SERVICE_NAME}" ./...

echo "Built ../build/${SERVICE_NAME}"

cd -


