package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/personal-excalidraw/backend/internal/domain/drawing"
)

// DrawingRepository implements the drawing.Repository interface using PostgreSQL
type DrawingRepository struct {
	pool *pgxpool.Pool
}

// NewDrawingRepository creates a new DrawingRepository
func NewDrawingRepository(pool *pgxpool.Pool) *DrawingRepository {
	return &DrawingRepository{
		pool: pool,
	}
}

// Create stores a new drawing in the database
func (r *DrawingRepository) Create(ctx context.Context, d *drawing.Drawing) error {
	// Convert drawing data to JSON bytes
	dataJSON, err := d.Data().ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal drawing data: %w", err)
	}

	// Execute insert query
	_, err = r.pool.Exec(
		ctx,
		queryCreateDrawing,
		d.ID(),
		d.Name(),
		dataJSON,
		d.CreatedAt(),
		d.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("failed to create drawing: %w", err)
	}

	return nil
}

// FindByID retrieves a drawing by its ID
func (r *DrawingRepository) FindByID(ctx context.Context, id uuid.UUID) (*drawing.Drawing, error) {
	var (
		drawingID uuid.UUID
		name      string
		dataJSON  []byte
		createdAt, updatedAt time.Time
	)

	// Execute select query
	err := r.pool.QueryRow(ctx, queryFindDrawingByID, id).Scan(
		&drawingID,
		&name,
		&dataJSON,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, drawing.ErrDrawingNotFound
		}
		return nil, fmt.Errorf("failed to find drawing: %w", err)
	}

	// Parse drawing data from JSON
	data, err := drawing.FromJSON(dataJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal drawing data: %w", err)
	}

	// Reconstitute the drawing entity
	d, err := drawing.Reconstitute(drawingID, name, data, createdAt, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstitute drawing: %w", err)
	}

	return d, nil
}

// FindAll retrieves all drawings with pagination
func (r *DrawingRepository) FindAll(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
	// Execute select query
	rows, err := r.pool.Query(ctx, queryFindAllDrawings, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find all drawings: %w", err)
	}
	defer rows.Close()

	// Collect drawings
	var drawings []*drawing.Drawing
	for rows.Next() {
		var (
			drawingID uuid.UUID
			name      string
			dataJSON  []byte
			createdAt, updatedAt time.Time
		)

		if err := rows.Scan(&drawingID, &name, &dataJSON, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan drawing row: %w", err)
		}

		// Parse drawing data from JSON
		data, err := drawing.FromJSON(dataJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal drawing data: %w", err)
		}

		// Reconstitute the drawing entity
		d, err := drawing.Reconstitute(drawingID, name, data, createdAt, updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstitute drawing: %w", err)
		}

		drawings = append(drawings, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating drawing rows: %w", err)
	}

	return drawings, nil
}

// Update updates an existing drawing in the database
func (r *DrawingRepository) Update(ctx context.Context, d *drawing.Drawing) error {
	// Convert drawing data to JSON bytes
	dataJSON, err := d.Data().ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal drawing data: %w", err)
	}

	// Execute update query
	result, err := r.pool.Exec(
		ctx,
		queryUpdateDrawing,
		d.Name(),
		dataJSON,
		d.UpdatedAt(),
		d.ID(),
	)
	if err != nil {
		return fmt.Errorf("failed to update drawing: %w", err)
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		return drawing.ErrDrawingNotFound
	}

	return nil
}

// Delete removes a drawing from the database
func (r *DrawingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Execute delete query
	result, err := r.pool.Exec(ctx, queryDeleteDrawing, id)
	if err != nil {
		return fmt.Errorf("failed to delete drawing: %w", err)
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		return drawing.ErrDrawingNotFound
	}

	return nil
}

// Count returns the total number of drawings in the database
func (r *DrawingRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	// Execute count query
	err := r.pool.QueryRow(ctx, queryCountDrawings).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count drawings: %w", err)
	}

	return count, nil
}
