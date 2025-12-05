package drawing

import "errors"

var (
	// ErrDrawingNotFound is returned when a drawing is not found
	ErrDrawingNotFound = errors.New("drawing not found")

	// ErrInvalidDrawingName is returned when a drawing name is invalid
	ErrInvalidDrawingName = errors.New("invalid drawing name")

	// ErrInvalidDrawingData is returned when drawing data is invalid
	ErrInvalidDrawingData = errors.New("invalid drawing data")

	// ErrEmptyName is returned when a drawing name is empty
	ErrEmptyName = errors.New("drawing name cannot be empty")

	// ErrNameTooLong is returned when a drawing name exceeds maximum length
	ErrNameTooLong = errors.New("drawing name exceeds maximum length")
)
