package http

import (
	"log/slog"
	"net/http"

	"github.com/personal-excalidraw/backend/internal/adapter/http/handler"
	"github.com/personal-excalidraw/backend/internal/adapter/http/middleware"
	"github.com/personal-excalidraw/backend/internal/infrastructure/config"
)

// NewRouter creates a new HTTP router with all routes and middleware
func NewRouter(
	cfg *config.Config,
	healthHandler *handler.HealthHandler,
	drawingHandler *handler.DrawingHandler,
	logger *slog.Logger,
) http.Handler {
	// Create new ServeMux with Go 1.22+ routing
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", healthHandler.Check)

	// Drawing API endpoints (nginx strips /api prefix)
	mux.HandleFunc("GET /drawings", drawingHandler.ListDrawings)

	// Apply middleware stack (in reverse order - outermost first)
	var handler http.Handler = mux
	handler = middleware.CORS(cfg)(handler)
	handler = middleware.Logger(logger)(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recover(logger)(handler)

	return handler
}
