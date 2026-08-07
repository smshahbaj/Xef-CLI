#!/bin/bash
set -euo pipefail

echo "Running tests..."
go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
echo "Tests complete"
