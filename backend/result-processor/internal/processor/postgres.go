package processor

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("monitor not found")

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ApplyCheckResult(ctx context.Context, monitorID uuid.UUID, success bool) error {
	updateQuery := `
		UPDATE monitors
		SET consecutive_failures = CASE
				WHEN $2 THEN 0
				ELSE consecutive_failures + 1
			END,
			consecutive_successes = CASE
				WHEN $2 THEN consecutive_successes + 1
				ELSE 0
			END,
			current_status = CASE
				WHEN $2 AND consecutive_successes + 1 >= 2 THEN 'up'
				WHEN NOT $2 AND consecutive_failures + 1 >= 3 THEN 'down'
				ELSE current_status
			END
		WHERE id = $1
		RETURNING current_status
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.Begin(opCtx)
	if err != nil {
		return err
	}
	defer tx.Rollback(opCtx)

	var previousStatus string
	err = tx.QueryRow(opCtx, `SELECT current_status FROM monitors WHERE id = $1 FOR UPDATE`, monitorID).Scan(&previousStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	var currentStatus string
	err = tx.QueryRow(opCtx, updateQuery, monitorID, success).Scan(&currentStatus)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	if previousStatus != "down" && currentStatus == "down" {
		_, err = tx.Exec(
			opCtx,
			`
				INSERT INTO incidents (id, monitor_id, started_at, status, reason)
				VALUES ($1, $2, NOW(), 'open', 'monitor marked down after consecutive failures')
			`,
			uuid.New(),
			monitorID,
		)
		if err != nil {
			return err
		}
	}

	if previousStatus == "down" && currentStatus == "up" {
		_, err = tx.Exec(
			opCtx,
			`
				UPDATE incidents
				SET ended_at = NOW(),
					status = 'resolved'
				WHERE monitor_id = $1
					AND status = 'open'
					AND ended_at IS NULL
			`,
			monitorID,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(opCtx); err != nil {
		return err
	}

	return nil
}
