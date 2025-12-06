.PHONY: help dev build test docker-up docker-down lint fmt tidy migrate-up migrate-down migrate-status migrate-force

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Run development server
	go run cmd/server/main.go

build: ## Build the application
	go build -o bin/server cmd/server/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

docker-up: ## Start Docker containers
	cd .. && docker-compose up -d

docker-down: ## Stop Docker containers
	cd .. && docker-compose down

docker-logs: ## View Docker logs
	cd .. && docker-compose logs -f

docker-logs-backend: ## View backend logs only
	cd .. && docker-compose logs -f backend

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

tidy: ## Tidy go modules
	go mod tidy

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

migrate-up: ## Run all pending migrations inside Docker
	docker exec personal-excalidraw-backend sh -c "cd migrations && go run migrate.go -up"

migrate-down: ## Rollback all migrations inside Docker
	docker exec personal-excalidraw-backend sh -c "cd migrations && go run migrate.go -down"

migrate-status: ## Check current migration version inside Docker
	docker exec personal-excalidraw-backend sh -c "cd migrations && go run migrate.go"

migrate-force: ## Force migration version (usage: make migrate-force VERSION=1)
	docker exec personal-excalidraw-backend sh -c "cd migrations && go run migrate.go -force $(VERSION)"
