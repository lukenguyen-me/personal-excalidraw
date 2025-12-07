# Personal Excalidraw

A self-hosted, open-source drawing application built on Excalidraw. Create and manage your diagrams with a clean, modern interface.

## Why Personal Excalidraw?

Personal Excalidraw gives you a self-hosted alternative to cloud-based drawing tools. Your drawings stay on your machine (for now) and will sync to your own server when the backend is ready. It's built with modern web technologies and designed for easy deployment.

**Use Cases:**
- Architecture diagrams
- Wireframes and mockups
- Mind maps and brainstorming
- Technical documentation
- Teaching and presentations

## Key Features

- Free drawing application based on Excalidraw
- Auto-save with 1-second debounce
- Full CRUD operations (Create, Read, Update, Delete) for drawings
- Inline name editing with keyboard shortcuts (Enter to save, Escape to cancel)
- API-backed data storage with Go backend
- TanStack Query for efficient data fetching and caching
- Modern UI with SvelteKit and Tailwind CSS
- Responsive design with clean table layout
- Clean Architecture backend with comprehensive validation
- Access key authentication for secure self-hosting

## Quick Start

### Prerequisites

- Node.js 18+ (for frontend development)
- pnpm (package manager)
- Go 1.24+ (for backend development)
- PostgreSQL 14+ (for database)
- Docker & Docker Compose (optional, for containerized deployment)

### Development

#### Configuration

Before running the application, configure your environment:

```bash
# Backend configuration
cd backend
cp .env.example .env

# Edit .env and set your access key
# AUTH_ENABLED=true
# ACCESS_KEY=your-secret-key-here
```

#### Full Stack Development

```bash
# Terminal 1 - Start Backend
cd backend
go run cmd/server/main.go

# Terminal 2 - Start Frontend
cd frontend
pnpm install
pnpm run dev
```

This starts:
- Backend API server on `http://localhost:8080`
- Frontend dev server on `http://localhost:5173`

When you first access the application, you'll be prompted to enter your access key if authentication is enabled.

#### With Docker Compose

```bash
docker-compose up
```

This starts:
- PostgreSQL database on port `5432`
- Backend server on `http://localhost:8080`
- Frontend dev server on `http://localhost:5173`

### Building for Production

```bash
docker build -t personal-excalidraw:latest .
docker run -p 8080:8080 \
  -e ACCESS_KEY=your-secret-key-here \
  -e AUTH_ENABLED=true \
  personal-excalidraw:latest
```

## Authentication

Personal Excalidraw includes access key authentication to secure your self-hosted instance.

### Configuration

Set these environment variables in your `.env` file or pass them to Docker:

```bash
# Enable/disable authentication
AUTH_ENABLED=true

# Your secret access key (change this!)
ACCESS_KEY=your-secret-key-here
```

### Usage

1. When authentication is enabled, you'll see an authentication modal on first visit
2. Enter your access key to gain access to the application
3. The key is stored in localStorage and persists across browser sessions
4. The key is synchronized across tabs for convenience

### Security Notes

- Change the default access key before deploying to production
- Use a strong, randomly generated key for better security
- The access key is validated using constant-time comparison to prevent timing attacks
- Authentication can be disabled by setting `AUTH_ENABLED=false` for local development

## Documentation

- [ARCHITECTURE.md](ARCHITECTURE.md) - Technical architecture, tech stack, and implementation details
- [PLAN.md](PLAN.md) - Development roadmap and feature planning

## Contributing

This is an open-source project. Contributions are welcome!

## License

MIT License - See LICENSE file for details

## Acknowledgments

- Built on top of [Excalidraw](https://excalidraw.com/)
- Uses [SvelteKit](https://svelte.dev/) for the frontend framework
- Styled with [Tailwind CSS](https://tailwindcss.com/) and [DaisyUI](https://daisyui.com/)
