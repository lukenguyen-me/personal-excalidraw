# Personal Excalidraw

A self-hosted, open-source drawing application built on Excalidraw. Store your diagrams locally with no cloud dependency.

## Features

- Free drawing application based on Excalidraw
- Local-first data storage (IndexedDB/localStorage)
- Self-hosted deployment with Docker
- Modern UI with Svelte and Daisy UI
- Go backend for future cloud synchronization features

## Project Structure

```
personal-excalidraw/
├── frontend/          # Svelte + TypeScript frontend
├── backend/           # Go backend (stub)
├── docker-compose.yml # Development environment
└── Dockerfile         # Production build
```

## Quick Start

### Prerequisites

- Node.js 18+ (for frontend development)
- Go 1.21+ (for backend development)
- Docker & Docker Compose (for containerized deployment)

### Development

#### Frontend Only

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`

#### With Docker Compose

```bash
docker-compose up
```

This starts:
- Frontend dev server on `http://localhost:5173`
- Backend server on `http://localhost:8080`

### Building for Production

```bash
docker build -t personal-excalidraw:latest .
docker run -p 8080:8080 personal-excalidraw:latest
```

## Tech Stack

- **Frontend**: Svelte + TypeScript + Vite
- **Styling**: Tailwind CSS + Daisy UI
- **Drawing**: @excalidraw/excalidraw (React wrapper)
- **State Management**: Svelte Stores
- **Backend**: Go (net/http)
- **Database**: PostgreSQL (future)
- **Deployment**: Docker + Docker Compose

## Development Roadmap

### Phase 1: Frontend Infrastructure ✅
- [x] Project initialization
- [x] SvelteKit setup with TypeScript
- [x] Tailwind CSS + Daisy UI integration
- [x] ExcalidrawWrapper component scaffold
- [ ] Toolbar and Sidebar components
- [ ] Keyboard shortcuts

### Phase 2: Local Storage
- [ ] IndexedDB integration
- [ ] Export/Import functionality
- [ ] Drawing history management

### Phase 3: Backend Integration
- [ ] Go API endpoints
- [ ] Database schema
- [ ] User authentication
- [ ] Cloud synchronization

## Contributing

This is an open-source project. Contributions are welcome!

## License

MIT License - See LICENSE file for details

## Acknowledgments

- Built on top of [Excalidraw](https://excalidraw.com/)
- Uses [Svelte](https://svelte.dev/) for the frontend
- Styled with [Tailwind CSS](https://tailwindcss.com/) and [Daisy UI](https://daisyui.com/)
