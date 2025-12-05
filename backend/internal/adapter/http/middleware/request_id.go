package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const RequestIDKey contextKey = "request_id"

// RequestID creates a middleware that generates or extracts request IDs
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get request ID from header
		requestID := r.Header.Get("X-Request-ID")

		// Generate new ID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
