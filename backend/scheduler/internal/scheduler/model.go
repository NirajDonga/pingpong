package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type DueMonitor struct {
	ID              uuid.UUID
	URL             string
	IntervalSeconds int
	TimeoutSeconds  int
	ExpectedStatus  int
	NextCheckAt     time.Time
}

type CheckJob struct {
	MonitorID      string `json:"monitorId"`
	URL            string `json:"url"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
	ExpectedStatus int    `json:"expectedStatus"`
}
