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
- Local-first data storage (localStorage)
- Modern UI with SvelteKit and DaisyUI wireframe theme
- Responsive design with clean table layout
- Future: Go backend for cloud synchronization

## Quick Start

### Prerequisites

- Node.js 18+ (for frontend development)
- pnpm (package manager)
- Docker & Docker Compose (optional, for containerized deployment)

### Development

#### Frontend Only

```bash
cd frontend
pnpm install
pnpm run dev
```

The frontend will be available at `http://localhost:5173`

#### With Docker Compose

```bash
docker-compose up
```

This starts:
- Frontend dev server on `http://localhost:5173`
- Backend server on `http://localhost:8080` (stub, for future use)

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
