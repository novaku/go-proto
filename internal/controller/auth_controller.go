package controller

import (
	"encoding/json"
	"net/http"

	"github.com/novaherdi/go-proto/internal/dto"
	"github.com/novaherdi/go-proto/internal/service"
)

// AuthController handles HTTP authentication routes. It maps HTTP to AuthService calls only
// (Single Responsibility); business rules live in service, tokens in TokenIssuer.
type AuthController struct {
	authService service.AuthService
}

// NewAuthController constructs an AuthController with an injected AuthService.
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// LoginHandler handles POST /auth/login.
func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	resp, err := c.authService.Login(&req)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON(w, http.StatusOK, resp)
}

// RegisterHandler handles POST /auth/register.
func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		RespondWithError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	if len(req.Password) < 6 {
		RespondWithError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	resp, err := c.authService.Register(&req)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			RespondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, resp)
}
