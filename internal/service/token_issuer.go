package service

import "github.com/novaherdi/go-proto/pkg/auth"

// TokenIssuer abstracts JWT (or any token) creation. Auth business logic depends on
// this interface, not on pkg/auth.JWTService directly (Dependency Inversion Principle).
type TokenIssuer interface {
	// GenerateToken produces a signed token for the authenticated principal.
	GenerateToken(userID uint, username, email string) (string, error)
}

// jwtTokenIssuer adapts pkg/auth.JWTService to TokenIssuer without leaking the
// concrete type into authServiceImpl (Liskov substitution: any TokenIssuer works).
type jwtTokenIssuer struct {
	jwt *auth.JWTService
}

// NewJWTTokenIssuer wraps the shared JWT service as a TokenIssuer for dependency injection.
func NewJWTTokenIssuer(jwt *auth.JWTService) TokenIssuer {
	return &jwtTokenIssuer{jwt: jwt}
}

// GenerateToken delegates to the underlying JWT service.
func (a *jwtTokenIssuer) GenerateToken(userID uint, username, email string) (string, error) {
	return a.jwt.GenerateToken(userID, username, email)
}
