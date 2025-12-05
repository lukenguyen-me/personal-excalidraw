package handler

import (
	"net/http"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "ok",
	}
	respondJSON(w, http.StatusOK, response)
}
