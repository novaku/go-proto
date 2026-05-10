package auth

import "github.com/golang-jwt/jwt/v5"

// JWTClaims holds custom and standard JWT claims for this application.
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWTConfig configures signing and lifetime for issued tokens.
type JWTConfig struct {
	SecretKey     string `json:"secret_key"`
	TokenDuration int    `json:"token_duration"` // hours
	Issuer        string `json:"issuer"`
}
