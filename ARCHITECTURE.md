# Architecture

This document describes the technical architecture of Personal Excalidraw.

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

## State Management

### Drawing Store
Manages drawing elements, app state, and files with ID-aware persistence:
- Centralized drawing data management
- Automatic localStorage persistence
- Per-drawing storage with unique keys (`excalidraw-drawing-{id}`)
- Auto-save with 1-second debounce for performance

### Excalidraw Store
Handles Excalidraw API reference:
- Provides access to Excalidraw instance methods
- Enables programmatic control of the canvas
- Manages Excalidraw lifecycle

### UI Store
Tracks UI state:
- Sidebar visibility
- Zoom level
- Active tools
- View mode preferences

### Mock Drawings Store
Manages drawing metadata with reactive updates:
- CRUD operations for drawing metadata
- Timestamp synchronization
- ID-based drawing lookup

### Type System
Centralized `ID` type for consistent identifier handling:
- Type-safe string-based IDs
- Compile-time validation
- Consistent across all stores

## Data Storage

### LocalStorage Schema (Current - Phase 2)

Drawings are stored in localStorage with the following schema:

```
Key: excalidraw-drawing-{id}
Value: {
  elements: [...],      // Excalidraw drawing elements (shapes, lines, text, etc.)
  appState: {...},      // Canvas state (zoom, scroll position, view mode, etc.)
  files: {...}          // Embedded images and files
}
```

**Features:**
- Per-drawing storage with unique keys
- Automatic JSON serialization/deserialization
- Data validation on load (validates structure and element arrays)
- Error handling for corrupted data and quota exceeded
- Automatic cleanup of invalid entries

**Limitations:**
- localStorage has ~5-10MB limit per domain
- Data is browser-specific (not synced across devices)
- Will be replaced with backend API in Phase 3

### Database Schema (Future - Phase 3)

PostgreSQL schema for cloud storage:

```sql
-- Drawings table
CREATE TABLE drawings (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id UUID REFERENCES users(id),
  data JSONB NOT NULL  -- Stores elements, appState, and files
);

-- Users table (for future multi-user support)
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## API Endpoints (Future - Phase 3)

RESTful API design:

```
GET    /api/drawings          # List all drawings
GET    /api/drawings/:id      # Get specific drawing
POST   /api/drawings          # Create new drawing
PUT    /api/drawings/:id      # Update drawing
DELETE /api/drawings/:id      # Delete drawing
```

**Request/Response Examples:**

```json
// GET /api/drawings
{
  "drawings": [
    {
      "id": "uuid",
      "name": "Architecture Diagram",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T14:20:00Z"
    }
  ]
}

// GET /api/drawings/:id
{
  "id": "uuid",
  "name": "Architecture Diagram",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T14:20:00Z",
  "data": {
    "elements": [...],
    "appState": {...},
    "files": {...}
  }
}
```

## Component Architecture

### ExcalidrawWrapper
A Svelte component that wraps the React-based Excalidraw library:
- Handles React-Svelte interop
- Manages Excalidraw lifecycle
- Provides Svelte-friendly API
- Auto-save integration
- State synchronization with stores

### Drawing Editor Page
- Dynamic routing based on drawing ID
- Loads drawing data from localStorage
- Integrates ExcalidrawWrapper
- Provides navigation and status UI
- Implements auto-save with debounce

### Drawings List Page
- Table-based layout with DaisyUI wireframe theme
- CRUD operations (Create, Read, Update, Delete)
- Responsive design
- Navigation to editor

## Build & Deployment

### Development Environment
```bash
# Frontend only
cd frontend
pnpm install
pnpm run dev

# With Docker Compose
docker-compose up
```

### Production Build
```bash
# Multi-stage Docker build
docker build -t personal-excalidraw:latest .
docker run -p 8080:8080 personal-excalidraw:latest
```

The production build:
1. Builds frontend with Vite (SvelteKit adapter-static)
2. Serves static files (future: Go backend will serve)
3. Optimizes for self-hosting

## Performance Considerations

### Auto-Save Optimization
- 1-second debounce prevents excessive writes
- Only saves when changes are detected
- Batches rapid edits into single save operation

### Data Validation
- Validates drawing data structure on load
- Handles corrupted data gracefully
- Provides fallback to empty canvas

### Error Handling
- localStorage quota exceeded detection
- Corrupted data recovery
- Automatic cleanup of invalid entries
- User-friendly error messages
