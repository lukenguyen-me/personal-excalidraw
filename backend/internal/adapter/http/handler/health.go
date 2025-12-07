package handler

import (
	"net/http"

	"github.com/personal-excalidraw/backend/internal/adapter/http/util"
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
	util.RespondJSON(w, http.StatusOK, response)
}
