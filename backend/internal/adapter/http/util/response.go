package util

import (
	"encoding/json"
	"io"
	"net/http"
)

// RespondJSON sends a JSON response with the given status code and data
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			// Log error but can't change status code at this point
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// DecodeJSON decodes JSON from an io.Reader into the provided interface
func DecodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
