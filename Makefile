.PHONY: build run test lint clean fmt vet help

# Binary name
BINARY_NAME=trustica
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'

## build: Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/trustica

## run: Run the application
run:
	$(GORUN) ./cmd/trustica

## test: Run tests
test:
	$(GOTEST) -v -race ./...

## test-cover: Run tests with coverage
test-cover:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	$(GOFMT) ./...

## vet: Run go vet
vet:
	$(GOVET) ./...

## tidy: Tidy and verify dependencies
tidy:
	$(GOMOD) tidy
	$(GOMOD) verify

## clean: Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## all: Format, vet, lint, test, and build
all: fmt vet lint test build

# =============================================================================
# Integration Tests (Docker-based)
# =============================================================================

## integration-certs: Generate test certificates
integration-certs:
	@echo "Generating test certificates..."
	@chmod +x testcontainer/generate_certs.sh
	@./testcontainer/generate_certs.sh

## integration-build: Build test containers
integration-build: build integration-certs
	@echo "Building test containers..."
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker compose -f testcontainer/docker-compose.yml build

## integration-test: Run full integration test suite
integration-test: integration-build
	@echo "Running integration tests..."
	@docker compose -f testcontainer/docker-compose.yml up --abort-on-container-exit --exit-code-from testrunner
	@docker compose -f testcontainer/docker-compose.yml down

## integration-shell: Start a shell in the test container (for debugging)
integration-shell: integration-build
	@docker compose -f testcontainer/docker-compose.yml run --rm testrunner /bin/bash

## integration-server: Start just the test server (for local testing)
integration-server: integration-certs
	@echo "Starting test server on https://localhost:8443..."
	@cd testcontainer/testserver && TLS_CERT=../certs/server.pem TLS_KEY=../certs/server-key.pem go run main.go

## integration-clean: Clean up test containers and certificates
integration-clean:
	@docker compose -f testcontainer/docker-compose.yml down --rmi local --volumes 2>/dev/null || true
	@rm -rf testcontainer/certs/
