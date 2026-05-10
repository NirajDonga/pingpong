package monitor

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	UserID               uuid.UUID `json:"user_id" db:"user_id"`
	Name                 string    `json:"name" db:"name"`
	URL                  string    `json:"url" db:"url"`
	IntervalSeconds      int       `json:"interval_seconds" db:"interval_seconds"`
	TimeoutSeconds       int       `json:"timeout_seconds" db:"timeout_seconds"`
	ExpectedStatus       int       `json:"expected_status" db:"expected_status"`
	Enabled              bool      `json:"enabled" db:"enabled"`
	CurrentStatus        string    `json:"current_status" db:"current_status"`
	ConsecutiveFailures  int       `json:"consecutive_failures" db:"consecutive_failures"`
	ConsecutiveSuccesses int       `json:"consecutive_successes" db:"consecutive_successes"`
	NextCheckAt          time.Time `json:"next_check_at" db:"next_check_at"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

type CreateRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	ExpectedStatus  int    `json:"expected_status"`
}

type UpdateRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
	ExpectedStatus  int    `json:"expected_status"`
}

type EnabledRequest struct {
	Enabled bool `json:"enabled"`
}
