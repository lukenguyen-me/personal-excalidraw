package drawing

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// MaxNameLength is the maximum allowed length for a drawing name
	MaxNameLength = 255
)

// Drawing represents the drawing aggregate root
type Drawing struct {
	id        uuid.UUID
	slug      string
	name      string
	data      DrawingData
	createdAt time.Time
	updatedAt time.Time
}

// NewDrawing creates a new drawing with validation
func NewDrawing(name string, data DrawingData) (*Drawing, error) {
	d := &Drawing{
		id:        uuid.New(),
		slug:      "", // Slug will be set by the service layer
		name:      name,
		data:      data,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}

// SetSlug sets the slug for the drawing (to be called by service layer)
func (d *Drawing) SetSlug(slug string) {
	d.slug = slug
}

// Reconstitute creates a drawing from persisted data (for repository use)
func Reconstitute(id uuid.UUID, slug, name string, data DrawingData, createdAt, updatedAt time.Time) (*Drawing, error) {
	d := &Drawing{
		id:        id,
		slug:      slug,
		name:      name,
		data:      data,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}

// Update updates the drawing with new data
func (d *Drawing) Update(name string, data DrawingData) error {
	d.name = name
	d.data = data
	d.updatedAt = time.Now().UTC()

	return d.Validate()
}

// Validate ensures the drawing is in a valid state
func (d *Drawing) Validate() error {
	// Validate name
	if err := d.validateName(); err != nil {
		return err
	}

	// Validate data
	if err := d.data.Validate(); err != nil {
		return err
	}

	return nil
}

// validateName checks if the name is valid
func (d *Drawing) validateName() error {
	trimmed := strings.TrimSpace(d.name)

	if trimmed == "" {
		return ErrEmptyName
	}

	if len(d.name) > MaxNameLength {
		return ErrNameTooLong
	}

	return nil
}

// ID returns the drawing ID
func (d *Drawing) ID() uuid.UUID {
	return d.id
}

// Slug returns the drawing slug
func (d *Drawing) Slug() string {
	return d.slug
}

// Name returns the drawing name
func (d *Drawing) Name() string {
	return d.name
}

// Data returns the drawing data
func (d *Drawing) Data() DrawingData {
	return d.data
}

// CreatedAt returns the creation timestamp
func (d *Drawing) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns the last update timestamp
func (d *Drawing) UpdatedAt() time.Time {
	return d.updatedAt
}
