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

## Quick Start

### Prerequisites

- Node.js 18+ (for frontend development)
- pnpm (package manager)
- Go 1.24+ (for backend development)
- PostgreSQL 14+ (for database)
- Docker & Docker Compose (optional, for containerized deployment)

### Development

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
docker run -p 8080:8080 personal-excalidraw:latest
```

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
