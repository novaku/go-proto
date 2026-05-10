package repository

import (
	"github.com/novaherdi/go-proto/internal/model"
	"gorm.io/gorm"
)

// userGormRepository implements UserRepository using GORM (Single Responsibility: user persistence).
type userGormRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userGormRepository{db: db}
}

// Create creates a new user in the database.
func (r *userGormRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByUsername finds a user by their username.
func (r *userGormRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by their email.
func (r *userGormRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID finds a user by their ID.
func (r *userGormRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
