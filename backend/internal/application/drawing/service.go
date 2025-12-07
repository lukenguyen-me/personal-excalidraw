package drawing

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/personal-excalidraw/backend/internal/domain/drawing"
)

// Service handles drawing use cases
type Service struct {
	repo   drawing.Repository
	logger *slog.Logger
}

// NewService creates a new drawing service
func NewService(repo drawing.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateDrawing creates a new drawing
func (s *Service) CreateDrawing(ctx context.Context, input CreateDrawingInput) (*DrawingOutput, error) {
	s.logger.Info("creating drawing", "name", input.Name)

	// Create domain drawing
	d, err := drawing.NewDrawing(input.Name, input.Data)
	if err != nil {
		s.logger.Error("failed to create drawing domain object", "error", err)
		return nil, fmt.Errorf("failed to create drawing: %w", err)
	}

	// Persist to repository
	if err := s.repo.Create(ctx, d); err != nil {
		s.logger.Error("failed to persist drawing", "error", err)
		return nil, fmt.Errorf("failed to save drawing: %w", err)
	}

	s.logger.Info("drawing created successfully", "id", d.ID())

	return ToOutput(d), nil
}

// GetDrawing retrieves a single drawing by ID
func (s *Service) GetDrawing(ctx context.Context, id string) (*DrawingOutput, error) {
	s.logger.Info("getting drawing", "id", id)

	// Parse UUID from string
	drawingID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("invalid drawing ID format", "id", id, "error", err)
		return nil, fmt.Errorf("invalid drawing ID: %w", err)
	}

	// Retrieve from repository
	d, err := s.repo.FindByID(ctx, drawingID)
	if err != nil {
		s.logger.Error("failed to get drawing", "id", drawingID, "error", err)
		return nil, err
	}

	s.logger.Info("drawing retrieved successfully", "id", drawingID)

	return ToOutput(d), nil
}

// ListDrawings retrieves all drawings with pagination
func (s *Service) ListDrawings(ctx context.Context, input ListDrawingsInput) (*DrawingListOutput, error) {
	s.logger.Info("listing drawings", "limit", input.Limit, "offset", input.Offset)

	// Set default limit if not provided
	if input.Limit <= 0 {
		input.Limit = 10
	}

	// Ensure offset is not negative
	if input.Offset < 0 {
		input.Offset = 0
	}

	// Find all drawings with pagination
	drawings, err := s.repo.FindAll(ctx, input.Limit, input.Offset)
	if err != nil {
		s.logger.Error("failed to list drawings", "error", err)
		return nil, fmt.Errorf("failed to retrieve drawings: %w", err)
	}

	// Get total count
	total, err := s.repo.Count(ctx)
	if err != nil {
		s.logger.Error("failed to count drawings", "error", err)
		return nil, fmt.Errorf("failed to count drawings: %w", err)
	}

	s.logger.Info("drawings listed successfully", "count", len(drawings), "total", total)

	return &DrawingListOutput{
		Drawings: ToOutputList(drawings),
		Total:    total,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}, nil
}

// UpdateDrawing updates an existing drawing
func (s *Service) UpdateDrawing(ctx context.Context, id string, input UpdateDrawingInput) (*DrawingOutput, error) {
	s.logger.Info("updating drawing", "id", id)

	// Parse UUID from string
	drawingID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("invalid drawing ID format", "id", id, "error", err)
		return nil, fmt.Errorf("invalid drawing ID: %w", err)
	}

	// Retrieve from repository
	d, err := s.repo.FindByID(ctx, drawingID)
	if err != nil {
		s.logger.Error("failed to get drawing", "id", drawingID, "error", err)
		return nil, err
	}

	// Update the domain entity
	// If name is not provided (empty), keep the existing name
	nameToUpdate := input.Name
	if nameToUpdate == "" {
		nameToUpdate = d.Name()
	}

	// If data is not provided (nil), keep the existing data
	dataToUpdate := input.Data
	if dataToUpdate == nil {
		dataToUpdate = d.Data()
	}

	if err := d.Update(nameToUpdate, dataToUpdate); err != nil {
		s.logger.Error("failed to update drawing domain object", "error", err)
		return nil, fmt.Errorf("failed to update drawing: %w", err)
	}

	// Persist to repository
	if err := s.repo.Update(ctx, d); err != nil {
		s.logger.Error("failed to persist updated drawing", "error", err)
		return nil, fmt.Errorf("failed to save drawing: %w", err)
	}

	s.logger.Info("drawing updated successfully", "id", drawingID)

	return ToOutput(d), nil
}

// DeleteDrawing deletes an existing drawing
func (s *Service) DeleteDrawing(ctx context.Context, id string) error {
	s.logger.Info("deleting drawing", "id", id)

	// Parse UUID from string
	drawingID, err := uuid.Parse(id)
	if err != nil {
		s.logger.Error("invalid drawing ID format", "id", id, "error", err)
		return fmt.Errorf("invalid drawing ID: %w", err)
	}

	// Check if drawing exists
	_, err = s.repo.FindByID(ctx, drawingID)
	if err != nil {
		s.logger.Error("failed to get drawing", "id", drawingID, "error", err)
		return err
	}

	// Delete from repository
	if err := s.repo.Delete(ctx, drawingID); err != nil {
		s.logger.Error("failed to delete drawing", "id", drawingID, "error", err)
		return fmt.Errorf("failed to delete drawing: %w", err)
	}

	s.logger.Info("drawing deleted successfully", "id", drawingID)

	return nil
}
