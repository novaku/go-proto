package auth

import "errors"

// Errors returned by token validation and parsing.
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)
