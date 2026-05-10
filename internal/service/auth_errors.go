package service

import "errors"

// Sentinel errors for authentication use cases. Callers (e.g. HTTP layer) map these
// to status codes without string-matching implementation details (stable contract).
var (
	// ErrInvalidCredentials indicates username/password did not match a valid user.
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrUserAlreadyExists indicates a conflicting username or email on register.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound is returned when a user lookup by ID fails in flows that need it.
	ErrUserNotFound = errors.New("user not found")
)
