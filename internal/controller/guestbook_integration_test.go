// Integration tests: real SQLite, repository, service, and controller wired together.
package controller

import (
	"context"
	"testing"
	"time"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
	"github.com/novaherdi/go-proto/internal/model"
	"github.com/novaherdi/go-proto/internal/repository"
	"github.com/novaherdi/go-proto/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for integration testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&model.GuestbookEntry{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

// setupController creates a fully wired controller with all dependencies for integration testing
func setupController(t *testing.T) *GuestbookController {
	db := setupTestDB(t)
	repo := repository.NewGormGuestbookRepository(db)
	validator := service.NewDefaultValidator()
	svc := service.NewGuestbookService(repo, validator)
	controller := NewGuestbookController(svc)

	return controller
}

func TestGuestbookController_Integration_AddAndListEntries(t *testing.T) {
	controller := setupController(t)

	// Add multiple entries
	entries := []*pb.AddEntryRequest{
		{
			Name:    "John Doe",
			Email:   "john@example.com",
			Message: "Hello, World!",
		},
		{
			Name:    "Jane Smith",
			Email:   "jane@example.com",
			Message: "Great service!",
		},
		{
			Name:    "Bob Johnson",
			Message: "Nice to meet you!", // No email
		},
	}

	// Add entries with slight delay to ensure proper ordering
	for i, entry := range entries {
		resp, err := controller.AddEntry(context.Background(), entry)
		if err != nil {
			t.Fatalf("failed to add entry %d: %v", i, err)
		}
		if !resp.Success {
			t.Errorf("expected success for entry %d, got error: %s", i, resp.Error)
		}
		// Delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// List entries
	listResp, err := controller.ListEntries(context.Background(), &pb.ListEntriesRequest{})
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}

	if len(listResp.Data) != 3 {
		t.Errorf("expected 3 entries, got %d", len(listResp.Data))
	}

	// Verify entries are ordered by created_at DESC (newest first)
	// Since entries are added in order: John, Jane, Bob
	// and Bob is the last one added, Bob should be first in the response
	if listResp.Data[0].Name != "Bob Johnson" {
		// Log the actual order for debugging
		t.Logf("Actual order:")
		for i, entry := range listResp.Data {
			t.Logf("%d: %s (created_at: %s)", i, entry.Name, entry.CreatedAt)
		}
		t.Errorf("expected first entry to be 'Bob Johnson', got '%s'", listResp.Data[0].Name)
	}

	if listResp.Data[2].Name != "John Doe" {
		t.Errorf("expected last entry to be 'John Doe', got '%s'", listResp.Data[2].Name)
	}
}

func TestGuestbookController_Integration_ValidationErrors(t *testing.T) {
	controller := setupController(t)

	testCases := []struct {
		name        string
		request     *pb.AddEntryRequest
		expectError bool
	}{
		{
			name: "missing name",
			request: &pb.AddEntryRequest{
				Email:   "test@example.com",
				Message: "Hello",
			},
			expectError: true,
		},
		{
			name: "missing message",
			request: &pb.AddEntryRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			expectError: true,
		},
		{
			name: "valid request without email",
			request: &pb.AddEntryRequest{
				Name:    "Test User",
				Message: "Hello",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := controller.AddEntry(context.Background(), tc.request)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.expectError {
				if resp.Success {
					t.Error("expected validation error but got success")
				}
				if resp.Error == "" {
					t.Error("expected error message but got empty string")
				}
			} else {
				if !resp.Success {
					t.Errorf("expected success but got error: %s", resp.Error)
				}
			}
		})
	}
}

func TestGuestbookController_Integration_Pagination(t *testing.T) {
	controller := setupController(t)

	// Add more entries than the default limit
	for i := 0; i < 15; i++ {
		req := &pb.AddEntryRequest{
			Name:    "User " + string(rune('A'+i)),
			Message: "Message from user " + string(rune('A'+i)),
		}

		resp, err := controller.AddEntry(context.Background(), req)
		if err != nil {
			t.Fatalf("failed to add entry %d: %v", i, err)
		}
		if !resp.Success {
			t.Errorf("failed to add entry %d: %s", i, resp.Error)
		}
		// Delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)
	}

	// Test default pagination (should get 10 entries)
	listResp, err := controller.ListEntries(context.Background(), &pb.ListEntriesRequest{})
	if err != nil {
		t.Fatalf("failed to list entries: %v", err)
	}

	if len(listResp.Data) != 10 {
		t.Errorf("expected 10 entries with default pagination, got %d", len(listResp.Data))
	}

	// Test custom pagination
	listResp, err = controller.ListEntries(context.Background(), &pb.ListEntriesRequest{
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		t.Fatalf("failed to list entries with custom pagination: %v", err)
	}

	if len(listResp.Data) != 5 {
		t.Errorf("expected 5 entries, got %d", len(listResp.Data))
	}

	// Test second page
	listResp, err = controller.ListEntries(context.Background(), &pb.ListEntriesRequest{
		Limit:  5,
		Offset: 5,
	})
	if err != nil {
		t.Fatalf("failed to list second page: %v", err)
	}

	if len(listResp.Data) != 5 {
		t.Errorf("expected 5 entries on second page, got %d", len(listResp.Data))
	}
}

func TestGuestbookController_Integration_EmptyDatabase(t *testing.T) {
	controller := setupController(t)

	// List entries from empty database
	listResp, err := controller.ListEntries(context.Background(), &pb.ListEntriesRequest{})
	if err != nil {
		t.Fatalf("failed to list entries from empty database: %v", err)
	}

	if len(listResp.Data) != 0 {
		t.Errorf("expected 0 entries from empty database, got %d", len(listResp.Data))
	}
}
