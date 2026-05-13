package incident

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	MonitorID uuid.UUID  `json:"monitor_id" db:"monitor_id"`
	StartedAt time.Time  `json:"started_at" db:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	Status    string     `json:"status" db:"status"`
	Reason    string     `json:"reason" db:"reason"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
