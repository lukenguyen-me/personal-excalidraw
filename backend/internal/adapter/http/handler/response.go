package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/personal-excalidraw/backend/internal/domain/drawing"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// parseJSON parses JSON from request body
func parseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errors.New("invalid JSON format")
	}

	return nil
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Log error but can't change status code at this point
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// respondError sends an error response with appropriate status code
func respondError(w http.ResponseWriter, err error, logger *slog.Logger) {
	status, errorType, message := mapErrorToHTTP(err)

	logger.Error("request error", "error", err, "status", status, "message", message)

	response := ErrorResponse{
		Error:   errorType,
		Message: message,
	}

	respondJSON(w, status, response)
}

// respondValidationError sends a validation error response
func respondValidationError(w http.ResponseWriter, validationErrors []ValidationError) {
	details := make(map[string]string)
	for _, ve := range validationErrors {
		details[ve.Field] = ve.Message
	}

	response := ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid request data",
		Details: details,
	}

	respondJSON(w, http.StatusBadRequest, response)
}

// respondNotFound sends a 404 response
func respondNotFound(w http.ResponseWriter, message string) {
	response := ErrorResponse{
		Error:   "not_found",
		Message: message,
	}

	respondJSON(w, http.StatusNotFound, response)
}

// mapErrorToHTTP maps domain errors to HTTP status codes
func mapErrorToHTTP(err error) (status int, errorType, message string) {
	switch {
	case errors.Is(err, drawing.ErrDrawingNotFound):
		return http.StatusNotFound, "not_found", "Drawing not found"
	case errors.Is(err, drawing.ErrInvalidDrawingName):
		return http.StatusBadRequest, "invalid_name", "Invalid drawing name"
	case errors.Is(err, drawing.ErrInvalidDrawingData):
		return http.StatusBadRequest, "invalid_data", "Invalid drawing data"
	case errors.Is(err, drawing.ErrEmptyName):
		return http.StatusBadRequest, "empty_name", "Drawing name cannot be empty"
	case errors.Is(err, drawing.ErrNameTooLong):
		return http.StatusBadRequest, "name_too_long", "Drawing name exceeds maximum length"
	default:
		return http.StatusInternalServerError, "internal_error", "Internal server error"
	}
}
