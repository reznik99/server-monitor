BINARY_NAME := server-monitor
VERSION := $(shell git describe --tags 2>/dev/null || echo "dev")
LDFLAGS := -ldflags="-X 'main.Version=$(VERSION)'"
DIST_DIR := dist

.PHONY: all lint test build build-amd64 build-arm64 clean

all: lint test build

lint:
	@echo "==> Linting..."
	golangci-lint run ./...

test:
	@echo "==> Running tests..."
	# TODO: add tests
	go test ./...

build: build-amd64 build-arm64

build-amd64:
	@echo "==> Building $(BINARY_NAME) (linux/amd64)..."
	GOOS=linux GOARCH=amd64 go build -v $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64

build-arm64:
	@echo "==> Building $(BINARY_NAME) (linux/arm64)..."
	GOOS=linux GOARCH=arm64 go build -v $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64

clean:
	@echo "==> Cleaning..."
	rm -rf $(DIST_DIR)/$(BINARY_NAME)-*
