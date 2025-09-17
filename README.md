# Let It Go

A 3-tier Go application built with Echo framework, MySQL database, and structured logging using slog.

## Architecture

This application follows a feature-based architecture pattern:


### Feature Structure
- **Repository Package**: `feature/<feature_name>/repository/`
  - `entity.go` - Domain models
  - `repository.go` - Repository struct and constructor
  - `repository_interface.go` - Repository interface
  - `count.go` - Count method implementation
  - `*_test.go` - Individual test files for each method
  - Individual service methods: `list.go`, `count.go`, etc.
- **Service Package**: `feature/<feature_name>/service/`
  - `service.go` - Service implementation
  - `service_interface.go` - Service interface
  - `*_dto.go` - Action-specific DTOs with conversion functions
  - Individual service methods: `create.go`, `get.go`, etc.
- **Handler Package**: `feature/<feature_name>/handler/`
  - `http_handler.go` - HTTP handlers with routing and Swagger annotations
  - `http_handler_v2.go` - Version 2 HTTP handlers (if using API versioning)
  - `http_response.go` - Response utilities and types


## Project Structure

```
project-root/
├── cmd/
│   ├── http_server/     # HTTP server entry point
│   │   └── main.go      # Main application with Swagger annotations
│   └── migrate/         # Database migration tool
│       └── main.go      # Migration commands (up, down, force, version, create)
├── config/              # Configuration management
│   └── config.go        # Config loading with .env file support
├── feature/             # Feature-based modules
│   └── <feature_name>/  # Example feature
│       ├── repository/  # Data access layer
│       │   ├── entity.go
│       │   ├── repository.go
│       │   ├── repository_interface.go
│       │   ├── <method>.go
│       │   └── *_test.go
│       ├── service/     # Business logic layer
│       │   ├── service.go
│       │   ├── service_interface.go
│       │   ├── <method>.go
│       │   └── *_dto.go
│       └── handler/     # Presentation layer
│           ├── http_handler.go      # Main HTTP handlers with Swagger docs
│           ├── http_handler_v2.go   # Version 2 handlers
│           └── http_response.go     # Response utilities and types
├── pkg/                 # Shared packages
│   ├── database/        # Database connection management
│   ├── logger/          # Structured logging utilities
│   └── http_server/     # Generic HTTP server with Swagger middleware
├── migrations/          # Database schema migrations (timestamp format)
│   ├── {timestamp}_{name}.up.sql    # Up migrations
│   └── {timestamp}_{name}.down.sql  # Down migrations
├── docs/                # Generated Swagger documentation
└── Makefile            # Build and development commands
```


## Prerequisites

- Go 1.25 or higher
- Make (optional, for build automation)

## Quick Start

### Local Development

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd let-it-go
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Setup environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your database connection details
   ```

4. **Setup database and run migrations**:
   ```bash
   make migrate-up
   ```

5. **Generate API documentation**:
   ```bash
   make swagger
   ```

6. **Start the server**:
   ```bash
   make run
   ```

## Testing

```bash
# Run all tests
make test

# Run unit tests only (short tests)
make test-unit

# Skip integration tests (uses SKIP_INTEGRATION_TESTS environment variable)
SKIP_INTEGRATION_TESTS=true go test ./...
```

## Running the Application

### With Docker

```bash
docker-compose up --build
```

### Local Build

**Development:**
```bash
make run
```

**Production Build:**
```bash
make build-production
./bin/server
```


