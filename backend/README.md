# Personal Excalidraw Backend

Go backend API for Personal Excalidraw - initial setup with middleware and infrastructure foundation.

## Features

- HTTP server with graceful shutdown
- Structured logging with slog
- CORS middleware
- Request ID tracking
- Recovery middleware for panic handling
- Logger middleware for HTTP requests
- Configuration management from environment variables
- Docker Compose setup

## Tech Stack

- **Language**: Go 1.22+
- **Router**: Standard library `net/http` with Go 1.22+ enhancements
- **Logging**: Standard library `log/slog`
- **Configuration**: godotenv for environment variables

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
│       └── main.go                    # Application entry point
├── internal/
│   ├── adapter/                       # Interface adapters
│   │   └── http/
│   │       ├── handler/               # HTTP handlers
│   │       │   ├── health.go          # Health check handler
│   │       │   └── response.go        # JSON response helpers
│   │       ├── middleware/            # HTTP middleware
│   │       │   ├── cors.go            # CORS middleware
│   │       │   ├── logger.go          # Logger middleware
│   │       │   ├── recover.go         # Recovery middleware
│   │       │   └── request_id.go      # Request ID middleware
│   │       └── router.go              # Route configuration
│   └── infrastructure/                # Infrastructure layer
│       ├── config/                    # Configuration management
│       │   └── config.go              # Config loading from env
│       └── logger/                    # Logging infrastructure
│           └── logger.go              # Structured logger setup
├── .env.example                       # Environment template
├── Makefile                           # Build commands
├── go.mod                             # Go dependencies
└── README.md                          # This file
```

## Architecture

This backend provides a clean foundation with:

### Infrastructure Layer
- Configuration management with smart defaults
- Structured logging with slog
- Environment variable support

### HTTP Layer
- Standard library router with Go 1.22+ enhancements
- Middleware stack:
  - Recovery: Panic recovery
  - Request ID: Request tracking
  - Logger: HTTP request logging
  - CORS: Cross-origin support
- Health check endpoint

### Benefits
- **Simple**: Minimal dependencies, standard library first
- **Maintainable**: Clear separation of concerns
- **Flexible**: Ready to add domain logic, database, etc.
- **Production-ready**: Graceful shutdown, structured logging

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
