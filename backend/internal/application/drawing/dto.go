package drawing

import (
	"time"

	"github.com/google/uuid"

	"github.com/personal-excalidraw/backend/internal/domain/drawing"
)

// CreateDrawingInput represents input for creating a drawing
type CreateDrawingInput struct {
	Name string
	Data map[string]interface{}
}

// UpdateDrawingInput represents input for updating a drawing
type UpdateDrawingInput struct {
	Name string
	Data map[string]interface{}
}

// ListDrawingsInput represents input for listing drawings
type ListDrawingsInput struct {
	Limit  int
	Offset int
}

// DrawingOutput represents a drawing response
type DrawingOutput struct {
	ID        uuid.UUID
	Name      string
	Data      map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DrawingListOutput represents a paginated list of drawings
type DrawingListOutput struct {
	Drawings []*DrawingOutput
	Total    int64
	Limit    int
	Offset   int
}

// ToOutput converts a domain drawing to a DrawingOutput DTO
func ToOutput(d *drawing.Drawing) *DrawingOutput {
	return &DrawingOutput{
		ID:        d.ID(),
		Name:      d.Name(),
		Data:      d.Data(),
		CreatedAt: d.CreatedAt(),
		UpdatedAt: d.UpdatedAt(),
	}
}

// ToOutputList converts a list of domain drawings to DrawingOutput DTOs
func ToOutputList(drawings []*drawing.Drawing) []*DrawingOutput {
	outputs := make([]*DrawingOutput, len(drawings))
	for i, d := range drawings {
		outputs[i] = ToOutput(d)
	}
	return outputs
}
