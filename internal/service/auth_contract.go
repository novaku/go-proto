package service

import "github.com/novaherdi/go-proto/internal/dto"

// AuthService defines authentication operations independent of HTTP or gRPC.
// Implementations depend on UserRepository and TokenIssuer (abstractions).
type AuthService interface {
	// Login validates credentials and returns a token plus user profile fields.
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	// Register creates a new user and returns a token plus profile fields.
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
}
