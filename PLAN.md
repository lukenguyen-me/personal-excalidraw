# Development Plan

This document outlines the development roadmap for Personal Excalidraw.

## Phase 1: Frontend Infrastructure ✅

**Status**: Complete

- [x] Project initialization
- [x] SvelteKit setup with TypeScript
- [x] Tailwind CSS + DaisyUI integration
- [x] ExcalidrawWrapper component with React integration
- [x] Svelte stores for state management
- [x] Home page with drawings list
- [x] Drawing editor with dynamic routing
- [x] Wireframe theme configuration
- [x] Mock data with 8 sample drawings

**Outcome**: Functional frontend with mock data and basic navigation.

## Phase 2: Local Storage & Real Data 🚧

**Status**: In Progress

### Completed
- [x] Connect drawing editor to load/save specific drawings
- [x] LocalStorage persistence with per-drawing storage (`excalidraw-drawing-{id}`)
- [x] Auto-save with 1-second debounce for performance
- [x] Implement actual delete functionality (removes both metadata and content)
- [x] ID type abstraction for consistent type safety
- [x] Drawing metadata timestamp synchronization
- [x] Data validation and error handling (corrupted data, quota exceeded)

### In Progress
- [ ] Add drawing name editing
- [ ] Export/Import functionality (JSON, PNG, SVG)
- [ ] Drawing history management (undo/redo persistence)

### Goals
- Replace all mock data with real localStorage-backed data
- Implement full CRUD operations for drawings
- Add data import/export capabilities
- Ensure data persistence and recovery

## Phase 3: Backend Integration

**Status**: Not Started

### Backend Development
- [ ] Set up Go project structure
- [ ] Implement PostgreSQL database schema
- [ ] Create RESTful API endpoints:
  - [ ] `GET /api/drawings` - List all drawings
  - [ ] `GET /api/drawings/:id` - Get specific drawing
  - [ ] `POST /api/drawings` - Create new drawing
  - [ ] `PUT /api/drawings/:id` - Update drawing
  - [ ] `DELETE /api/drawings/:id` - Delete drawing
- [ ] Add database migrations
- [ ] Implement error handling and validation

### Frontend Integration
- [ ] Create API client service
- [ ] Replace localStorage calls with API calls
- [ ] Add loading states and error handling
- [ ] Implement data synchronization

### Authentication & Security
- [ ] User authentication system
- [ ] JWT token management
- [ ] Protected API routes
- [ ] User session handling

### Cloud Synchronization
- [ ] Real-time sync between devices
- [ ] Conflict resolution strategy
- [ ] Offline mode support
- [ ] Multi-user support

### Goals
- Transition from local-only to cloud-backed storage
- Enable cross-device access
- Prepare for multi-user features

## Phase 4: Enhanced Features

**Status**: Not Started

### Search & Organization
- [ ] Search and filter drawings by name
- [ ] Advanced search (by date, tags, content)
- [ ] Sorting options (name, created date, updated date)
- [ ] Pagination for large drawing lists
- [ ] Bulk operations (delete, export multiple)

### Visual Improvements
- [ ] Drawing thumbnails/previews in list view
- [ ] Grid view option (alternative to table)
- [ ] Preview modal on hover
- [ ] Canvas thumbnails generation

### Categorization
- [ ] Tags and categories system
- [ ] Color-coded labels
- [ ] Custom folders/collections
- [ ] Favorites/starred drawings

### Collaboration
- [ ] Sharing capabilities (public links)
- [ ] Collaborative editing (real-time)
- [ ] Comments and annotations
- [ ] Version history
- [ ] Access control (view/edit permissions)

### Export Options
- [ ] Batch export (multiple drawings)
- [ ] Multiple format support (PNG, SVG, PDF)
- [ ] Custom export settings (resolution, background)
- [ ] Automated backups

### Goals
- Enhance user experience with advanced features
- Support collaboration workflows
- Provide professional export options

## Future Considerations

### Performance
- [ ] Lazy loading for large drawing lists
- [ ] Canvas rendering optimization
- [ ] Image compression for embedded files
- [ ] Cache management

### Mobile Support
- [ ] Responsive mobile UI
- [ ] Touch gestures optimization
- [ ] Mobile-friendly toolbar
- [ ] Progressive Web App (PWA)

### Integrations
- [ ] Import from other drawing tools
- [ ] Export to Figma, Miro, etc.
- [ ] Markdown embedding
- [ ] API for third-party integrations

### Analytics
- [ ] Usage statistics
- [ ] Drawing metrics
- [ ] Performance monitoring
- [ ] Error tracking

## Milestones

| Phase | Target | Status |
|-------|--------|--------|
| Phase 1: Frontend Infrastructure | ✅ Complete | Done |
| Phase 2: Local Storage | 🚧 ~75% | In Progress |
| Phase 3: Backend Integration | 📅 Planned | Not Started |
| Phase 4: Enhanced Features | 📅 Planned | Not Started |

## Current Priority

**Focus**: Complete Phase 2 by finishing:
1. Drawing name editing functionality
2. Export/Import capabilities
3. Drawing history persistence

These features will provide a solid foundation before moving to backend integration.
