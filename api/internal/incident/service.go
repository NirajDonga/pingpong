package incident

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID string, limit int) ([]Incident, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.List(ctx, uid, normalizeLimit(limit))
}

func (s *Service) ListForMonitor(ctx context.Context, userID string, monitorID string, limit int) ([]Incident, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	mid, err := uuid.Parse(monitorID)
	if err != nil {
		return nil, errors.New("invalid monitor id")
	}

	return s.repo.ListForMonitor(ctx, uid, mid, normalizeLimit(limit))
}

func normalizeLimit(limit int) int {
	if limit <= 0 || limit > 500 {
		return 100
	}

	return limit
}
