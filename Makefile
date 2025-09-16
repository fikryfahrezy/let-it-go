.PHONY: build run test clean fmt lint tidy deps dev-setup health db-create db-migrate docker-build docker-run help

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

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

# Development setup
dev-setup:
	cp .env.example .env
	go mod download

# Health check
health:
	curl -f http://localhost:8080/api/health || echo "Service is not running or unhealthy"

# Database operations
db-create:
	mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS letitgo;"

db-migrate:
	mysql -u root -p letitgo < migrations/001_create_users_table.sql

# Docker operations
docker-build:
	docker build -t let-it-go .

docker-run:
	docker run -p 8080:8080 let-it-go

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  tidy        - Tidy dependencies"
	@echo "  deps        - Download dependencies"
	@echo "  dev-setup   - Setup development environment"
	@echo "  health      - Check service health"
	@echo "  db-create   - Create database"
	@echo "  db-migrate  - Run migrations"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"