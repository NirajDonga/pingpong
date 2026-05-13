package result

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Repository interface {
	Insert(ctx context.Context, result CheckResult) error
	History(ctx context.Context, monitorID uuid.UUID, limit int) ([]CheckResult, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Process(ctx context.Context, result CheckResult) error {
	if err := validate(result); err != nil {
		return err
	}

	return s.repo.Insert(ctx, result)
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

func validate(result CheckResult) error {
	if _, err := uuid.Parse(result.MonitorID); err != nil {
		return errors.New("invalid monitorId")
	}
	if result.CheckedAt.IsZero() {
		return errors.New("checkedAt is required")
	}
	if result.StatusCode < 0 || result.StatusCode > 599 {
		return errors.New("statusCode must be a valid HTTP status code")
	}
	if result.ResponseTimeMS < 0 {
		return errors.New("responseTimeMs must be non-negative")
	}
	if result.DNSMS < 0 || result.TCPMS < 0 || result.TLSMS < 0 || result.TTFBMS < 0 {
		return errors.New("timing fields must be non-negative")
	}

	return nil
}
