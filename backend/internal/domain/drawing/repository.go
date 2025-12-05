package drawing

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the contract for drawing persistence
type Repository interface {
	// Create stores a new drawing
	Create(ctx context.Context, drawing *Drawing) error

	// FindByID retrieves a drawing by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Drawing, error)

	// FindAll retrieves all drawings with pagination
	FindAll(ctx context.Context, limit, offset int) ([]*Drawing, error)

	// Update updates an existing drawing
	Update(ctx context.Context, drawing *Drawing) error

	// Delete removes a drawing by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// Count returns the total number of drawings
	Count(ctx context.Context) (int64, error)
}
