package controller

import (
	"encoding/json"
	"net/http"
)

// RespondWithError writes a JSON error envelope { "error": message } with the given status.
// Single Responsibility: HTTP serialization for failures, reusable across handlers.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON marshals payload to JSON, sets Content-Type, status, and body.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
