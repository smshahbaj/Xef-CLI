#!/bin/bash
set -euo pipefail

export CGO_ENABLED=0

echo "Running tests..."
go test -p 1 -count=1 -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
echo "Tests complete"
