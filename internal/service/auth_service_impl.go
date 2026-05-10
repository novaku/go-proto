package service

import (
	"errors"

	"github.com/novaherdi/go-proto/internal/dto"
	"github.com/novaherdi/go-proto/internal/model"
	"github.com/novaherdi/go-proto/internal/repository"
	"gorm.io/gorm"
)

// authServiceImpl implements AuthService using a user store and token issuer.
// Single Responsibility: orchestrate login/register; persistence and signing are delegated.
type authServiceImpl struct {
	userRepo    repository.UserRepository
	tokenIssuer TokenIssuer
}

// NewAuthService wires authentication with injected dependencies (Dependency Inversion).
func NewAuthService(userRepo repository.UserRepository, issuer TokenIssuer) AuthService {
	return &authServiceImpl{
		userRepo:    userRepo,
		tokenIssuer: issuer,
	}
}

// Login authenticates a user and returns a JWT token.
func (s *authServiceImpl) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenIssuer.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

// Register creates a new user after uniqueness checks, then issues a token.
func (s *authServiceImpl) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, ErrUserAlreadyExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.tokenIssuer.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
