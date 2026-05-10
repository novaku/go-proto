package auth

import (
	"testing"
	"time"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1,
		Issuer:        "test-issuer",
	}
	jwtService := NewJWTService(config)

	userID := uint(1)
	username := "testuser"
	email := "test@example.com"

	token, err := jwtService.GenerateToken(userID, username, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Generated token is empty")
	}

	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Username != username {
		t.Errorf("Expected Username %s, got %s", username, claims.Username)
	}
	if claims.Email != email {
		t.Errorf("Expected Email %s, got %s", email, claims.Email)
	}
}

func TestJWTService_InvalidToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1,
		Issuer:        "test-issuer",
	}
	jwtService := NewJWTService(config)

	_, err := jwtService.ValidateToken("invalid-token")
	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTService_WrongSecretKey(t *testing.T) {
	jwtService1 := NewJWTService(JWTConfig{
		SecretKey:     "secret-key-1",
		TokenDuration: 1,
		Issuer:        "test-issuer",
	})
	jwtService2 := NewJWTService(JWTConfig{
		SecretKey:     "secret-key-2",
		TokenDuration: 1,
		Issuer:        "test-issuer",
	})

	token, err := jwtService1.GenerateToken(1, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	_, err = jwtService2.ValidateToken(token)
	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken, got %v", err)
	}
}

func TestJWTService_RefreshToken(t *testing.T) {
	config := JWTConfig{
		SecretKey:     "test-secret-key",
		TokenDuration: 1,
		Issuer:        "test-issuer",
	}
	jwtService := NewJWTService(config)

	originalToken, err := jwtService.GenerateToken(1, "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Wait long enough to ensure different IssuedAt timestamp (JWT uses seconds)
	time.Sleep(1100 * time.Millisecond)

	newToken, err := jwtService.RefreshToken(originalToken)
	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	// Validate the new token works
	claims, err := jwtService.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("Failed to validate refreshed token: %v", err)
	}

	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}
}
