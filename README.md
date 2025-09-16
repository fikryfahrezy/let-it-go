# Let It Go - 3-Tier Go Application

A future-proof, object-oriented 3-tier Go application built with Echo framework, MySQL database, and structured logging using slog.

## Architecture

This application follows a feature-based architecture pattern:

### Feature-Based Organization
- **Location**: `feature/<feature_name>/`
- **Structure**: Each feature contains all related components in a single package
- **Components**: Entities, repositories, services, and HTTP handlers grouped by feature

### Feature Structure (User Example)
- **Repository Package**: `feature/user/repository/`
  - `entity.go` - Domain models (User entity)
  - `repository.go` - Data access implementation
  - `repository_interface.go` - Repository interface
- **Service Package**: `feature/user/service/`
  - `service.go` - Service implementation
  - `service_interface.go` - Service interface
  - `*_dto.go` - Action-specific DTOs with conversion functions
  - Individual service methods: `create_user.go`, `get_user.go`, etc.
- **Handler Package**: `feature/user/handler/`
  - `http_handler.go` - Main HTTP handlers with routing
  - `http_handler_v2.go` - Version 2 specific endpoints
  - `http_response.go` - HTTP response utilities

### API Layer Architecture
- **Versioned APIs**: Each API version (v1, v2) has its own package
- **HTTP/REST**: RESTful API implementation using Echo framework
- **Interface-Based**: Common interface for all API server implementations
- **Server Management**: Centralized server lifecycle management in `pkg/server`

## Project Structure

```
let-it-go/
├── cmd/server/          # Application entry point
├── config/              # Configuration management
├── feature/             # Feature-based modules
│   ├── health/          # Health check feature
│   │   └── handler/     # Health HTTP handlers
│   │       └── http_handler.go
│   └── user/            # User management feature
│       ├── repository/  # Data access layer
│       │   ├── entity.go
│       │   ├── repository.go
│       │   └── repository_interface.go
│       ├── service/     # Business logic layer
│       │   ├── service.go
│       │   ├── service_interface.go
│       │   ├── create_user.go
│       │   ├── get_user.go
│       │   ├── update_user.go
│       │   ├── delete_user.go
│       │   ├── list_users.go
│       │   └── *_dto.go
│       └── handler/     # HTTP presentation layer
│           ├── http_handler.go
│           ├── http_handler_v2.go
│           └── http_response.go
├── pkg/                 # Shared packages
│   ├── database/        # Database connection management
│   ├── logger/          # Structured logging utilities
│   └── server/          # Generic HTTP server
├── migrations/          # Database schema migrations
└── Makefile            # Build and development commands
```

## Features

- **Layered Architecture**: Repository, Service, and Handler layers with clear separation
- **Feature-Based Organization**: Each feature has its own repository, service, and handler packages  
- **Action-Specific DTOs**: Separate DTOs for create, update, get, and list operations
- **API Versioning**: Built-in support for v1, v2, and future API versions
- **Health Monitoring**: Dedicated health check feature with database connectivity testing
- **Type-Safe Configuration**: Enum-based configuration for logger levels and formats
- **Structured Logging**: Using Go's built-in `slog` package with contextual logging
- **Database Integration**: MySQL with DSN-based configuration and connection pooling
- **RESTful API**: Echo framework with comprehensive middleware support
- **Graceful Shutdown**: Proper resource cleanup for server and database connections
- **Interface-Based Design**: Dependency inversion with repository and service interfaces
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Development Tools**: Makefile with build, test, lint, and health check commands

## Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher
- Make (optional, for build automation)

## Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd let-it-go
   ```

2. **Install dependencies**:
   ```bash
   make deps
   # or
   go mod download
   ```

3. **Setup environment**:
   ```bash
   make dev-setup
   # or manually
   cp .env.example .env
   ```

4. **Configure database** (edit `.env` file):
   ```env
   DB_DSN=root:password@tcp(localhost:3306)/letitgo?charset=utf8mb4&parseTime=True&loc=Local
   ```

5. **Create database and run migrations**:
   ```bash
   make db-create
   make db-migrate
   ```

## Running the Application

### Development
```bash
make run
# or
go run cmd/server/main.go
```

### Production Build
```bash
make build
./bin/server
```

## API Endpoints

### Health Check
- **GET** `/api/health` - Service health check with database connectivity test

### Users (Version 1)
- **POST** `/api/v1/users` - Create a new user
- **GET** `/api/v1/users` - List users (with pagination)
- **GET** `/api/v1/users/:id` - Get user by ID
- **PUT** `/api/v1/users/:id` - Update user
- **DELETE** `/api/v1/users/:id` - Delete user

### Users (Version 2 - Enhanced)
- **POST** `/api/v2/users` - Create a new user
- **GET** `/api/v2/users` - List users (with pagination)
- **GET** `/api/v2/users/:id` - Get user by ID
- **PUT** `/api/v2/users/:id` - Update user
- **DELETE** `/api/v2/users/:id` - Delete user
- **GET** `/api/v2/users/:id/profile` - Get enhanced user profile
- **POST** `/api/v2/users/batch` - Batch user operations

### Example API Usage

**Health Check:**
```bash
curl http://localhost:8080/api/health
```

**Create User:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

**Get Users:**
```bash
curl http://localhost:8080/api/v1/users?page=1&page_size=10
```

**Check Service Health via Makefile:**
```bash
make health
```

## Configuration

### Database Connection

The application uses a complete DSN (Data Source Name) for database configuration:

```env
DB_DSN=username:password@tcp(host:port)/database?param1=value1&param2=value2
```

### Common DSN Examples

**Local Development:**
```
root:password@tcp(localhost:3306)/letitgo?charset=utf8mb4&parseTime=True&loc=Local
```

**Production with SSL:**
```
user:pass@tcp(prod-host:3306)/letitgo?charset=utf8mb4&parseTime=True&tls=true
```

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | Server host | `localhost` |
| `SERVER_PORT` | Server port | `8080` |
| `DB_DSN` | Complete database connection string | `` |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |
| `LOG_FORMAT` | Log format (text, json) | `text` |

## Development

### Code Formatting
```bash
make fmt
```

### Linting
```bash
make lint
```

### Testing
```bash
make test
```

### Health Check
```bash
make health
```

### Clean Build Artifacts
```bash
make clean
```

## Available Make Commands

Run `make help` to see all available commands:

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make run` | Run the application in development mode |
| `make test` | Run all tests |
| `make clean` | Clean build artifacts |
| `make fmt` | Format code using go fmt |
| `make lint` | Run linter (requires golangci-lint) |
| `make tidy` | Tidy go modules |
| `make deps` | Download dependencies |
| `make dev-setup` | Setup development environment |
| `make health` | Check service health status |
| `make db-create` | Create database |
| `make db-migrate` | Run database migrations |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |

## Key Design Patterns

1. **Feature-Based Architecture**: Related functionality grouped in feature packages
2. **Layered Architecture**: Clear separation between Repository, Service, and Handler layers
3. **Dependency Injection**: Services depend on interfaces, not concrete implementations
4. **Repository Pattern**: Abstracts data access layer with interface-based contracts
5. **Service Layer Pattern**: Encapsulates business logic with action-specific methods
6. **DTO Pattern**: Separate DTOs for different operations (Create, Update, Get, List)
7. **Interface Segregation**: Small, focused interfaces for better maintainability
8. **Conversion Functions**: Dedicated functions for entity-to-DTO transformations
9. **Route Handler Interface**: Generic interface for feature route registration
10. **Single Responsibility**: Each layer and file has a focused, single purpose

## Future Enhancements

- [ ] Add authentication and authorization
- [ ] Implement caching layer
- [ ] Add metrics and monitoring
- [ ] Add comprehensive testing
- [ ] Add API documentation (Swagger/OpenAPI)
- [ ] Add Docker support
- [ ] Add CI/CD pipeline

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run linting and tests
6. Submit a pull request

## License

This project is licensed under the MIT License.