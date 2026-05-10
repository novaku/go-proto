package service

import (
	"testing"

	"github.com/novaherdi/go-proto/pkg/auth"
)

func TestNewJWTTokenIssuer_GenerateToken(t *testing.T) {
	jwtSvc := auth.NewJWTService(auth.JWTConfig{
		SecretKey:     "test-secret-key-for-unit-tests-only",
		TokenDuration: 24,
		Issuer:        "test",
	})
	issuer := NewJWTTokenIssuer(jwtSvc)

	token, err := issuer.GenerateToken(1, "user", "u@example.com")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := jwtSvc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != 1 || claims.Username != "user" || claims.Email != "u@example.com" {
		t.Errorf("claims = %+v", claims)
	}
}
