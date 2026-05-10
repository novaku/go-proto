package model

import (
	"time"
)

// GuestbookEntry is the persisted guestbook row (domain + GORM mapping).
type GuestbookEntry struct {
	// ID is the primary key assigned by the database.
	ID uint `gorm:"primaryKey;autoIncrement"`
	// Name is the visitor display name.
	Name string `gorm:"type:varchar(255);not null"`
	// Email is stored for contact; may be empty string if optional at API level.
	Email string `gorm:"type:varchar(255);not null"`
	// Message is the guestbook body text.
	Message string `gorm:"type:text;not null"`
	// CreatedAt is set on insert (millisecond precision).
	CreatedAt time.Time `gorm:"autoCreateTime:milli;not null"`
	// UpdatedAt is maintained by GORM on updates.
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
}

// TableName returns the GORM table name for GuestbookEntry.
func (GuestbookEntry) TableName() string {
	return "guestbook_entries"
}
