package auth

import "context"

// ContextKey avoids collisions for context.Value keys.
type ContextKey string

const (
	// UserContextKey stores *JWTClaims in context after successful gRPC auth.
	UserContextKey ContextKey = "user_claims"
	// AuthorizationHeader is the gRPC metadata key for the bearer token (lowercase per gRPC convention).
	AuthorizationHeader = "authorization"
)

// GetClaimsFromContext returns claims set by JWTAuthInterceptor, if present.
func GetClaimsFromContext(ctx context.Context) (*JWTClaims, bool) {
	claims, ok := ctx.Value(UserContextKey).(*JWTClaims)
	return claims, ok
}
