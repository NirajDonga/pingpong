package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	CreateUser(ctx context.Context, u User) (uuid.UUID, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, u User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var id uuid.UUID
	err := r.db.QueryRow(opCtx, query, u.ID, u.Email, u.PasswordHash).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var u User
	err := r.db.QueryRow(opCtx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}
