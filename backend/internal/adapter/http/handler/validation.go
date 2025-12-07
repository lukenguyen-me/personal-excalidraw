package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ValidationError represents a validation error with field-specific details
type ValidationError struct {
	Field   string
	Message string
}

// ValidateCreateDrawingRequest validates the create drawing request
func ValidateCreateDrawingRequest(name string, data map[string]interface{}) []ValidationError {
	var errs []ValidationError

	// Validate name
	if strings.TrimSpace(name) == "" {
		errs = append(errs, ValidationError{
			Field:   "name",
			Message: "name cannot be empty",
		})
	}

	if len(name) > 255 {
		errs = append(errs, ValidationError{
			Field:   "name",
			Message: "name exceeds maximum length of 255 characters",
		})
	}

	// Validate data
	if data == nil {
		errs = append(errs, ValidationError{
			Field:   "data",
			Message: "data cannot be null",
		})
	}

	return errs
}

// ValidateUpdateDrawingRequest validates the update drawing request
func ValidateUpdateDrawingRequest(name string, data map[string]interface{}) []ValidationError {
	var errs []ValidationError

	// Both name and data are optional in updates
	// Only validate name if it's provided (non-empty)
	if name != "" {
		if len(name) > 255 {
			errs = append(errs, ValidationError{
				Field:   "name",
				Message: "name exceeds maximum length of 255 characters",
			})
		}
	}

	// Data is optional in updates (only validate if provided)
	// No validation needed for data field in updates

	return errs
}

// validateCreateDrawingRequest validates the CreateDrawingRequest
func validateCreateDrawingRequest(req *CreateDrawingRequest) error {
	errs := ValidateCreateDrawingRequest(req.Name, req.Data)
	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", errs[0].Message)
	}
	return nil
}

// validateUpdateDrawingRequest validates the UpdateDrawingRequest
func validateUpdateDrawingRequest(req *UpdateDrawingRequest) error {
	errs := ValidateUpdateDrawingRequest(req.Name, req.Data)
	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", errs[0].Message)
	}
	return nil
}

// ParseUUID safely parses a UUID string
func ParseUUID(s string) (uuid.UUID, error) {
	if s == "" {
		return uuid.Nil, errors.New("UUID cannot be empty")
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	return id, nil
}
