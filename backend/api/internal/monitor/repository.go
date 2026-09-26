package monitor

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("monitor not found")

const monitorColumns = `
	id, user_id, name, url, interval_seconds, timeout_seconds,
	expected_status, enabled, current_status, consecutive_failures,
	consecutive_successes, next_check_at, created_at
`

type Repository interface {
	Create(ctx context.Context, m Monitor) (*Monitor, error)
	List(ctx context.Context, userID uuid.UUID) ([]Monitor, error)
	Get(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID) (*Monitor, error)
	Update(ctx context.Context, m Monitor) (*Monitor, error)
	SetEnabled(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID, enabled bool) (*Monitor, error)
	Delete(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID) error
	ApplyCheckResult(ctx context.Context, monitorID uuid.UUID, success bool) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, m Monitor) (*Monitor, error) {
	query := `
		INSERT INTO monitors (
			id, user_id, name, url, interval_seconds,
			timeout_seconds, expected_status, next_check_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING ` + monitorColumns

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var out Monitor
	row := r.db.QueryRow(opCtx, query,
		m.ID, m.UserID, m.Name, m.URL,
		m.IntervalSeconds, m.TimeoutSeconds, m.ExpectedStatus, m.NextCheckAt,
	)
	err := scanMonitor(row, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (r *repository) List(ctx context.Context, userID uuid.UUID) ([]Monitor, error) {
	query := `
		SELECT ` + monitorColumns + `
		FROM monitors
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := r.db.Query(opCtx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitors := []Monitor{}
	for rows.Next() {
		var m Monitor
		if err := scanMonitor(rows, &m); err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monitors, nil
}

func (r *repository) Get(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID) (*Monitor, error) {
	query := `
		SELECT ` + monitorColumns + `
		FROM monitors
		WHERE user_id = $1 AND id = $2
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var m Monitor
	row := r.db.QueryRow(opCtx, query, userID, monitorID)
	err := scanMonitor(row, &m)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *repository) Update(ctx context.Context, m Monitor) (*Monitor, error) {
	query := `
		UPDATE monitors
		SET name = $3,
			url = $4,
			interval_seconds = $5,
			timeout_seconds = $6,
			expected_status = $7
		WHERE user_id = $1 AND id = $2
		RETURNING ` + monitorColumns

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var out Monitor
	row := r.db.QueryRow(
		opCtx,
		query,
		m.UserID,
		m.ID,
		m.Name,
		m.URL,
		m.IntervalSeconds,
		m.TimeoutSeconds,
		m.ExpectedStatus,
	)
	err := scanMonitor(row, &out)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &out, nil
}

func (r *repository) SetEnabled(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID, enabled bool) (*Monitor, error) {
	query := `
		UPDATE monitors
		SET enabled = $3
		WHERE user_id = $1 AND id = $2
		RETURNING ` + monitorColumns

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var m Monitor
	row := r.db.QueryRow(opCtx, query, userID, monitorID, enabled)
	err := scanMonitor(row, &m)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &m, nil
}

func (r *repository) Delete(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID) error {
	query := `DELETE FROM monitors WHERE user_id = $1 AND id = $2`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tag, err := r.db.Exec(opCtx, query, userID, monitorID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *repository) ApplyCheckResult(ctx context.Context, monitorID uuid.UUID, success bool) error {
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

type monitorScanner interface {
	Scan(dest ...any) error
}

func scanMonitor(row monitorScanner, m *Monitor) error {
	return row.Scan(
		&m.ID,
		&m.UserID,
		&m.Name,
		&m.URL,
		&m.IntervalSeconds,
		&m.TimeoutSeconds,
		&m.ExpectedStatus,
		&m.Enabled,
		&m.CurrentStatus,
		&m.ConsecutiveFailures,
		&m.ConsecutiveSuccesses,
		&m.NextCheckAt,
		&m.CreatedAt,
	)
}
