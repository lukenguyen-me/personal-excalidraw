package handler

import (
	"net/http"

	"github.com/personal-excalidraw/backend/internal/adapter/http/util"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct{}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Validate handles the GET /auth/validate endpoint
// If this endpoint is reached, the auth middleware has already validated the access key
func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	util.RespondJSON(w, http.StatusOK, map[string]bool{
		"authenticated": true,
	})
}
