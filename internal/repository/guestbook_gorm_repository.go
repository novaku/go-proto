package repository

import (
	"context"

	"github.com/novaherdi/go-proto/internal/model"
	"gorm.io/gorm"
)

// gormGuestbookRepository implements GuestbookRepository with GORM (Single Responsibility: SQL access only).
type gormGuestbookRepository struct {
	db *gorm.DB
}

// NewGormGuestbookRepository builds a repository backed by the given DB handle.
func NewGormGuestbookRepository(db *gorm.DB) GuestbookRepository {
	return &gormGuestbookRepository{db: db}
}

// Create inserts a new guestbook entry into the database.
func (r *gormGuestbookRepository) Create(ctx context.Context, entry *model.GuestbookEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

// FindWithPagination retrieves guestbook entries with pagination, ordered by created_at DESC.
func (r *gormGuestbookRepository) FindWithPagination(ctx context.Context, limit, offset int) ([]model.GuestbookEntry, error) {
	var entries []model.GuestbookEntry
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error

	return entries, err
}
