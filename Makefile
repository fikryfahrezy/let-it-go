.PHONY: build run test clean fmt lint tidy deps health db-create db-migrate migrate-up migrate-down migrate-force migrate-version migrate-create swagger docker-build docker-run help

# Build the application
build:
	go build -ldflags="-X 'main.version=$(shell git describe --tags --always --dirty)' -X 'main.commit=$(shell git rev-parse HEAD)' -X 'main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" -o bin/server cmd/http_server/main.go

# Run the application
run:
	go run cmd/http_server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Download dependencies
deps:
	go mod download


# Health check
health:
	curl -f http://localhost:8080/api/health || echo "Service is not running or unhealthy"

# Database operations
db-create:
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS letitgo;"

# Migration operations using custom migrate command
migrate-up:
	go run cmd/migrate/main.go -action=up

migrate-down:
	go run cmd/migrate/main.go -action=down

migrate-force:
	go run cmd/migrate/main.go -action=force -version=$(VERSION)

migrate-version:
	go run cmd/migrate/main.go -action=version

migrate-create:
	go run cmd/migrate/main.go -action=create -name=$(NAME)

# Legacy migration (deprecated)
db-migrate:
	mysql -u root -p letitgo < migrations/001_create_users_table.sql

# Swagger documentation generation
swagger:
	go run github.com/swaggo/swag/cmd/swag init -g cmd/http_server/main.go -o docs

# Docker operations
docker-build:
	docker build -t let-it-go .

docker-run:
	docker run -p 8080:8080 let-it-go

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  deps          - Download dependencies"
	@echo "  health        - Check service health"
	@echo "  db-create     - Create database"
	@echo "  db-migrate    - Run migrations (legacy)"
	@echo "  migrate-up    - Run all migrations up"
	@echo "  migrate-down  - Run all migrations down"
	@echo "  migrate-force - Force migration version (use VERSION=N)"
	@echo "  migrate-version - Show current migration version"
	@echo "  migrate-create - Create new migration (use NAME=migration_name)"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"