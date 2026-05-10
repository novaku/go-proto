// Tests for GORM guestbook repository against an in-memory SQLite database.
package repository

import (
	"context"
	"testing"
	"time"

	"github.com/novaherdi/go-proto/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := db.AutoMigrate(&model.GuestbookEntry{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

func TestGormGuestbookRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormGuestbookRepository(db)

	entry := &model.GuestbookEntry{
		Name:      "John Doe",
		Email:     "john@example.com",
		Message:   "Hello, World!",
		CreatedAt: time.Now(),
	}

	err := repo.Create(context.Background(), entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entry.ID == 0 {
		t.Error("expected entry ID to be set after creation")
	}
}

func TestGormGuestbookRepository_FindWithPagination(t *testing.T) {
	db := setupTestDB(t)
	repo := NewGormGuestbookRepository(db)

	// Insert test data
	now := time.Now()
	entries := []model.GuestbookEntry{
		{Name: "John", Message: "Hello", CreatedAt: now},
		{Name: "Jane", Message: "Hi", CreatedAt: now.Add(-100 * time.Second)},
	}

	for _, entry := range entries {
		if err := db.Create(&entry).Error; err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	results, err := repo.FindWithPagination(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 entries, got %d", len(results))
	}

	// Verify ordering (newest first)
	if results[0].Name != "John" {
		t.Errorf("expected first entry to be John, got %s", results[0].Name)
	}
}
