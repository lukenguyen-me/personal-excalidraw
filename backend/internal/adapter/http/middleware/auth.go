package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/personal-excalidraw/backend/internal/infrastructure/config"
)

// Auth creates a middleware for handling authentication
func Auth(cfg *config.Config, publicPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth for public paths
			for _, path := range publicPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Skip if auth disabled
			if !cfg.Auth.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error":   "Unauthorized",
					"message": "Access key required",
					"code":    "AUTH_REQUIRED",
				})
				return
			}

			// Check Bearer format
			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error":   "Unauthorized",
					"message": "Invalid authorization format. Use: Bearer <key>",
					"code":    "INVALID_AUTH_FORMAT",
				})
				return
			}

			// Extract and validate token
			token := strings.TrimPrefix(authHeader, prefix)

			// Use constant-time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(token), []byte(cfg.Auth.AccessKey)) != 1 {
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error":   "Unauthorized",
					"message": "Invalid access key",
					"code":    "INVALID_ACCESS_KEY",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// respondJSON is a helper function to send JSON responses
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
