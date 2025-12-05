package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	drawingapp "github.com/personal-excalidraw/backend/internal/application/drawing"
)

// DrawingHandler handles drawing HTTP requests
type DrawingHandler struct {
	service *drawingapp.Service
	logger  *slog.Logger
}

// NewDrawingHandler creates a new drawing handler
func NewDrawingHandler(service *drawingapp.Service, logger *slog.Logger) *DrawingHandler {
	return &DrawingHandler{
		service: service,
		logger:  logger,
	}
}

// DrawingResponse represents the HTTP response for a drawing
type DrawingResponse struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

// DrawingListResponse represents a paginated list response
type DrawingListResponse struct {
	Drawings []*DrawingResponse `json:"drawings"`
	Total    int64              `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}

// ListDrawings handles GET /api/drawings
func (h *DrawingHandler) ListDrawings(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling list drawings request")

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Call service
	input := drawingapp.ListDrawingsInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.service.ListDrawings(r.Context(), input)
	if err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Convert to HTTP response
	response := DrawingListResponse{
		Drawings: make([]*DrawingResponse, len(output.Drawings)),
		Total:    output.Total,
		Limit:    output.Limit,
		Offset:   output.Offset,
	}

	for i, d := range output.Drawings {
		response.Drawings[i] = &DrawingResponse{
			ID:        d.ID.String(),
			Name:      d.Name,
			Data:      d.Data,
			CreatedAt: d.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: d.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	respondJSON(w, http.StatusOK, response)
}
