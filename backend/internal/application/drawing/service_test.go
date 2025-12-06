package drawing

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/personal-excalidraw/backend/internal/domain/drawing"
)

// mockDrawingRepository is a mock implementation of the drawing repository
type mockDrawingRepository struct {
	createFunc   func(ctx context.Context, d *drawing.Drawing) error
	findAllFunc  func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error)
	countFunc    func(ctx context.Context) (int64, error)
	findByIDFunc func(ctx context.Context, id uuid.UUID) (*drawing.Drawing, error)
	updateFunc   func(ctx context.Context, d *drawing.Drawing) error
	deleteFunc   func(ctx context.Context, id uuid.UUID) error
}

func (m *mockDrawingRepository) Create(ctx context.Context, d *drawing.Drawing) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, d)
	}
	return errors.New("not implemented")
}

func (m *mockDrawingRepository) FindAll(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx, limit, offset)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDrawingRepository) Count(ctx context.Context) (int64, error) {
	if m.countFunc != nil {
		return m.countFunc(ctx)
	}
	return 0, errors.New("not implemented")
}

func (m *mockDrawingRepository) FindByID(ctx context.Context, id uuid.UUID) (*drawing.Drawing, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockDrawingRepository) Update(ctx context.Context, d *drawing.Drawing) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, d)
	}
	return errors.New("not implemented")
}

func (m *mockDrawingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestCreateDrawing(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name        string
		input       CreateDrawingInput
		mockRepo    *mockDrawingRepository
		expectError bool
		validateOut func(t *testing.T, out *DrawingOutput)
	}{
		{
			name: "successful creation",
			input: CreateDrawingInput{
				Name: "Test Drawing",
				Data: map[string]interface{}{
					"elements": []interface{}{},
					"appState": map[string]interface{}{},
				},
			},
			mockRepo: &mockDrawingRepository{
				createFunc: func(ctx context.Context, d *drawing.Drawing) error {
					return nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingOutput) {
				if out == nil {
					t.Fatal("expected non-nil output")
				}
				if out.Name != "Test Drawing" {
					t.Errorf("expected name 'Test Drawing', got '%s'", out.Name)
				}
				if out.ID == uuid.Nil {
					t.Error("expected non-nil UUID")
				}
				if out.Data == nil {
					t.Error("expected non-nil data")
				}
				if out.CreatedAt.IsZero() {
					t.Error("expected non-zero created_at timestamp")
				}
				if out.UpdatedAt.IsZero() {
					t.Error("expected non-zero updated_at timestamp")
				}
			},
		},
		{
			name: "empty name validation error",
			input: CreateDrawingInput{
				Name: "",
				Data: map[string]interface{}{"elements": []interface{}{}},
			},
			mockRepo:    &mockDrawingRepository{},
			expectError: true,
			validateOut: func(t *testing.T, out *DrawingOutput) {
				if out != nil {
					t.Error("expected nil output on error")
				}
			},
		},
		{
			name: "name too long validation error",
			input: CreateDrawingInput{
				Name: string(make([]byte, 300)), // Exceeds 255 character limit
				Data: map[string]interface{}{"elements": []interface{}{}},
			},
			mockRepo:    &mockDrawingRepository{},
			expectError: true,
			validateOut: func(t *testing.T, out *DrawingOutput) {
				if out != nil {
					t.Error("expected nil output on error")
				}
			},
		},
		{
			name: "nil data validation error",
			input: CreateDrawingInput{
				Name: "Test",
				Data: nil,
			},
			mockRepo:    &mockDrawingRepository{},
			expectError: true,
			validateOut: func(t *testing.T, out *DrawingOutput) {
				if out != nil {
					t.Error("expected nil output on error")
				}
			},
		},
		{
			name: "repository error",
			input: CreateDrawingInput{
				Name: "Test Drawing",
				Data: map[string]interface{}{"elements": []interface{}{}},
			},
			mockRepo: &mockDrawingRepository{
				createFunc: func(ctx context.Context, d *drawing.Drawing) error {
					return errors.New("database connection failed")
				},
			},
			expectError: true,
			validateOut: func(t *testing.T, out *DrawingOutput) {
				if out != nil {
					t.Error("expected nil output on error")
				}
			},
		},
		{
			name: "complex drawing data",
			input: CreateDrawingInput{
				Name: "Complex Drawing",
				Data: map[string]interface{}{
					"elements": []interface{}{
						map[string]interface{}{
							"type": "rectangle",
							"x":    100,
							"y":    200,
						},
					},
					"appState": map[string]interface{}{
						"viewBackgroundColor": "#ffffff",
					},
					"files": map[string]interface{}{},
				},
			},
			mockRepo: &mockDrawingRepository{
				createFunc: func(ctx context.Context, d *drawing.Drawing) error {
					return nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingOutput) {
				if out == nil {
					t.Fatal("expected non-nil output")
				}
				elements, ok := out.Data["elements"]
				if !ok {
					t.Error("expected elements in data")
				}
				if elements == nil {
					t.Error("expected non-nil elements")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.mockRepo, logger)
			ctx := context.Background()

			output, err := service.CreateDrawing(ctx, tt.input)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.validateOut != nil {
				tt.validateOut(t, output)
			}
		})
	}
}

func TestListDrawings(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name        string
		input       ListDrawingsInput
		mockRepo    *mockDrawingRepository
		expectError bool
		validateOut func(t *testing.T, out *DrawingListOutput)
	}{
		{
			name: "successful list with default limit",
			input: ListDrawingsInput{
				Limit:  0, // Should default to 10
				Offset: 0,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					if limit != 10 {
						t.Errorf("expected limit 10, got %d", limit)
					}
					d1, _ := drawing.NewDrawing("Drawing 1", map[string]interface{}{"elements": []interface{}{}})
					d2, _ := drawing.NewDrawing("Drawing 2", map[string]interface{}{"elements": []interface{}{}})
					return []*drawing.Drawing{d1, d2}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 2, nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out == nil {
					t.Fatal("expected non-nil output")
				}
				if len(out.Drawings) != 2 {
					t.Errorf("expected 2 drawings, got %d", len(out.Drawings))
				}
				if out.Total != 2 {
					t.Errorf("expected total 2, got %d", out.Total)
				}
				if out.Limit != 10 {
					t.Errorf("expected limit 10, got %d", out.Limit)
				}
			},
		},
		{
			name: "successful list with custom pagination",
			input: ListDrawingsInput{
				Limit:  5,
				Offset: 10,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					if limit != 5 {
						t.Errorf("expected limit 5, got %d", limit)
					}
					if offset != 10 {
						t.Errorf("expected offset 10, got %d", offset)
					}
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 100, nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out == nil {
					t.Fatal("expected non-nil output")
				}
				if out.Total != 100 {
					t.Errorf("expected total 100, got %d", out.Total)
				}
				if out.Limit != 5 {
					t.Errorf("expected limit 5, got %d", out.Limit)
				}
				if out.Offset != 10 {
					t.Errorf("expected offset 10, got %d", out.Offset)
				}
			},
		},
		{
			name: "empty list",
			input: ListDrawingsInput{
				Limit:  10,
				Offset: 0,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 0, nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out == nil {
					t.Fatal("expected non-nil output")
				}
				if len(out.Drawings) != 0 {
					t.Errorf("expected 0 drawings, got %d", len(out.Drawings))
				}
				if out.Total != 0 {
					t.Errorf("expected total 0, got %d", out.Total)
				}
			},
		},
		{
			name: "negative offset normalized to 0",
			input: ListDrawingsInput{
				Limit:  10,
				Offset: -5,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					if offset != 0 {
						t.Errorf("expected offset 0, got %d", offset)
					}
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 0, nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out.Offset != 0 {
					t.Errorf("expected offset 0, got %d", out.Offset)
				}
			},
		},
		{
			name: "repository error on FindAll",
			input: ListDrawingsInput{
				Limit:  10,
				Offset: 0,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return nil, errors.New("database connection failed")
				},
			},
			expectError: true,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out != nil {
					t.Error("expected nil output on error")
				}
			},
		},
		{
			name: "repository error on Count",
			input: ListDrawingsInput{
				Limit:  10,
				Offset: 0,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 0, errors.New("database connection failed")
				},
			},
			expectError: true,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out != nil {
					t.Error("expected nil output on error")
				}
			},
		},
		{
			name: "large dataset pagination",
			input: ListDrawingsInput{
				Limit:  100,
				Offset: 500,
			},
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					drawings := make([]*drawing.Drawing, 100)
					for i := 0; i < 100; i++ {
						d, _ := drawing.NewDrawing("Drawing", map[string]interface{}{"elements": []interface{}{}})
						drawings[i] = d
					}
					return drawings, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 10000, nil
				},
			},
			expectError: false,
			validateOut: func(t *testing.T, out *DrawingListOutput) {
				if out == nil {
					t.Fatal("expected non-nil output")
				}
				if len(out.Drawings) != 100 {
					t.Errorf("expected 100 drawings, got %d", len(out.Drawings))
				}
				if out.Total != 10000 {
					t.Errorf("expected total 10000, got %d", out.Total)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.mockRepo, logger)
			ctx := context.Background()

			output, err := service.ListDrawings(ctx, tt.input)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.validateOut != nil {
				tt.validateOut(t, output)
			}
		})
	}
}

func TestToOutput(t *testing.T) {
	now := time.Now()
	data := map[string]interface{}{"elements": []interface{}{}}
	d, err := drawing.NewDrawing("Test", data)
	if err != nil {
		t.Fatalf("failed to create drawing: %v", err)
	}

	output := ToOutput(d)

	if output == nil {
		t.Fatal("expected non-nil output")
	}
	if output.ID != d.ID() {
		t.Error("ID mismatch")
	}
	if output.Name != d.Name() {
		t.Error("Name mismatch")
	}
	if output.CreatedAt.Before(now.Add(-time.Second)) {
		t.Error("CreatedAt should be recent")
	}
}

func TestToOutputList(t *testing.T) {
	d1, _ := drawing.NewDrawing("Drawing 1", map[string]interface{}{"elements": []interface{}{}})
	d2, _ := drawing.NewDrawing("Drawing 2", map[string]interface{}{"elements": []interface{}{}})

	drawings := []*drawing.Drawing{d1, d2}
	outputs := ToOutputList(drawings)

	if len(outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(outputs))
	}

	if outputs[0].Name != "Drawing 1" {
		t.Errorf("expected 'Drawing 1', got '%s'", outputs[0].Name)
	}
	if outputs[1].Name != "Drawing 2" {
		t.Errorf("expected 'Drawing 2', got '%s'", outputs[1].Name)
	}
}
