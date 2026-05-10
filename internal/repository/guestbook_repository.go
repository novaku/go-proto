package repository

import (
	"context"

	"github.com/novaherdi/go-proto/internal/model"
)

// GuestbookRepository is the persistence port for guestbook entries (Dependency Inversion).
// Domain and application layers depend on this interface; GORM (or another store) implements it.
type GuestbookRepository interface {
	// Create persists a new guestbook entry.
	Create(ctx context.Context, entry *model.GuestbookEntry) error
	// FindWithPagination returns a page of entries (newest first, per implementation).
	FindWithPagination(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error)
}
