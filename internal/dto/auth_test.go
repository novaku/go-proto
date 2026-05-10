package dto

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestAuthDTOs_JSONRoundTrip ensures HTTP DTOs marshal and unmarshal consistently.
func TestAuthDTOs_JSONRoundTrip(t *testing.T) {
	t.Run("LoginRequest", func(t *testing.T) {
		orig := LoginRequest{Username: "alice", Password: "secret"}
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		var got LoginRequest
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !reflect.DeepEqual(orig, got) {
			t.Errorf("got %+v, want %+v", got, orig)
		}
	})

	t.Run("RegisterRequest", func(t *testing.T) {
		orig := RegisterRequest{Username: "bob", Email: "bob@example.com", Password: "hunter2!"}
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		var got RegisterRequest
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !reflect.DeepEqual(orig, got) {
			t.Errorf("got %+v, want %+v", got, orig)
		}
	})

	t.Run("AuthResponse", func(t *testing.T) {
		orig := AuthResponse{
			Token:    "jwt-here",
			UserID:   42,
			Username: "carol",
			Email:    "carol@example.com",
		}
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		var got AuthResponse
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if !reflect.DeepEqual(orig, got) {
			t.Errorf("got %+v, want %+v", got, orig)
		}
	})
}
