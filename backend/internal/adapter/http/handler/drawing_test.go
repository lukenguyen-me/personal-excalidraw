package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	drawingapp "github.com/personal-excalidraw/backend/internal/application/drawing"
	"github.com/personal-excalidraw/backend/internal/domain/drawing"
)

// mockDrawingRepository is a mock implementation for testing
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
		name           string
		requestBody    interface{}
		mockRepo       *mockDrawingRepository
		expectedStatus int
		validateResp   func(t *testing.T, body []byte)
	}{
		{
			name: "successful creation",
			requestBody: CreateDrawingRequest{
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
			expectedStatus: http.StatusCreated,
			validateResp: func(t *testing.T, body []byte) {
				var resp DrawingResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Name != "Test Drawing" {
					t.Errorf("expected name 'Test Drawing', got '%s'", resp.Name)
				}
				if resp.ID == "" {
					t.Error("expected non-empty ID")
				}
			},
		},
		{
			name:           "empty request body",
			requestBody:    nil,
			mockRepo:       &mockDrawingRepository{},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ErrorResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
			},
		},
		{
			name: "invalid request - empty name",
			requestBody: CreateDrawingRequest{
				Name: "",
				Data: map[string]interface{}{"elements": []interface{}{}},
			},
			mockRepo:       &mockDrawingRepository{},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ErrorResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
			},
		},
		{
			name: "invalid request - nil data",
			requestBody: CreateDrawingRequest{
				Name: "Test",
				Data: nil,
			},
			mockRepo:       &mockDrawingRepository{},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ErrorResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
			},
		},
		{
			name: "repository error",
			requestBody: CreateDrawingRequest{
				Name: "Test Drawing",
				Data: map[string]interface{}{"elements": []interface{}{}},
			},
			mockRepo: &mockDrawingRepository{
				createFunc: func(ctx context.Context, d *drawing.Drawing) error {
					return errors.New("database connection failed")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ErrorResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if resp.Error != "internal_error" {
					t.Errorf("expected error type 'internal_error', got '%s'", resp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := drawingapp.NewService(tt.mockRepo, logger)
			handler := NewDrawingHandler(service, logger)

			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/drawings", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateDrawing(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, w.Body.Bytes())
			}
		})
	}
}

func TestListDrawings(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	tests := []struct {
		name           string
		queryParams    string
		mockRepo       *mockDrawingRepository
		expectedStatus int
		validateResp   func(t *testing.T, body []byte)
	}{
		{
			name:        "successful list with defaults",
			queryParams: "",
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					d1, _ := drawing.NewDrawing("Drawing 1", map[string]interface{}{"elements": []interface{}{}})
					d2, _ := drawing.NewDrawing("Drawing 2", map[string]interface{}{"elements": []interface{}{}})
					return []*drawing.Drawing{d1, d2}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 2, nil
				},
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, body []byte) {
				var resp DrawingListResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(resp.Drawings) != 2 {
					t.Errorf("expected 2 drawings, got %d", len(resp.Drawings))
				}
				if resp.Total != 2 {
					t.Errorf("expected total 2, got %d", resp.Total)
				}
				if resp.Limit != 10 {
					t.Errorf("expected limit 10, got %d", resp.Limit)
				}
				if resp.Offset != 0 {
					t.Errorf("expected offset 0, got %d", resp.Offset)
				}
			},
		},
		{
			name:        "successful list with custom pagination",
			queryParams: "?limit=5&offset=10",
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
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, body []byte) {
				var resp DrawingListResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if resp.Limit != 5 {
					t.Errorf("expected limit 5, got %d", resp.Limit)
				}
				if resp.Offset != 10 {
					t.Errorf("expected offset 10, got %d", resp.Offset)
				}
			},
		},
		{
			name:        "empty list",
			queryParams: "",
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 0, nil
				},
			},
			expectedStatus: http.StatusOK,
			validateResp: func(t *testing.T, body []byte) {
				var resp DrawingListResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}
				if len(resp.Drawings) != 0 {
					t.Errorf("expected 0 drawings, got %d", len(resp.Drawings))
				}
				if resp.Total != 0 {
					t.Errorf("expected total 0, got %d", resp.Total)
				}
			},
		},
		{
			name:        "invalid pagination parameters - negative values ignored",
			queryParams: "?limit=-5&offset=-10",
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 0, nil
				},
			},
			expectedStatus: http.StatusOK,
			validateResp:   func(t *testing.T, body []byte) {},
		},
		{
			name:        "repository error on FindAll",
			queryParams: "",
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return nil, errors.New("database connection failed")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ErrorResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if resp.Error != "internal_error" {
					t.Errorf("expected error type 'internal_error', got '%s'", resp.Error)
				}
			},
		},
		{
			name:        "repository error on Count",
			queryParams: "",
			mockRepo: &mockDrawingRepository{
				findAllFunc: func(ctx context.Context, limit, offset int) ([]*drawing.Drawing, error) {
					return []*drawing.Drawing{}, nil
				},
				countFunc: func(ctx context.Context) (int64, error) {
					return 0, errors.New("database connection failed")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			validateResp: func(t *testing.T, body []byte) {
				var resp ErrorResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal error response: %v", err)
				}
				if resp.Error != "internal_error" {
					t.Errorf("expected error type 'internal_error', got '%s'", resp.Error)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := drawingapp.NewService(tt.mockRepo, logger)
			handler := NewDrawingHandler(service, logger)

			req := httptest.NewRequest(http.MethodGet, "/drawings"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			handler.ListDrawings(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.validateResp != nil {
				tt.validateResp(t, w.Body.Bytes())
			}
		})
	}
}
