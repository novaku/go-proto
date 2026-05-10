// Package dto holds transport-layer data shapes (HTTP JSON, etc.) separate from
// domain models and persistence. This supports Single Responsibility: API payloads
// can evolve without entangling database or protobuf definitions.
package dto

// LoginRequest is the JSON body for POST /auth/login.
// Validation rules live in the controller or a dedicated validator, not on this struct.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest is the JSON body for POST /auth/register.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is returned after successful login or registration (token + profile fields).
type AuthResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
