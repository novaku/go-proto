package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/novaherdi/go-proto/internal/dto"
	"github.com/novaherdi/go-proto/internal/service"
)

type mockAuthService struct {
	loginFunc    func(*dto.LoginRequest) (*dto.AuthResponse, error)
	registerFunc func(*dto.RegisterRequest) (*dto.AuthResponse, error)
}

func (m *mockAuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	if m.loginFunc != nil {
		return m.loginFunc(req)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if m.registerFunc != nil {
		return m.registerFunc(req)
	}
	return nil, errors.New("not implemented")
}

var _ service.AuthService = (*mockAuthService)(nil)

func TestAuthController_LoginHandler(t *testing.T) {
	t.Run("method not allowed", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{})
		req := httptest.NewRequest(http.MethodGet, "/auth/login", nil)
		rec := httptest.NewRecorder()
		c.LoginHandler(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("code %d", rec.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("not-json"))
		rec := httptest.NewRecorder()
		c.LoginHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("code %d", rec.Code)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{})
		body, _ := json.Marshal(dto.LoginRequest{Username: ""})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		c.LoginHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("code %d", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{
			loginFunc: func(r *dto.LoginRequest) (*dto.AuthResponse, error) {
				return &dto.AuthResponse{Token: "t", UserID: 1, Username: r.Username}, nil
			},
		})
		body, _ := json.Marshal(dto.LoginRequest{Username: "a", Password: "b"})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		c.LoginHandler(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("code %d body %s", rec.Code, rec.Body.String())
		}
		var out dto.AuthResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		if out.Token != "t" || out.Username != "a" {
			t.Errorf("%+v", out)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{
			loginFunc: func(*dto.LoginRequest) (*dto.AuthResponse, error) {
				return nil, service.ErrInvalidCredentials
			},
		})
		body, _ := json.Marshal(dto.LoginRequest{Username: "a", Password: "b"})
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		c.LoginHandler(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("code %d", rec.Code)
		}
	})
}

func TestAuthController_RegisterHandler(t *testing.T) {
	t.Run("password too short", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{})
		body, _ := json.Marshal(dto.RegisterRequest{Username: "u", Email: "e@e.com", Password: "12345"})
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		c.RegisterHandler(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("code %d", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{
			registerFunc: func(r *dto.RegisterRequest) (*dto.AuthResponse, error) {
				return &dto.AuthResponse{Token: "rt", UserID: 2, Username: r.Username, Email: r.Email}, nil
			},
		})
		body, _ := json.Marshal(dto.RegisterRequest{Username: "u", Email: "u@u.com", Password: "123456"})
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		c.RegisterHandler(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("code %d %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("conflict", func(t *testing.T) {
		c := NewAuthController(&mockAuthService{
			registerFunc: func(*dto.RegisterRequest) (*dto.AuthResponse, error) {
				return nil, service.ErrUserAlreadyExists
			},
		})
		body, _ := json.Marshal(dto.RegisterRequest{Username: "u", Email: "u@u.com", Password: "123456"})
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
		rec := httptest.NewRecorder()
		c.RegisterHandler(rec, req)
		if rec.Code != http.StatusConflict {
			t.Errorf("code %d", rec.Code)
		}
	})
}
