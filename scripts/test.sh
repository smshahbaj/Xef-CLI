#!/bin/bash
set -e

echo "Running tests..."
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
echo "Tests complete"
