package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/NirajDonga/pingpong/api/internal/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, input RegisterRequest) (UserResponse, error)
	Login(ctx context.Context, input LoginRequest) (AuthResponse, error)
}

type service struct {
	repo Repository
	auth auth.Service
}

func NewService(repo Repository, authSvc auth.Service) Service {
	return &service{
		repo: repo,
		auth: authSvc,
	}
}

func (s *service) Register(ctx context.Context, input RegisterRequest) (UserResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return UserResponse{}, errors.New("missing required fields")
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return UserResponse{}, errors.New("email already registered")
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("hash password: %w", err)
	}

	u := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	id, err := s.repo.CreateUser(ctx, u)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID:    id.String(),
		Email: u.Email,
	}, nil
}

func (s *service) Login(ctx context.Context, input LoginRequest) (AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return AuthResponse{}, errors.New("missing credentials")
	}

	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return AuthResponse{}, errors.New("invalid credentials")
		}
		return AuthResponse{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password)) != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	token, err := s.auth.GenerateToken(u.ID.String())
	if err != nil {
		return AuthResponse{}, fmt.Errorf("generate token: %w", err)
	}

	return AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    u.ID.String(),
			Email: u.Email,
		},
	}, nil
}
