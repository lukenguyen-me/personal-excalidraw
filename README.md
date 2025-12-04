# Personal Excalidraw

A self-hosted, open-source drawing application built on Excalidraw. Create and manage your diagrams with a clean, modern interface.

## Features

- Free drawing application based on Excalidraw
- Drawings list with create, edit, and delete operations
- Local-first data storage (localStorage/future API integration)
- Modern UI with SvelteKit and DaisyUI wireframe theme
- Responsive design with clean table layout
- Future: Go backend for cloud synchronization and PostgreSQL storage

## Project Structure

```
personal-excalidraw/
├── frontend/                  # SvelteKit + TypeScript frontend
│   ├── src/
│   │   ├── routes/
│   │   │   ├── +layout.svelte           # Root layout
│   │   │   ├── +page.svelte             # Home page (drawings list)
│   │   │   └── drawing/[id]/
│   │   │       └── +page.svelte         # Drawing editor page
│   │   ├── lib/
│   │   │   ├── components/
│   │   │   │   └── ExcalidrawWrapper.svelte  # Excalidraw integration
│   │   │   └── stores/
│   │   │       ├── mockDrawings.ts      # Mock data store
│   │   │       ├── drawing.ts           # Drawing state
│   │   │       ├── excalidraw.ts        # Excalidraw API
│   │   │       └── ui.ts                # UI state
│   │   └── routes/layout.css            # Global styles
├── backend/                   # Go backend (stub for future)
├── docker-compose.yml         # Development environment
└── Dockerfile                 # Production build
```

## Quick Start

### Prerequisites

- Node.js 18+ (for frontend development)
- pnpm (package manager)
- Go 1.21+ (for future backend development)
- Docker & Docker Compose (for containerized deployment)

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
- Backend server on `http://localhost:8080` (stub)

### Building for Production

```bash
docker build -t personal-excalidraw:latest .
docker run -p 8080:8080 personal-excalidraw:latest
```

## Tech Stack

### Frontend
- **Framework**: SvelteKit 2.48.5 with Svelte 5.43.8
- **Build Tool**: Vite 7.2.6
- **Styling**: Tailwind CSS 4.1.17 + DaisyUI 5.5.5 (Wireframe theme)
- **Drawing Engine**: @excalidraw/excalidraw (React wrapper)
- **State Management**: Svelte 5 Stores (runes-based)
- **Language**: TypeScript 5.7.3

### Backend (Future)
- **Language**: Go (net/http)
- **Database**: PostgreSQL
- **API**: RESTful API for drawing CRUD operations

### Deployment
- **Containerization**: Docker + Docker Compose
- **Hosting**: Self-hosted

## Current Features

### Home Page (/)
- List of all drawings in a clean table format
- Columns: Name, Created At, Updated At, Actions
- **New** button to create a new drawing
- **Edit** button (icon) to open drawing editor
- **Delete** button (icon) for removing drawings (placeholder)
- Wireframe theme with responsive layout

### Drawing Editor (/drawing/[id])
- Full Excalidraw canvas integration
- Header with:
  - Back button to return to drawings list
  - Status indicator (placeholder for save/load states)
- Dynamic routing with drawing ID parameter
- Empty canvas (data persistence coming in future phases)

### State Management
- **Drawing Store**: Manages drawing elements, app state, and files
- **Excalidraw Store**: Handles Excalidraw API reference
- **UI Store**: Tracks UI state (sidebar, zoom, tools, etc.)
- **Mock Drawings Store**: 8 sample drawings with incremental IDs

## Development Roadmap

### Phase 1: Frontend Infrastructure ✅

- [x] Project initialization
- [x] SvelteKit setup with TypeScript
- [x] Tailwind CSS + DaisyUI integration
- [x] ExcalidrawWrapper component with React integration
- [x] Svelte stores for state management
- [x] Home page with drawings list
- [x] Drawing editor with dynamic routing
- [x] Wireframe theme configuration
- [x] Mock data with 8 sample drawings

### Phase 2: Local Storage & Real Data

- [ ] Connect drawing editor to load/save specific drawings
- [ ] Implement actual delete functionality
- [ ] Add drawing name editing
- [ ] LocalStorage persistence for drawings
- [ ] Export/Import functionality
- [ ] Drawing history management

### Phase 3: Backend Integration

- [ ] Go API endpoints (CRUD operations)
- [ ] PostgreSQL database schema
- [ ] Replace mock data with real API calls
- [ ] User authentication
- [ ] Cloud synchronization
- [ ] Multi-user support

### Phase 4: Enhanced Features

- [ ] Search and filter drawings
- [ ] Sorting by name/date
- [ ] Pagination for large drawing lists
- [ ] Drawing thumbnails/previews
- [ ] Tags and categories
- [ ] Sharing capabilities

## API Endpoints (Future)

```
GET    /api/drawings          # List all drawings
GET    /api/drawings/:id      # Get specific drawing
POST   /api/drawings          # Create new drawing
PUT    /api/drawings/:id      # Update drawing
DELETE /api/drawings/:id      # Delete drawing
```

## Contributing

This is an open-source project. Contributions are welcome!

## License

MIT License - See LICENSE file for details

## Acknowledgments

- Built on top of [Excalidraw](https://excalidraw.com/)
- Uses [SvelteKit](https://svelte.dev/) for the frontend framework
- Styled with [Tailwind CSS](https://tailwindcss.com/) and [DaisyUI](https://daisyui.com/)
