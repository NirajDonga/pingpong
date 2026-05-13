package incident

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const incidentColumns = `
	i.id, i.monitor_id, i.started_at, i.ended_at, i.status, i.reason, i.created_at
`

type Repository interface {
	List(ctx context.Context, userID uuid.UUID, limit int) ([]Incident, error)
	ListForMonitor(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID, limit int) ([]Incident, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) List(ctx context.Context, userID uuid.UUID, limit int) ([]Incident, error) {
	query := `
		SELECT ` + incidentColumns + `
		FROM incidents i
		INNER JOIN monitors m ON m.id = i.monitor_id
		WHERE m.user_id = $1
		ORDER BY i.started_at DESC
		LIMIT $2
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := r.db.Query(opCtx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanIncidents(rows)
}

func (r *repository) ListForMonitor(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID, limit int) ([]Incident, error) {
	query := `
		SELECT ` + incidentColumns + `
		FROM incidents i
		INNER JOIN monitors m ON m.id = i.monitor_id
		WHERE m.user_id = $1 AND m.id = $2
		ORDER BY i.started_at DESC
		LIMIT $3
	`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := r.db.Query(opCtx, query, userID, monitorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanIncidents(rows)
}

type incidentRows interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
}

func scanIncidents(rows incidentRows) ([]Incident, error) {
	incidents := []Incident{}
	for rows.Next() {
		var i Incident
		if err := rows.Scan(
			&i.ID,
			&i.MonitorID,
			&i.StartedAt,
			&i.EndedAt,
			&i.Status,
			&i.Reason,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}
