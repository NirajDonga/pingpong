package result

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	History(ctx context.Context, monitorID uuid.UUID, limit int) ([]CheckResult, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) History(ctx context.Context, monitorID string, limit int) ([]CheckResult, error) {
	mid, err := uuid.Parse(monitorID)
	if err != nil {
		return nil, errors.New("invalid monitor id")
	}

	if limit <= 0 || limit > 500 {
		limit = 100
	}

	return s.repo.History(ctx, mid, limit)
}
