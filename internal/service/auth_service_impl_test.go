package service

import (
	"errors"
	"testing"

	"github.com/novaherdi/go-proto/internal/dto"
	"github.com/novaherdi/go-proto/internal/model"
	"github.com/novaherdi/go-proto/internal/repository"
	"gorm.io/gorm"
)

type mockUserRepository struct {
	findByUsername func(string) (*model.User, error)
	findByEmail    func(string) (*model.User, error)
	create         func(*model.User) error
	findByID       func(uint) (*model.User, error)
}

func (m *mockUserRepository) Create(user *model.User) error {
	if m.create != nil {
		return m.create(user)
	}
	return nil
}

func (m *mockUserRepository) FindByUsername(username string) (*model.User, error) {
	if m.findByUsername != nil {
		return m.findByUsername(username)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) FindByEmail(email string) (*model.User, error) {
	if m.findByEmail != nil {
		return m.findByEmail(email)
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockUserRepository) FindByID(id uint) (*model.User, error) {
	if m.findByID != nil {
		return m.findByID(id)
	}
	return nil, gorm.ErrRecordNotFound
}

var _ repository.UserRepository = (*mockUserRepository)(nil)

type mockTokenIssuer struct {
	generateToken func(uint, string, string) (string, error)
}

func (m *mockTokenIssuer) GenerateToken(userID uint, username, email string) (string, error) {
	if m.generateToken != nil {
		return m.generateToken(userID, username, email)
	}
	return "mock-token", nil
}

func TestAuthService_Login_Success(t *testing.T) {
	u := &model.User{ID: 1, Username: "alice", Email: "a@b.c"}
	if err := u.SetPassword("correct"); err != nil {
		t.Fatal(err)
	}
	repo := &mockUserRepository{
		findByUsername: func(s string) (*model.User, error) {
			if s == "alice" {
				return u, nil
			}
			return nil, gorm.ErrRecordNotFound
		},
	}
	issuer := &mockTokenIssuer{
		generateToken: func(id uint, user, em string) (string, error) {
			if id != 1 || user != "alice" {
				t.Errorf("unexpected args %d %s", id, user)
			}
			return "tok", nil
		},
	}
	svc := NewAuthService(repo, issuer)

	out, err := svc.Login(&dto.LoginRequest{Username: "alice", Password: "correct"})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if out.Token != "tok" || out.UserID != 1 || out.Username != "alice" {
		t.Errorf("response %+v", out)
	}
}

func TestAuthService_Login_InvalidCredentials_NotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Login(&dto.LoginRequest{Username: "nope", Password: "x"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_InvalidCredentials_BadPassword(t *testing.T) {
	u := &model.User{ID: 1, Username: "alice", Email: "a@b.c"}
	if err := u.SetPassword("right"); err != nil {
		t.Fatal(err)
	}
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return u, nil },
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Login(&dto.LoginRequest{Username: "alice", Password: "wrong"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_TokenIssuerError(t *testing.T) {
	u := &model.User{ID: 1, Username: "alice", Email: "a@b.c"}
	if err := u.SetPassword("ok"); err != nil {
		t.Fatal(err)
	}
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return u, nil },
	}
	tokenErr := errors.New("signing failed")
	issuer := &mockTokenIssuer{
		generateToken: func(uint, string, string) (string, error) { return "", tokenErr },
	}
	svc := NewAuthService(repo, issuer)

	_, err := svc.Login(&dto.LoginRequest{Username: "alice", Password: "ok"})
	if !errors.Is(err, tokenErr) {
		t.Fatalf("got %v", err)
	}
}

func TestAuthService_Login_RepoError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return nil, dbErr },
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Login(&dto.LoginRequest{Username: "a", Password: "b"})
	if !errors.Is(err, dbErr) {
		t.Fatalf("want db error, got %v", err)
	}
}

func TestAuthService_Register_Success(t *testing.T) {
	var created *model.User
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		findByEmail:    func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		create: func(user *model.User) error {
			user.ID = 99
			created = user
			return nil
		},
	}
	issuer := &mockTokenIssuer{
		generateToken: func(id uint, username, email string) (string, error) {
			if id != 99 {
				t.Errorf("user id %d", id)
			}
			return "reg-token", nil
		},
	}
	svc := NewAuthService(repo, issuer)

	out, err := svc.Register(&dto.RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "secret12",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if out.Token != "reg-token" || out.UserID != 99 {
		t.Errorf("response %+v", out)
	}
	if created == nil || created.Username != "newuser" {
		t.Fatalf("user not created correctly: %+v", created)
	}
}

func TestAuthService_Register_UserAlreadyExists_Username(t *testing.T) {
	existing := &model.User{Username: "taken"}
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return existing, nil },
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Register(&dto.RegisterRequest{Username: "taken", Email: "e@e.com", Password: "123456"})
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("got %v", err)
	}
}

func TestAuthService_Register_UserAlreadyExists_Email(t *testing.T) {
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		findByEmail: func(string) (*model.User, error) {
			return &model.User{Email: "dup@example.com"}, nil
		},
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Register(&dto.RegisterRequest{Username: "u", Email: "dup@example.com", Password: "123456"})
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("got %v", err)
	}
}

func TestAuthService_Register_FindByUsername_Error(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return nil, dbErr },
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Register(&dto.RegisterRequest{Username: "u", Email: "e@e.com", Password: "123456"})
	if !errors.Is(err, dbErr) {
		t.Fatalf("got %v", err)
	}
}

func TestAuthService_Register_CreateError(t *testing.T) {
	repo := &mockUserRepository{
		findByUsername: func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		findByEmail:    func(string) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
		create:         func(*model.User) error { return errors.New("insert failed") },
	}
	svc := NewAuthService(repo, &mockTokenIssuer{})

	_, err := svc.Register(&dto.RegisterRequest{Username: "u", Email: "e@e.com", Password: "123456"})
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("got %v", err)
	}
}
