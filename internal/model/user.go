package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User is the authentication principal stored in the database.
type User struct {
	// ID is the primary key.
	ID uint `gorm:"primaryKey" json:"id"`
	// Username is unique and used at login.
	Username string `gorm:"uniqueIndex;size:100;not null" json:"username"`
	// Email is unique and used for registration and profile.
	Email string `gorm:"uniqueIndex;size:255;not null" json:"email"`
	// Password holds the bcrypt hash; omitted from JSON responses.
	Password string `gorm:"size:255;not null" json:"-"`
	// CreatedAt is the record creation time.
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the last update time.
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt enables soft deletes via GORM.
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the GORM table name for User.
func (User) TableName() string {
	return "users"
}

// SetPassword hashes the plaintext password and stores it on the user.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword returns true if the plaintext matches the stored bcrypt hash.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
