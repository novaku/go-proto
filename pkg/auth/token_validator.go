package auth

// TokenValidator abstracts JWT validation (Dependency Inversion). Interceptors and
// middleware depend on this interface; *JWTService is the typical implementation.
type TokenValidator interface {
	ValidateToken(tokenString string) (*JWTClaims, error)
}
