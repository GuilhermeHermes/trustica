.PHONY: certs build test shell server clean help

# Docker compose file
COMPOSE=docker compose -f docker-compose.yml

## help: Show this help message
help:
	@echo "Trustica Integration Tests"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'

## certs: Generate test certificates (root CA + server cert)
certs:
	@echo "Generating test certificates..."
	@chmod +x generate_certs.sh
	@./generate_certs.sh

## build: Build test containers
build: certs
	@echo "Building test containers..."
	@COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 $(COMPOSE) build

## test: Run full integration test suite
test: build
	@echo "Running integration tests..."
	@$(COMPOSE) up --abort-on-container-exit --exit-code-from testrunner
	@$(COMPOSE) down

## shell: Start a shell in the test container (for debugging)
shell: build
	@$(COMPOSE) run --rm testrunner /bin/bash

## server: Start just the test server locally (for debugging)
server: certs
	@echo "Starting test server on https://localhost:8443..."
	@cd testserver && TLS_CERT=../certs/server.pem TLS_KEY=../certs/server-key.pem go run main.go

## clean: Clean up containers, images, and certificates
clean:
	@$(COMPOSE) down --rmi local --volumes 2>/dev/null || true
	@rm -rf certs/

## logs: Show logs from last test run
logs:
	@$(COMPOSE) logs

## status: Show status of test containers
status:
	@$(COMPOSE) ps
