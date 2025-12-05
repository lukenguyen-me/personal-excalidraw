# Personal Excalidraw Backend

Go backend API for Personal Excalidraw - initial setup with middleware and infrastructure foundation.

## Features

- **Drawing CRUD API**: Complete REST API for managing drawings
  - Create, Read, Update, Delete, List operations
  - Pagination support for list endpoint
  - JSON data storage with PostgreSQL JSONB
- **Clean Architecture**: Domain-driven design with clear layer separation
- **HTTP server** with graceful shutdown
- **Structured logging** with slog
- **PostgreSQL database** with connection pooling (pgx/v5)
- **Database migrations** with SQL scripts
- **Middleware stack**:
  - CORS middleware
  - Request ID tracking
  - Recovery middleware for panic handling
  - Logger middleware for HTTP requests
- **Configuration management** from environment variables
- **Docker Compose setup** for development

## Tech Stack

- **Language**: Go 1.25+
- **Database**: PostgreSQL 16 with pgx/v5 driver
- **Router**: Standard library `net/http` with Go 1.22+ enhancements
- **Logging**: Standard library `log/slog`
- **Configuration**: godotenv for environment variables
- **Migrations**: golang-migrate/v4
- **Architecture**: Clean Architecture + Domain-Driven Design

## Quick Start

### Prerequisites

- Go 1.22+ (for local development)
- Docker & Docker Compose (recommended)

### Option 1: Docker Compose (Recommended)

From the project root:

```bash
# Start all services (frontend and backend)
docker-compose up

# Backend will be available at http://localhost:8080
# Frontend will be available at http://localhost:5173
```

### Option 2: Local Development

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Run the server**:
   ```bash
   make dev
   # Or: go run cmd/server/main.go
   ```

## Configuration

The backend uses environment variables for configuration with smart defaults.

### Default Settings

- **Server**: `0.0.0.0:8080`
- **CORS**: `http://localhost:5173` (frontend)
- **Logging**: Text format, info level

### Custom Configuration (Optional)

Create a `.env` file (copy from `.env.example`):

```bash
cp .env.example .env
```

Edit `.env` to customize settings:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173,https://mydomain.com

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
```

## Database Migrations

### Using the Migration Tool

The project uses `golang-migrate/v4` with embedded SQL files for database migrations.

**From the backend directory:**

```bash
# Run all pending migrations (most common)
go run migrations/migrate.go -up

# Check current migration version
go run migrations/migrate.go

# Rollback all migrations
go run migrations/migrate.go -down

# Apply one migration
go run migrations/migrate.go -steps 1

# Rollback one migration
go run migrations/migrate.go -steps -1

# Migrate to specific version
go run migrations/migrate.go -version 1

# Force set version (use with caution if migrations are in dirty state)
go run migrations/migrate.go -force 1
```

**In Docker environment:**

```bash
# Run migrations inside the backend container
docker exec personal-excalidraw-backend go run migrations/migrate.go -up
```

### Verify Migration

```bash
# Check if drawings table exists
docker exec personal-excalidraw-postgres psql -U postgres -d personal_excalidraw -c "\d drawings"
```

### Adding New Migrations

Create new SQL files in `migrations/` with sequential numbering:
```
migrations/
├── 001_create_drawings.sql
├── 002_add_user_tables.sql  # New migration
└── migrate.go
```

Each migration file should include both `Up` and `Down` sections:
```sql
-- +migrate Up
CREATE TABLE ...;

-- +migrate Down
DROP TABLE ...;
```

## API Endpoints

### Health Check
```http
GET /health
```

**Response** (200 OK):
```json
{
  "status": "ok"
}
```

### Drawing CRUD Operations

#### List Drawings
```http
GET /api/drawings?limit=10&offset=0
```

**Response** (200 OK):
```json
{
  "drawings": [
    {
      "id": "uuid",
      "name": "My Drawing",
      "data": {...},
      "created_at": "2025-12-05T10:30:00Z",
      "updated_at": "2025-12-05T10:30:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

#### Get Drawing
```http
GET /api/drawings/{id}
```

**Response** (200 OK):
```json
{
  "id": "uuid",
  "name": "My Drawing",
  "data": {...},
  "created_at": "2025-12-05T10:30:00Z",
  "updated_at": "2025-12-05T10:30:00Z"
}
```

#### Create Drawing
```http
POST /api/drawings
Content-Type: application/json

{
  "name": "New Drawing",
  "data": {
    "elements": [],
    "appState": {}
  }
}
```

**Response** (201 Created):
```json
{
  "id": "uuid",
  "name": "New Drawing",
  "data": {...},
  "created_at": "2025-12-05T10:30:00Z",
  "updated_at": "2025-12-05T10:30:00Z"
}
```

#### Update Drawing
```http
PUT /api/drawings/{id}
Content-Type: application/json

{
  "name": "Updated Drawing",
  "data": {
    "elements": [],
    "appState": {}
  }
}
```

**Response** (200 OK):
```json
{
  "id": "uuid",
  "name": "Updated Drawing",
  "data": {...},
  "created_at": "2025-12-05T10:30:00Z",
  "updated_at": "2025-12-05T10:35:00Z"
}
```

#### Delete Drawing
```http
DELETE /api/drawings/{id}
```

**Response** (204 No Content)

## Development

### Makefile Commands

```bash
make help              # Show all available commands
make dev               # Run development server
make build             # Build the application
make test              # Run tests
make test-coverage     # Run tests with coverage report
make docker-up         # Start Docker containers
make docker-down       # Stop Docker containers
make docker-logs       # View all logs
make lint              # Run linter
make fmt               # Format code
make tidy              # Tidy go modules
make install-tools     # Install development tools
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                            # Application entry point
├── internal/
│   ├── domain/                                # Domain layer (business logic)
│   │   └── drawing/
│   │       ├── drawing.go                     # Drawing entity
│   │       ├── errors.go                      # Domain errors
│   │       ├── repository.go                  # Repository interface
│   │       └── value_objects.go               # Value objects
│   ├── application/                           # Application layer (use cases)
│   │   └── drawing/
│   │       ├── dto.go                         # Data transfer objects
│   │       └── service.go                     # Application service
│   ├── adapter/                               # Interface adapters
│   │   ├── http/
│   │   │   ├── handler/                       # HTTP handlers
│   │   │   │   ├── drawing.go                 # Drawing CRUD handlers
│   │   │   │   ├── health.go                  # Health check handler
│   │   │   │   ├── response.go                # JSON response helpers
│   │   │   │   └── validation.go              # Request validation
│   │   │   ├── middleware/                    # HTTP middleware
│   │   │   │   ├── cors.go                    # CORS middleware
│   │   │   │   ├── logger.go                  # Logger middleware
│   │   │   │   ├── recover.go                 # Recovery middleware
│   │   │   │   └── request_id.go              # Request ID middleware
│   │   │   └── router.go                      # Route configuration
│   │   └── repository/
│   │       └── postgres/
│   │           ├── drawing_repository.go      # PostgreSQL repository
│   │           └── queries.go                 # SQL queries
│   └── infrastructure/                        # Infrastructure layer
│       ├── config/                            # Configuration management
│       │   └── config.go                      # Config loading from env
│       ├── database/                          # Database infrastructure
│       │   └── postgres.go                    # PostgreSQL connection
│       └── logger/                            # Logging infrastructure
│           └── logger.go                      # Structured logger setup
├── migrations/                                # Database migrations
│   ├── 001_create_drawings.sql                # Create drawings table
│   └── migrate.go                             # Migration tool
├── .env.example                               # Environment template
├── Makefile                                   # Build commands
├── go.mod                                     # Go dependencies
└── README.md                                  # This file
```

## Architecture

This backend follows **Clean Architecture** and **Domain-Driven Design (DDD)** principles:

### Domain Layer
- **Entities**: Rich domain models with business logic (Drawing)
- **Value Objects**: Immutable data structures (DrawingData)
- **Repository Interfaces**: Define data access contracts
- **Domain Errors**: Business rule violations
- No dependencies on external frameworks

### Application Layer
- **Use Cases/Services**: Orchestrate business operations
- **DTOs**: Data transfer objects for input/output
- **Transaction boundaries**: Service layer manages transactions
- Depend only on domain layer

### Adapter Layer
- **HTTP Handlers**: REST API endpoints for drawing CRUD operations
- **Repository Implementations**: PostgreSQL repository using pgx
- **Request Validation**: Input sanitization and validation
- Depend on application and domain layers

### Infrastructure Layer
- **Configuration**: Environment-based configuration with smart defaults
- **Database**: PostgreSQL connection pool management
- **Logging**: Structured logging with slog
- External dependencies (pgx, environment variables)

### HTTP Middleware Stack
- **Recovery**: Panic recovery with stack traces
- **Request ID**: Request tracking (X-Request-ID header)
- **Logger**: HTTP request/response logging
- **CORS**: Cross-origin support

### Benefits
- **Clean Separation**: Each layer has clear responsibilities
- **Testable**: Easy to mock dependencies for unit tests
- **Maintainable**: Changes in one layer don't affect others
- **Domain-Focused**: Business logic is independent of frameworks
- **Flexible**: Easy to swap implementations (e.g., different database)
- **Production-ready**: Graceful shutdown, error handling, structured logging

## Logging

Structured logging with `log/slog`:

```
2025-01-15 10:30:45 INFO HTTP request method=GET path=/api/drawings status=200 duration_ms=45
2025-01-15 10:30:46 ERROR Database error error="connection refused"
```

**Production** (JSON format):
```json
{"time":"2025-01-15T10:30:45Z","level":"INFO","msg":"HTTP request","method":"GET","path":"/api/drawings","status":200}
```

## Contributing

1. Follow clean architecture principles
2. Write tests for new features
3. Use `make fmt` before committing
4. Run `make lint` to check code quality
5. Update documentation

## License

MIT License - See LICENSE file for details
