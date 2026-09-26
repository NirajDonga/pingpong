package processor

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

type Service struct {
	tbRepo *TinybirdRepository
	pgRepo *PostgresRepository
}

func NewService(tbRepo *TinybirdRepository, pgRepo *PostgresRepository) *Service {
	return &Service{
		tbRepo: tbRepo,
		pgRepo: pgRepo,
	}
}

func (s *Service) Process(ctx context.Context, result CheckResult) error {
	if err := validate(result); err != nil {
		return err
	}

	if err := s.tbRepo.Insert(ctx, result); err != nil {
		return err
	}

	monitorID, err := uuid.Parse(result.MonitorID)
	if err != nil {
		return errors.New("invalid monitorId")
	}

	if err := s.pgRepo.ApplyCheckResult(ctx, monitorID, result.Success); err != nil {
		log.Printf("failed to apply check result to postgres for monitor %s: %v", monitorID, err)
		return err
	}

	return nil
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
