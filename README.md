# 3-Tier Go Application

A future-proof, object-oriented 3-tier Go application built with Echo framework, MySQL database, and structured logging using slog.

## Architecture

This application follows a feature-based architecture pattern:


### Feature Structure
- **Repository Package**: `feature/<feature_name>/repository/`
  - `entity.go` - Domain models
  - `repository.go` - Repository struct and constructor
  - `repository_interface.go` - Repository interface
  - `create.go` - Create method implementation
  - `get_by_id.go` - Get by ID method implementation
  - `get_by_email.go` - Get by email method implementation
  - `update.go` - Update method implementation
  - `delete.go` - Delete method implementation
  - `list.go` - List method implementation
  - `count.go` - Count method implementation
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
│       │   ├── create.go
│       │   ├── get_by_id.go
│       │   ├── get_by_email.go
│       │   ├── update.go
│       │   ├── delete.go
│       │   ├── list.go
│       │   └── count.go
│       ├── service/     # Business logic layer
│       │   ├── service.go
│       │   ├── service_interface.go
│       │   ├── create.go
│       │   ├── get.go
│       │   ├── update.go
│       │   ├── delete.go
│       │   ├── list.go
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

## API Documentation

This application includes interactive API documentation using Swagger/OpenAPI:

- **Generate docs**: `make swagger`
- **View docs**: Start the server and visit `http://localhost:8080/swagger/`
- **API Base**: All endpoints are under `/api` path

## Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd <project-name>
   ```

2. **Install dependencies**:
   ```bash
   make deps
   # or
   go mod download
   ```

3. **Setup environment** (optional):
   ```bash
   cp .env.example .env
   ```

4. **Configure environment variables**: Edit `.env` file with your specific values for database connection, server settings, and logging preferences.

5. **Setup database and run migrations**:
   ```bash
   # Create database (if using MySQL)
   make db-create
   
   # Run migrations using custom migration tool
   make migrate-up
   
   # Check migration version
   make migrate-version
   ```

6. **Generate API documentation**:
   ```bash
   make swagger
   ```

## Running the Application

### Development
```bash
make run
# or
go run cmd/http_server/main.go
```

### Production Build
```bash
make build
./bin/server
```


