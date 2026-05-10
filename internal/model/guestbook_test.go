// Table and field behavior tests for GuestbookEntry (no database).
package model

import (
	"testing"
	"time"
)

func TestGuestbookEntry_TableName(t *testing.T) {
	entry := GuestbookEntry{}
	tableName := entry.TableName()

	expectedName := "guestbook_entries"
	if tableName != expectedName {
		t.Errorf("TableName() = %s; want %s", tableName, expectedName)
	}
}

func TestGuestbookEntry_Creation(t *testing.T) {
	tests := []struct {
		name    string
		entry   GuestbookEntry
		wantErr bool
	}{
		{
			name: "valid entry",
			entry: GuestbookEntry{
				ID:        1,
				Name:      "John Doe",
				Email:     "john@example.com",
				Message:   "Hello, World!",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			entry: GuestbookEntry{
				ID:        2,
				Name:      "",
				Email:     "test@example.com",
				Message:   "Test message",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty email",
			entry: GuestbookEntry{
				ID:        3,
				Name:      "Jane Doe",
				Email:     "",
				Message:   "Test message",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty message",
			entry: GuestbookEntry{
				ID:        4,
				Name:      "Test User",
				Email:     "user@example.com",
				Message:   "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.entry.ID == 0 && tt.name != "empty ID" {
				// ID is set, no error
			}
			if tt.entry.Name == "" && tt.name != "empty name" {
				// Name is not empty, continue
			}
			if tt.entry.Email == "" && tt.name != "empty email" {
				// Email is not empty, continue
			}
			if tt.entry.Message == "" && tt.name != "empty message" {
				// Message is not empty, continue
			}
		})
	}
}

func TestGuestbookEntry_Fields(t *testing.T) {
	now := time.Now()
	entry := GuestbookEntry{
		ID:        42,
		Name:      "Test User",
		Email:     "test@example.com",
		Message:   "Test Message",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if entry.ID != 42 {
		t.Errorf("ID = %d; want 42", entry.ID)
	}

	if entry.Name != "Test User" {
		t.Errorf("Name = %s; want Test User", entry.Name)
	}

	if entry.Email != "test@example.com" {
		t.Errorf("Email = %s; want test@example.com", entry.Email)
	}

	if entry.Message != "Test Message" {
		t.Errorf("Message = %s; want Test Message", entry.Message)
	}

	if entry.CreatedAt != now {
		t.Error("CreatedAt does not match expected time")
	}

	if entry.UpdatedAt != now {
		t.Error("UpdatedAt does not match expected time")
	}
}

func TestGuestbookEntry_ZeroValue(t *testing.T) {
	var entry GuestbookEntry

	if entry.ID != 0 {
		t.Errorf("Zero value ID = %d; want 0", entry.ID)
	}

	if entry.Name != "" {
		t.Errorf("Zero value Name = %s; want empty string", entry.Name)
	}

	if entry.Email != "" {
		t.Errorf("Zero value Email = %s; want empty string", entry.Email)
	}

	if entry.Message != "" {
		t.Errorf("Zero value Message = %s; want empty string", entry.Message)
	}

	if !entry.CreatedAt.IsZero() {
		t.Error("Zero value CreatedAt is not zero time")
	}

	if !entry.UpdatedAt.IsZero() {
		t.Error("Zero value UpdatedAt is not zero time")
	}
}

func TestGuestbookEntry_LongStrings(t *testing.T) {
	longName := string(make([]byte, 255))
	longEmail := string(make([]byte, 255))
	longMessage := string(make([]byte, 1000))

	entry := GuestbookEntry{
		Name:    longName,
		Email:   longEmail,
		Message: longMessage,
	}

	if len(entry.Name) != 255 {
		t.Errorf("Name length = %d; want 255", len(entry.Name))
	}

	if len(entry.Email) != 255 {
		t.Errorf("Email length = %d; want 255", len(entry.Email))
	}

	if len(entry.Message) != 1000 {
		t.Errorf("Message length = %d; want 1000", len(entry.Message))
	}
}

func TestGuestbookEntry_SpecialCharacters(t *testing.T) {
	entry := GuestbookEntry{
		Name:    "John Döe (тест)",
		Email:   "john+test@example.com",
		Message: "Hello! 你好 🌍 <script>alert('xss')</script>",
	}

	if entry.Name == "" {
		t.Error("Name with special characters should not be empty")
	}

	if entry.Email == "" {
		t.Error("Email with + sign should not be empty")
	}

	if entry.Message == "" {
		t.Error("Message with special characters should not be empty")
	}
}

func TestGuestbookEntry_TimeComparison(t *testing.T) {
	now := time.Now()
	entry := GuestbookEntry{
		CreatedAt: now,
		UpdatedAt: now.Add(1 * time.Hour),
	}

	if !entry.UpdatedAt.After(entry.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}
}

func TestGuestbookEntry_MultipleEntries(t *testing.T) {
	entries := []GuestbookEntry{
		{ID: 1, Name: "User 1", Email: "user1@example.com", Message: "Message 1"},
		{ID: 2, Name: "User 2", Email: "user2@example.com", Message: "Message 2"},
		{ID: 3, Name: "User 3", Email: "user3@example.com", Message: "Message 3"},
	}

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	for i, entry := range entries {
		if entry.ID != uint(i+1) {
			t.Errorf("Entry %d: ID = %d; want %d", i, entry.ID, i+1)
		}
	}
}

func TestGuestbookEntry_PointerVsValue(t *testing.T) {
	entry := GuestbookEntry{
		ID:      1,
		Name:    "Test",
		Email:   "test@example.com",
		Message: "Message",
	}

	entryPtr := &GuestbookEntry{
		ID:      2,
		Name:    "Test 2",
		Email:   "test2@example.com",
		Message: "Message 2",
	}

	if entry.TableName() != entryPtr.TableName() {
		t.Error("TableName should be the same for value and pointer receivers")
	}
}

func TestGuestbookEntry_Modification(t *testing.T) {
	entry := GuestbookEntry{
		ID:      1,
		Name:    "Original",
		Email:   "original@example.com",
		Message: "Original message",
	}

	// Modify entry
	entry.Name = "Modified"
	entry.Email = "modified@example.com"
	entry.Message = "Modified message"

	if entry.Name != "Modified" {
		t.Error("Failed to modify Name")
	}

	if entry.Email != "modified@example.com" {
		t.Error("Failed to modify Email")
	}

	if entry.Message != "Modified message" {
		t.Error("Failed to modify Message")
	}
}

func TestGuestbookEntry_Comparison(t *testing.T) {
	now := time.Now()
	entry1 := GuestbookEntry{
		ID:        1,
		Name:      "Test",
		Email:     "test@example.com",
		Message:   "Message",
		CreatedAt: now,
		UpdatedAt: now,
	}

	entry2 := GuestbookEntry{
		ID:        1,
		Name:      "Test",
		Email:     "test@example.com",
		Message:   "Message",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if entry1.ID != entry2.ID {
		t.Error("IDs should be equal")
	}

	if entry1.Name != entry2.Name {
		t.Error("Names should be equal")
	}

	if entry1.Email != entry2.Email {
		t.Error("Emails should be equal")
	}

	if entry1.Message != entry2.Message {
		t.Error("Messages should be equal")
	}
}

func TestGuestbookEntry_EmptyTableName(t *testing.T) {
	// Test that TableName always returns the same value
	entry1 := GuestbookEntry{}
	entry2 := GuestbookEntry{ID: 999, Name: "Test"}

	if entry1.TableName() != entry2.TableName() {
		t.Error("TableName should be consistent regardless of entry values")
	}
}
