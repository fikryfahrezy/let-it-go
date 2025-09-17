.PHONY: build build-production run test test-unit test-clean-cache clean fmt health migrate-up migrate-down migrate-force migrate-version migrate-create swagger generate cron

# Build the application for development
build:
	go build -ldflags="-X 'main.version=$(shell git describe --tags --always --dirty)' -X 'main.commit=$(shell git rev-parse HEAD)' -X 'main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" -o bin/server ./cmd/http_server

# Build the application for production
build-production:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -X 'main.version=$(shell git describe --tags --always --dirty)' -X 'main.commit=$(shell git rev-parse HEAD)' -X 'main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" -a -installsuffix cgo -o bin/server ./cmd/http_server

# Run the application
run:
	go run cmd/http_server/main.go

# Run tests
test:
	go test -v ./...

# Run unit tests only
test-unit:
	go test -v ./... -short

# Clear test cache
test-clean-cache:
	go clean -testcache

# Run unit tests with coverage
test-unit-coverage:
	go test -v -coverprofile=coverage.out -run "Unit" ./...
	go tool cover -func=coverage.out

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...

# Run tests with coverage and generate HTML report
test-coverage-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Show coverage summary
test-coverage-summary:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Health check
health:
	curl -f http://localhost:8080/api/health || echo "Service is not running or unhealthy"

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

# Swagger documentation generation
swagger:
	go tool swag init -g cmd/http_server/main.go -o docs

# Generate code 
generate:
	go generate ./...

# Run scheduled tasks
cron:
	go run cmd/cron/main.go

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application for development"
	@echo "  build-production - Build the application for production (optimized)"
	@echo "  run           - Run the application"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-clean-cache - Clear test cache"
	@echo "  test-unit-coverage - Run unit tests with coverage"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-coverage-html - Run tests with coverage and generate HTML report"
	@echo "  test-coverage-summary - Run tests with coverage and show summary"
	@echo "  clean         - Clean build artifacts"
	@echo "  fmt           - Format code"
	@echo "  health        - Check service health"
	@echo "  migrate-up    - Run all migrations up"
	@echo "  migrate-down  - Run all migrations down"
	@echo "  migrate-force - Force migration version (use VERSION=N)"
	@echo "  migrate-version - Show current migration version"
	@echo "  migrate-create - Create new migration (use NAME=migration_name)"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  generate      - Run code generation"
	@echo "  cron          - Run scheduled tasks"