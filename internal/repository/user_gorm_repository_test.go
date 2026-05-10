package repository

import (
	"errors"
	"testing"

	"github.com/novaherdi/go-proto/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestUserGormRepository_CreateAndFind(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	u := &model.User{Username: "alice", Email: "alice@example.com"}
	if err := u.SetPassword("password123"); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u.ID == 0 {
		t.Fatal("expected ID after create")
	}

	byName, err := repo.FindByUsername("alice")
	if err != nil {
		t.Fatalf("FindByUsername: %v", err)
	}
	if byName.Email != "alice@example.com" {
		t.Errorf("email %s", byName.Email)
	}
	if !byName.CheckPassword("password123") {
		t.Error("password should verify")
	}

	byEmail, err := repo.FindByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("FindByEmail: %v", err)
	}
	if byEmail.Username != "alice" {
		t.Errorf("username %s", byEmail.Username)
	}

	byID, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if byID.Username != "alice" {
		t.Errorf("FindByID %+v", byID)
	}
}

func TestUserGormRepository_FindByUsername_NotFound(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByUsername("nobody")
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("want record not found, got %v", err)
	}
}
