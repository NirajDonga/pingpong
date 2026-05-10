package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) DueMonitors(ctx context.Context, limit int) ([]DueMonitor, error) {
	query := `
		SELECT id, url, interval_seconds, timeout_seconds, expected_status, next_check_at
		FROM monitors
		WHERE enabled = true
		AND next_check_at <= NOW()
		ORDER BY next_check_at
		LIMIT $1
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.db.Query(opCtx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitors := []DueMonitor{}
	for rows.Next() {
		var m DueMonitor
		if err := rows.Scan(&m.ID, &m.URL, &m.IntervalSeconds, &m.TimeoutSeconds, &m.ExpectedStatus, &m.NextCheckAt); err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (s *Store) UpdateNextCheck(ctx context.Context, monitorID uuid.UUID, intervalSeconds int) error {
	query := `
		UPDATE monitors
		SET next_check_at = NOW() + ($2 * INTERVAL '1 second')
		WHERE id = $1
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.Exec(opCtx, query, monitorID, intervalSeconds)
	return err
}
