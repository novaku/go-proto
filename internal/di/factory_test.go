package di

import (
	"context"
	"testing"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/model"
	"github.com/novaherdi/go-proto/internal/repository"
	"github.com/novaherdi/go-proto/internal/service"
	"github.com/novaherdi/go-proto/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDBForFactory prepares an in-memory SQLite DB for factory integration tests.
func setupTestDBForFactory(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := db.AutoMigrate(&model.GuestbookEntry{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestNewGuestbookControllerFactory(t *testing.T) {
	db := setupTestDBForFactory(t)

	factory := NewGuestbookControllerFactory(db)

	if factory == nil {
		t.Fatal("NewGuestbookControllerFactory returned nil")
	}

	gc := factory.CreateController()
	if gc == nil {
		t.Fatal("CreateController returned nil")
	}
}

func TestNewGuestbookControllerFactoryWithCache_Disabled(t *testing.T) {
	db := setupTestDBForFactory(t)

	redisConfig := config.RedisConfig{
		Host:    "localhost",
		Port:    6379,
		Enabled: false,
	}

	factory := NewGuestbookControllerFactoryWithCache(db, redisConfig)

	if factory == nil {
		t.Fatal("NewGuestbookControllerFactoryWithCache returned nil")
	}

	if factory.CreateController() == nil {
		t.Fatal("CreateController returned nil")
	}
}

func TestNewGuestbookControllerFactoryWithCache_Enabled(t *testing.T) {
	db := setupTestDBForFactory(t)

	redisConfig := config.RedisConfig{
		Host:    "localhost",
		Port:    6379,
		Enabled: true,
	}

	factory := NewGuestbookControllerFactoryWithCache(db, redisConfig)

	if factory == nil {
		t.Fatal("NewGuestbookControllerFactoryWithCache returned nil")
	}

	t.Log("Factory created with cache config (cache may be nil if Redis unavailable)")
}

func TestGuestbookControllerFactory_CreateController(t *testing.T) {
	db := setupTestDBForFactory(t)

	factory := NewGuestbookControllerFactory(db)
	gc := factory.CreateController()

	if gc == nil {
		t.Fatal("CreateController returned nil")
	}

	ctx := context.Background()
	req := &pb.AddEntryRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Test message",
	}

	resp, err := gc.AddEntry(ctx, req)

	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got error: %s", resp.Error)
	}
}

func TestGuestbookControllerFactory_CreateControllerWithCache(t *testing.T) {
	db := setupTestDBForFactory(t)

	redisConfig := config.RedisConfig{
		Host:    "localhost",
		Port:    6379,
		Enabled: false,
	}

	factory := NewGuestbookControllerFactoryWithCache(db, redisConfig)
	gc := factory.CreateController()

	if gc == nil {
		t.Fatal("CreateController returned nil")
	}

	ctx := context.Background()
	req := &pb.AddEntryRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Test message",
	}

	resp, err := gc.AddEntry(ctx, req)

	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got error: %s", resp.Error)
	}
}

// MockValidator is a test double for GuestbookRequestValidator.
type MockValidator struct {
	ValidateAddEntryRequestFunc    func(req *pb.AddEntryRequest) error
	ValidateListEntriesRequestFunc func(req *pb.ListEntriesRequest) (*service.PaginationParams, error)
}

func (m *MockValidator) ValidateAddEntryRequest(req *pb.AddEntryRequest) error {
	if m.ValidateAddEntryRequestFunc != nil {
		return m.ValidateAddEntryRequestFunc(req)
	}
	return nil
}

func (m *MockValidator) ValidateListEntriesRequest(req *pb.ListEntriesRequest) (*service.PaginationParams, error) {
	if m.ValidateListEntriesRequestFunc != nil {
		return m.ValidateListEntriesRequestFunc(req)
	}
	return &service.PaginationParams{Limit: 10, Offset: 0}, nil
}

func TestGuestbookControllerFactory_CreateControllerWithCustomValidator(t *testing.T) {
	db := setupTestDBForFactory(t)

	factory := NewGuestbookControllerFactory(db)

	customValidator := &MockValidator{
		ValidateAddEntryRequestFunc: func(req *pb.AddEntryRequest) error {
			return nil
		},
	}

	gc := factory.CreateControllerWithCustomValidator(customValidator)

	if gc == nil {
		t.Fatal("CreateControllerWithCustomValidator returned nil")
	}

	ctx := context.Background()
	req := &pb.AddEntryRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Test message",
	}

	resp, err := gc.AddEntry(ctx, req)

	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got error: %s", resp.Error)
	}
}

// MockRepository is a test double for GuestbookRepository.
type MockRepository struct {
	CreateFunc             func(ctx context.Context, entry *model.GuestbookEntry) error
	FindWithPaginationFunc func(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error)
}

func (m *MockRepository) Create(ctx context.Context, entry *model.GuestbookEntry) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, entry)
	}
	return nil
}

func (m *MockRepository) FindWithPagination(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error) {
	if m.FindWithPaginationFunc != nil {
		return m.FindWithPaginationFunc(ctx, limit, offset)
	}
	return []model.GuestbookEntry{}, nil
}

func TestGuestbookControllerFactory_CreateControllerWithCustomRepository(t *testing.T) {
	db := setupTestDBForFactory(t)

	factory := NewGuestbookControllerFactory(db)

	customRepo := &MockRepository{
		CreateFunc: func(ctx context.Context, entry *model.GuestbookEntry) error {
			return nil
		},
	}

	gc := factory.CreateControllerWithCustomRepository(customRepo)

	if gc == nil {
		t.Fatal("CreateControllerWithCustomRepository returned nil")
	}

	ctx := context.Background()
	req := &pb.AddEntryRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Test message",
	}

	resp, err := gc.AddEntry(ctx, req)

	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got error: %s", resp.Error)
	}
}

func TestGuestbookControllerFactory_WithCacheEnabled(t *testing.T) {
	db := setupTestDBForFactory(t)

	redisConfig := config.RedisConfig{
		Host:    "invalid-redis-host",
		Port:    6379,
		Enabled: true,
	}

	factory := NewGuestbookControllerFactoryWithCache(db, redisConfig)
	gc := factory.CreateController()

	if gc == nil {
		t.Fatal("CreateController returned nil")
	}

	ctx := context.Background()
	req := &pb.AddEntryRequest{
		Name:    "Test User",
		Email:   "test@example.com",
		Message: "Test message",
	}

	resp, err := gc.AddEntry(ctx, req)

	if err != nil {
		t.Fatalf("AddEntry failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected success, got error: %s", resp.Error)
	}
}

func TestGuestbookControllerFactory_MultipleControllers(t *testing.T) {
	db := setupTestDBForFactory(t)

	factory := NewGuestbookControllerFactory(db)

	gc1 := factory.CreateController()
	gc2 := factory.CreateController()

	if gc1 == nil || gc2 == nil {
		t.Fatal("Failed to create controllers")
	}

	ctx := context.Background()

	req1 := &pb.AddEntryRequest{
		Name:    "User 1",
		Email:   "user1@example.com",
		Message: "Message 1",
	}

	req2 := &pb.AddEntryRequest{
		Name:    "User 2",
		Email:   "user2@example.com",
		Message: "Message 2",
	}

	resp1, _ := gc1.AddEntry(ctx, req1)
	resp2, _ := gc2.AddEntry(ctx, req2)

	if !resp1.Success || !resp2.Success {
		t.Error("Both controllers should work successfully")
	}
}

func TestGuestbookControllerFactory_ListEntries(t *testing.T) {
	db := setupTestDBForFactory(t)

	factory := NewGuestbookControllerFactory(db)
	gc := factory.CreateController()

	ctx := context.Background()
	for i := 0; i < 3; i++ {
		req := &pb.AddEntryRequest{
			Name:    "Test User",
			Email:   "test@example.com",
			Message: "Test message",
		}
		_, _ = gc.AddEntry(ctx, req)
	}

	listReq := &pb.ListEntriesRequest{
		Limit:  10,
		Offset: 0,
	}

	listResp, err := gc.ListEntries(ctx, listReq)

	if err != nil {
		t.Fatalf("ListEntries failed: %v", err)
	}

	if len(listResp.Data) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(listResp.Data))
	}
}

func TestGuestbookControllerFactory_CustomValidatorWithCache(t *testing.T) {
	db := setupTestDBForFactory(t)

	redisConfig := config.RedisConfig{
		Host:    "localhost",
		Port:    6379,
		Enabled: false,
	}

	factory := NewGuestbookControllerFactoryWithCache(db, redisConfig)

	customValidator := &MockValidator{
		ValidateAddEntryRequestFunc: func(req *pb.AddEntryRequest) error {
			return nil
		},
	}

	gc := factory.CreateControllerWithCustomValidator(customValidator)

	if gc == nil {
		t.Fatal("CreateControllerWithCustomValidator returned nil")
	}
}

func TestGuestbookControllerFactory_CustomRepositoryWithCache(t *testing.T) {
	db := setupTestDBForFactory(t)

	redisConfig := config.RedisConfig{
		Host:    "localhost",
		Port:    6379,
		Enabled: false,
	}

	factory := NewGuestbookControllerFactoryWithCache(db, redisConfig)

	customRepo := repository.NewGormGuestbookRepository(db)

	gc := factory.CreateControllerWithCustomRepository(customRepo)

	if gc == nil {
		t.Fatal("CreateControllerWithCustomRepository returned nil")
	}
}
