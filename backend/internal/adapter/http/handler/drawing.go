package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	drawingapp "github.com/personal-excalidraw/backend/internal/application/drawing"
	"github.com/personal-excalidraw/backend/internal/adapter/http/util"
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

// CreateDrawingRequest represents the HTTP request for creating a drawing
type CreateDrawingRequest struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

// UpdateDrawingRequest represents the HTTP request for updating a drawing
type UpdateDrawingRequest struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
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

// CreateDrawing handles POST /api/drawings
func (h *DrawingHandler) CreateDrawing(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling create drawing request")

	// Parse request body
	var req CreateDrawingRequest
	if err := parseJSON(r, &req); err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Validate request
	if err := validateCreateDrawingRequest(&req); err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Call service
	input := drawingapp.CreateDrawingInput{
		Name: req.Name,
		Data: req.Data,
	}

	output, err := h.service.CreateDrawing(r.Context(), input)
	if err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Convert to HTTP response
	response := DrawingResponse{
		ID:        output.ID.String(),
		Name:      output.Name,
		Data:      output.Data,
		CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: output.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	util.RespondJSON(w, http.StatusCreated, response)
}

// GetDrawing handles GET /api/drawings/{id}
func (h *DrawingHandler) GetDrawing(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling get drawing request")

	// Extract ID from path
	id := r.PathValue("id")
	if id == "" {
		h.logger.Error("missing drawing ID in path")
		response := ErrorResponse{
			Error:   "invalid_request",
			Message: "missing drawing ID",
		}
		util.RespondJSON(w, http.StatusBadRequest, response)
		return
	}

	// Call service
	output, err := h.service.GetDrawing(r.Context(), id)
	if err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Convert to HTTP response
	response := DrawingResponse{
		ID:        output.ID.String(),
		Name:      output.Name,
		Data:      output.Data,
		CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: output.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	util.RespondJSON(w, http.StatusOK, response)
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

	util.RespondJSON(w, http.StatusOK, response)
}

// UpdateDrawing handles PUT /api/drawings/{id}
func (h *DrawingHandler) UpdateDrawing(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling update drawing request")

	// Extract ID from path
	id := r.PathValue("id")
	if id == "" {
		h.logger.Error("missing drawing ID in path")
		response := ErrorResponse{
			Error:   "invalid_request",
			Message: "missing drawing ID",
		}
		util.RespondJSON(w, http.StatusBadRequest, response)
		return
	}

	// Parse request body
	var req UpdateDrawingRequest
	if err := parseJSON(r, &req); err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Validate request
	if err := validateUpdateDrawingRequest(&req); err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Call service
	input := drawingapp.UpdateDrawingInput{
		Name: req.Name,
		Data: req.Data,
	}

	output, err := h.service.UpdateDrawing(r.Context(), id, input)
	if err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Convert to HTTP response
	response := DrawingResponse{
		ID:        output.ID.String(),
		Name:      output.Name,
		Data:      output.Data,
		CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: output.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	util.RespondJSON(w, http.StatusOK, response)
}

// DeleteDrawing handles DELETE /api/drawings/{id}
func (h *DrawingHandler) DeleteDrawing(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handling delete drawing request")

	// Extract ID from path
	id := r.PathValue("id")
	if id == "" {
		h.logger.Error("missing drawing ID in path")
		response := ErrorResponse{
			Error:   "invalid_request",
			Message: "missing drawing ID",
		}
		util.RespondJSON(w, http.StatusBadRequest, response)
		return
	}

	// Call service
	err := h.service.DeleteDrawing(r.Context(), id)
	if err != nil {
		respondError(w, err, h.logger)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
