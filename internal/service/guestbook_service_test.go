// Tests for guestbookService: mocks satisfy repository, validator, and cache ports (DIP).
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/model"
)

type MockGuestbookRepository struct {
	CreateFunc             func(ctx context.Context, entry *model.GuestbookEntry) error
	FindWithPaginationFunc func(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error)
}

func (m *MockGuestbookRepository) Create(ctx context.Context, entry *model.GuestbookEntry) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, entry)
	}
	return nil
}

func (m *MockGuestbookRepository) FindWithPagination(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error) {
	if m.FindWithPaginationFunc != nil {
		return m.FindWithPaginationFunc(ctx, limit, offset)
	}
	return []model.GuestbookEntry{}, nil
}

type MockValidator struct {
	ValidateAddEntryRequestFunc    func(req *pb.AddEntryRequest) error
	ValidateListEntriesRequestFunc func(req *pb.ListEntriesRequest) (*PaginationParams, error)
}

func (m *MockValidator) ValidateAddEntryRequest(req *pb.AddEntryRequest) error {
	if m.ValidateAddEntryRequestFunc != nil {
		return m.ValidateAddEntryRequestFunc(req)
	}
	return nil
}

func (m *MockValidator) ValidateListEntriesRequest(req *pb.ListEntriesRequest) (*PaginationParams, error) {
	if m.ValidateListEntriesRequestFunc != nil {
		return m.ValidateListEntriesRequestFunc(req)
	}
	return &PaginationParams{Limit: 10, Offset: 0}, nil
}

func TestGuestbookService_AddEntry_Success(t *testing.T) {
	mockRepo := &MockGuestbookRepository{
		CreateFunc: func(ctx context.Context, entry *model.GuestbookEntry) error {
			return nil
		},
	}

	mockValidator := &MockValidator{
		ValidateAddEntryRequestFunc: func(req *pb.AddEntryRequest) error {
			return nil
		},
	}

	service := NewGuestbookService(mockRepo, mockValidator)

	req := &pb.AddEntryRequest{
		Name:    "John Doe",
		Email:   "john@example.com",
		Message: "Hello, World!",
	}

	resp, err := service.AddEntry(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.Success {
		t.Error("expected success to be true")
	}
}
func TestGuestbookService_ListEntries_Success(t *testing.T) {
	mockRepo := &MockGuestbookRepository{
		FindWithPaginationFunc: func(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error) {
			return []model.GuestbookEntry{
				{Name: "John", Email: "john@example.com", Message: "Hello"},
			}, nil
		},
	}

	mockValidator := &MockValidator{
		ValidateListEntriesRequestFunc: func(req *pb.ListEntriesRequest) (*PaginationParams, error) {
			return &PaginationParams{Limit: 10, Offset: 0}, nil
		},
	}

	service := NewGuestbookService(mockRepo, mockValidator)

	req := &pb.ListEntriesRequest{
		Limit:  10,
		Offset: 0,
	}

	resp, err := service.ListEntries(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("expected 1 entry, got %d", len(resp.Data))
	}
}

func TestGuestbookService_ListEntries_ValidationError(t *testing.T) {
	mockRepo := &MockGuestbookRepository{}
	mockValidator := &MockValidator{
		ValidateListEntriesRequestFunc: func(req *pb.ListEntriesRequest) (*PaginationParams, error) {
			return nil, errors.New("invalid limit")
		},
	}

	service := NewGuestbookService(mockRepo, mockValidator)

	req := &pb.ListEntriesRequest{
		Limit:  -1,
		Offset: 0,
	}

	_, err := service.ListEntries(context.Background(), req)

	if err == nil {
		t.Error("expected validation error")
	}
}

func TestGuestbookService_ListEntries_RepositoryError(t *testing.T) {
	mockRepo := &MockGuestbookRepository{
		FindWithPaginationFunc: func(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error) {
			return nil, errors.New("database error")
		},
	}

	mockValidator := &MockValidator{
		ValidateListEntriesRequestFunc: func(req *pb.ListEntriesRequest) (*PaginationParams, error) {
			return &PaginationParams{Limit: 10, Offset: 0}, nil
		},
	}

	service := NewGuestbookService(mockRepo, mockValidator)

	req := &pb.ListEntriesRequest{
		Limit:  10,
		Offset: 0,
	}

	_, err := service.ListEntries(context.Background(), req)

	if err == nil {
		t.Error("expected repository error")
	}
}

type MockCache struct {
	SetFunc    func(ctx context.Context, key string, value any, expiration time.Duration) error
	GetFunc    func(ctx context.Context, key string, dest any) (bool, error)
	DeleteFunc func(ctx context.Context, key string) error
}

func (m *MockCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if m.SetFunc != nil {
		return m.SetFunc(ctx, key, value, expiration)
	}
	return nil
}

func (m *MockCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key, dest)
	}
	return false, nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key)
	}
	return nil
}

func TestGuestbookService_WithCache(t *testing.T) {
	mockRepo := &MockGuestbookRepository{
		CreateFunc: func(ctx context.Context, entry *model.GuestbookEntry) error {
			return nil
		},
		FindWithPaginationFunc: func(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error) {
			return []model.GuestbookEntry{}, nil
		},
	}

	mockValidator := &MockValidator{
		ValidateAddEntryRequestFunc: func(req *pb.AddEntryRequest) error {
			return nil
		},
		ValidateListEntriesRequestFunc: func(req *pb.ListEntriesRequest) (*PaginationParams, error) {
			return &PaginationParams{Limit: 10, Offset: 0}, nil
		},
	}

	mockCache := &MockCache{}

	service := NewGuestbookServiceWithCache(mockRepo, mockValidator, mockCache)

	// Test AddEntry with cache
	addReq := &pb.AddEntryRequest{
		Name:    "John Doe",
		Email:   "john@example.com",
		Message: "Hello",
	}

	_, err := service.AddEntry(context.Background(), addReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test ListEntries with cache
	listReq := &pb.ListEntriesRequest{
		Limit:  10,
		Offset: 0,
	}

	_, err = service.ListEntries(context.Background(), listReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGuestbookService_AddEntry_ValidationError(t *testing.T) {
	mockRepo := &MockGuestbookRepository{}
	mockValidator := &MockValidator{
		ValidateAddEntryRequestFunc: func(req *pb.AddEntryRequest) error {
			return errors.New("name is required")
		},
	}

	service := NewGuestbookService(mockRepo, mockValidator)

	req := &pb.AddEntryRequest{
		Message: "Hello, World!",
	}

	resp, err := service.AddEntry(context.Background(), req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Success {
		t.Error("expected success to be false")
	}

	if resp.Error != "name is required" {
		t.Errorf("expected error message 'name is required', got '%s'", resp.Error)
	}
}
