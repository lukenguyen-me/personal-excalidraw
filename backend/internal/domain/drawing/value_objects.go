package drawing

import (
	"encoding/json"
	"fmt"
)

// DrawingData represents the JSONB drawing content
// It stores the drawing elements, appState, and files
type DrawingData map[string]interface{}

// Validate ensures the drawing data structure is valid
func (d DrawingData) Validate() error {
	if d == nil {
		return fmt.Errorf("%w: data cannot be nil", ErrInvalidDrawingData)
	}

	// Ensure it's a valid map that can be marshaled to JSON
	if _, err := json.Marshal(d); err != nil {
		return fmt.Errorf("%w: cannot marshal to JSON: %v", ErrInvalidDrawingData, err)
	}

	return nil
}

// ToJSON converts DrawingData to JSON bytes for storage
func (d DrawingData) ToJSON() ([]byte, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	data, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to marshal: %v", ErrInvalidDrawingData, err)
	}

	return data, nil
}

// FromJSON creates DrawingData from JSON bytes
func FromJSON(data []byte) (DrawingData, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: JSON data is empty", ErrInvalidDrawingData)
	}

	var dd DrawingData
	if err := json.Unmarshal(data, &dd); err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal: %v", ErrInvalidDrawingData, err)
	}

	return dd, nil
}
