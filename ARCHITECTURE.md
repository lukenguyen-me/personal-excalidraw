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
│   │   │   ├── api/
│   │   │   │   └── drawings.ts          # API client with TypeScript types
│   │   │   ├── components/
│   │   │   │   ├── ExcalidrawWrapper.svelte  # Excalidraw integration
│   │   │   │   └── AuthModal.svelte     # Authentication modal
│   │   │   ├── stores/
│   │   │   │   ├── drawing.ts           # Drawing state
│   │   │   │   ├── excalidraw.ts        # Excalidraw API
│   │   │   │   ├── auth.ts              # Authentication state
│   │   │   │   └── ui.ts                # UI state
│   │   │   └── types/
│   │   │       └── index.ts             # Shared TypeScript types
│   │   └── routes/layout.css            # Global styles
├── backend/                   # Go backend with Clean Architecture
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                  # Application entry point
│   ├── internal/
│   │   ├── adapter/
│   │   │   └── http/
│   │   │       ├── handler/             # HTTP handlers (controllers)
│   │   │       │   ├── drawing.go       # Drawing endpoints
│   │   │       │   └── auth.go          # Auth validation endpoint
│   │   │       └── middleware/          # HTTP middleware
│   │   │           └── auth.go          # Authentication middleware
│   │   ├── application/
│   │   │   └── drawing/                 # Business logic services
│   │   ├── domain/
│   │   │   └── drawing/                 # Domain models
│   │   └── infrastructure/
│   │       ├── config/                  # Configuration management
│   │       └── database/                # Database repositories & migrations
│   └── go.mod
├── docker-compose.yml         # Development environment
└── Dockerfile                 # Production build
```

## Tech Stack

### Frontend
- **Framework**: SvelteKit 2.48.5 with Svelte 5.45.4
- **Build Tool**: Vite 7.2.6
- **Styling**: Tailwind CSS 4.1.17
- **Drawing Engine**: @excalidraw/excalidraw (React wrapper)
- **Data Fetching**: TanStack Query v5 for caching and synchronization
- **State Management**: Svelte 5 Stores (runes-based)
- **Language**: TypeScript 5.9.3

### Backend
- **Language**: Go 1.24+
- **Architecture**: Clean Architecture (Handler → Service → Repository)
- **Database**: PostgreSQL 14+ with database/sql
- **API**: RESTful API with comprehensive validation
- **Testing**: Table-driven tests with testify
- **Migrations**: Custom migration system with reversible migrations

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

### Auth Store
Manages authentication state with localStorage persistence:
- Access key storage and retrieval
- Authentication status tracking
- Cross-tab synchronization via storage events
- Automatic persistence to localStorage (`excalidraw_access_key`)
- Helper methods: `getAccessKey()`, `setAccessKey()`, `clearAccessKey()`, `hasAccessKey()`

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

### Type System
Centralized `ID` type for consistent identifier handling:
- Type-safe string-based IDs
- Compile-time validation
- Consistent across all stores

## Data Storage

### Database Schema (PostgreSQL)

Current production schema:

```sql
-- Drawings table
CREATE TABLE drawings (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  data JSONB NOT NULL,  -- Stores elements, appState, and files
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_drawings_slug ON drawings(slug);
CREATE INDEX idx_drawings_created_at ON drawings(created_at DESC);
CREATE INDEX idx_drawings_updated_at ON drawings(updated_at DESC);
```

**Features:**
- Auto-incrementing integer IDs
- Slug-based URLs for SEO (e.g., `/drawing/my-architecture-diagram`)
- JSONB storage for flexible Excalidraw data
- Automatic timestamp management
- Optimized indexes for common queries

**Migration System:**
- Version-controlled migrations in `backend/internal/infrastructure/database/migrations/`
- Reversible migrations (up/down)
- Automatic migration execution on startup
- Safe, tested migration scripts

### Future Schema Extensions

For multi-user support (Phase 4+):

```sql
-- Users table
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add user_id to drawings
ALTER TABLE drawings ADD COLUMN user_id INTEGER REFERENCES users(id);
```

## API Endpoints

RESTful API implementation with authentication:

```
# Authentication
GET    /api/auth/validate     # Validate access key (requires Bearer token)

# Drawings (all protected by auth middleware when enabled)
GET    /api/drawings          # List all drawings
GET    /api/drawings/:id      # Get specific drawing
POST   /api/drawings          # Create new drawing
PUT    /api/drawings/:id      # Update drawing (supports partial updates)
DELETE /api/drawings/:id      # Delete drawing
```

**Authentication:**
- All `/api/drawings/*` endpoints require Bearer token authentication when `AUTH_ENABLED=true`
- Access key is validated using constant-time comparison to prevent timing attacks
- Token format: `Authorization: Bearer <your-access-key>`
- Unauthenticated requests receive structured error responses with error codes

**Request/Response Examples:**

```json
// GET /api/drawings
{
  "drawings": [
    {
      "id": 1,
      "name": "Architecture Diagram",
      "slug": "architecture-diagram",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T14:20:00Z"
    }
  ]
}

// GET /api/drawings/1
{
  "id": 1,
  "name": "Architecture Diagram",
  "slug": "architecture-diagram",
  "data": {
    "elements": [...],
    "appState": {...},
    "files": {...}
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T14:20:00Z"
}

// POST /api/drawings
{
  "name": "New Diagram",
  "data": {
    "elements": [],
    "appState": {},
    "files": {}
  }
}

// PUT /api/drawings/1 (partial update)
{
  "name": "Updated Name"  // Only update name, leave data unchanged
}
```

**Validation:**
- Name: 1-255 characters, required
- Slug: Auto-generated from name, unique, URL-safe
- Data: Valid JSON object with elements, appState, and files
- Comprehensive error messages for validation failures

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

## Security

### Authentication System

**Access Key Authentication:**
- Simple Bearer token-based authentication using access keys
- Configurable via environment variables (`ACCESS_KEY`, `AUTH_ENABLED`)
- Can be disabled for local development or trusted environments

**Backend Security:**
- Authentication middleware protects all `/api/drawings/*` endpoints
- Constant-time comparison using `crypto/subtle.ConstantTimeCompare()` prevents timing attacks
- Structured error responses with specific error codes
- Public path exemptions for health checks and auth validation

**Frontend Security:**
- Access key stored securely in localStorage
- AuthModal component for user authentication
- Automatic token injection in API requests via Authorization header
- Cross-tab synchronization ensures consistent auth state
- Token validation before allowing access to protected features

**Error Responses:**
```json
// Missing token
{
  "error": "Unauthorized",
  "message": "Access key required",
  "code": "AUTH_REQUIRED"
}

// Invalid format
{
  "error": "Unauthorized",
  "message": "Invalid authorization format. Use: Bearer <key>",
  "code": "INVALID_AUTH_FORMAT"
}

// Wrong key
{
  "error": "Unauthorized",
  "message": "Invalid access key",
  "code": "INVALID_ACCESS_KEY"
}
```

**Configuration:**
```bash
# Enable authentication
AUTH_ENABLED=true
ACCESS_KEY=your-secret-key-here

# Disable for development
AUTH_ENABLED=false
```

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
