package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	RespondWithJSON(rec, http.StatusTeapot, map[string]int{"x": 1})

	if rec.Code != http.StatusTeapot {
		t.Errorf("status %d, want %d", rec.Code, http.StatusTeapot)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	var body map[string]int
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["x"] != 1 {
		t.Errorf("body %+v", body)
	}
}

func TestRespondWithError(t *testing.T) {
	rec := httptest.NewRecorder()
	RespondWithError(rec, http.StatusBadRequest, "bad thing")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status %d", rec.Code)
	}
	var env struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if env.Error != "bad thing" {
		t.Errorf("error = %q", env.Error)
	}
}
