// Tests for defaultGuestbookValidator: add-entry and list pagination rules.
package service

import (
	"testing"

	pb "github.com/novaherdi/go-proto/gen/guestbook/v1"
)

func TestDefaultValidator_ValidateAddEntryRequest(t *testing.T) {
	validator := NewDefaultValidator()

	testCases := []struct {
		name        string
		request     *pb.AddEntryRequest
		expectError bool
		expectedMsg string
	}{
		{
			name: "valid request",
			request: &pb.AddEntryRequest{
				Name:    "John Doe",
				Email:   "john@example.com",
				Message: "Hello, World!",
			},
			expectError: false,
		},
		{
			name: "valid request without email",
			request: &pb.AddEntryRequest{
				Name:    "John Doe",
				Message: "Hello, World!",
			},
			expectError: false,
		},
		{
			name:        "nil request",
			request:     nil,
			expectError: true,
			expectedMsg: "request cannot be nil",
		},
		{
			name: "empty name",
			request: &pb.AddEntryRequest{
				Name:    "",
				Message: "Hello, World!",
			},
			expectError: true,
			expectedMsg: "name is required",
		},
		{
			name: "empty message",
			request: &pb.AddEntryRequest{
				Name:    "John Doe",
				Message: "",
			},
			expectError: true,
			expectedMsg: "message is required",
		},
		{
			name: "whitespace only name",
			request: &pb.AddEntryRequest{
				Name:    "   ",
				Message: "Hello, World!",
			},
			expectError: true,
			expectedMsg: "name is required",
		},
		{
			name: "whitespace only message",
			request: &pb.AddEntryRequest{
				Name:    "John Doe",
				Message: "   ",
			},
			expectError: true,
			expectedMsg: "message is required",
		},
		{
			name: "name with surrounding whitespace",
			request: &pb.AddEntryRequest{
				Name:    "  John Doe  ",
				Message: "Hello, World!",
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateAddEntryRequest(tc.request)

			if tc.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if err.Error() != tc.expectedMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.expectedMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Check that whitespace was trimmed
				if tc.request != nil && tc.request.Name == "  John Doe  " {
					// Should be trimmed after validation
					if tc.request.Name != "John Doe" {
						t.Errorf("expected name to be trimmed to 'John Doe', got '%s'", tc.request.Name)
					}
				}
			}
		})
	}
}

func TestDefaultValidator_ValidateListEntriesRequest(t *testing.T) {
	validator := NewDefaultValidator()

	testCases := []struct {
		name           string
		request        *pb.ListEntriesRequest
		expectError    bool
		expectedMsg    string
		expectedLimit  int
		expectedOffset int
	}{
		{
			name: "default values",
			request: &pb.ListEntriesRequest{
				Limit:  0,
				Offset: 0,
			},
			expectError:    false,
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name: "valid custom values",
			request: &pb.ListEntriesRequest{
				Limit:  20,
				Offset: 10,
			},
			expectError:    false,
			expectedLimit:  20,
			expectedOffset: 10,
		},
		{
			name: "limit too high",
			request: &pb.ListEntriesRequest{
				Limit:  200,
				Offset: 0,
			},
			expectError:    false,
			expectedLimit:  100, // Should be capped at 100
			expectedOffset: 0,
		},
		{
			name: "negative limit",
			request: &pb.ListEntriesRequest{
				Limit:  -5,
				Offset: 0,
			},
			expectError:    false,
			expectedLimit:  10, // Should use default
			expectedOffset: 0,
		},
		{
			name: "negative offset",
			request: &pb.ListEntriesRequest{
				Limit:  10,
				Offset: -5,
			},
			expectError:    false,
			expectedLimit:  10,
			expectedOffset: 0, // Should be reset to 0
		},
		{
			name:        "nil request",
			request:     nil,
			expectError: true,
			expectedMsg: "request cannot be nil",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params, err := validator.ValidateListEntriesRequest(tc.request)

			if tc.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				} else if err.Error() != tc.expectedMsg {
					t.Errorf("expected error message '%s', got '%s'", tc.expectedMsg, err.Error())
				}
				if params != nil {
					t.Error("expected params to be nil on error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if params == nil {
					t.Fatal("expected params to be non-nil")
				}
				if params.Limit != tc.expectedLimit {
					t.Errorf("expected limit %d, got %d", tc.expectedLimit, params.Limit)
				}
				if params.Offset != tc.expectedOffset {
					t.Errorf("expected offset %d, got %d", tc.expectedOffset, params.Offset)
				}
			}
		})
	}
}
