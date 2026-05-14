package user

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/NirajDonga/pingpong/api/internal/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

type Service interface {
	Register(ctx context.Context, input RegisterRequest) (AuthResponse, error)
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

func (s *service) Register(ctx context.Context, input RegisterRequest) (AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return AuthResponse{}, errors.New("missing required fields")
	}

	if !emailRegex.MatchString(email) {
		return AuthResponse{}, errors.New("invalid email format")
	}

	existing, err := s.repo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return AuthResponse{}, errors.New("email already registered")
	}
	if err != nil && !errors.Is(err, ErrNotFound) {
		return AuthResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("hash password: %w", err)
	}

	u := User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
	}

	id, err := s.repo.CreateUser(ctx, u)
	if err != nil {
		return AuthResponse{}, err
	}

	token, err := s.auth.GenerateToken(id.String())
	if err != nil {
		return AuthResponse{}, fmt.Errorf("generate token: %w", err)
	}

	return AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    id.String(),
			Email: u.Email,
		},
	}, nil
}

func (s *service) Login(ctx context.Context, input LoginRequest) (AuthResponse, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" || input.Password == "" {
		return AuthResponse{}, errors.New("missing credentials")
	}

	if !emailRegex.MatchString(email) {
		return AuthResponse{}, errors.New("invalid credentials")
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
