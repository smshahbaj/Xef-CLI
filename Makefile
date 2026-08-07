.PHONY: build test test-race lint clean install run

BINARY_NAME := xef
BUILD_DIR := bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION) -s -w"

build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/xefcli

install:
	go install $(LDFLAGS) ./cmd/xefcli

test:
	CGO_ENABLED=0 go test -p 1 -count=1 -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-race:
	CGO_ENABLED=0 go test -p 1 -count=1 -v -race ./...

test-ci:
	CGO_ENABLED=0 go test -p 1 -count=1 -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out dist/

run:
	go run ./cmd/xefcli

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean

deps:
	go mod download
	go mod tidy

generate:
	go generate ./...

.DEFAULT_GOAL := build
