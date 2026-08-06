#!/bin/bash
set -e

echo "Building XefCLI..."
mkdir -p bin
go build -ldflags "-s -w" -o bin/xef ./cmd/xefcli
echo "Build complete: bin/xef"
