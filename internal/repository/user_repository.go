package repository

import "github.com/novaherdi/go-proto/internal/model"

// UserRepository is the persistence port for users (Dependency Inversion Principle).
// Auth and other features depend on this abstraction rather than GORM types.
type UserRepository interface {
	// Create inserts a new user record.
	Create(user *model.User) error
	// FindByUsername returns a user by unique username or an error if not found.
	FindByUsername(username string) (*model.User, error)
	// FindByEmail returns a user by unique email or an error if not found.
	FindByEmail(email string) (*model.User, error)
	// FindByID returns a user by primary key or an error if not found.
	FindByID(id uint) (*model.User, error)
}
