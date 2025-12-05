package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (rw *responseWriter) WriteHeader(status int) {
	if !rw.wrote {
		rw.status = status
		rw.wrote = true
		rw.ResponseWriter.WriteHeader(status)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// Logger creates a middleware for request/response logging
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status
			wrapped := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Log request details
			duration := time.Since(start)
			logger.Info("HTTP request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.status,
				"duration_ms", duration.Milliseconds(),
				"remote_addr", r.RemoteAddr,
			)
		})
	}
}
