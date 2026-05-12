package monitor

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID string, input CreateRequest) (*Monitor, error)
	List(ctx context.Context, userID string) ([]Monitor, error)
	Get(ctx context.Context, userID string, monitorID string) (*Monitor, error)
	Update(ctx context.Context, userID string, monitorID string, input UpdateRequest) (*Monitor, error)
	SetEnabled(ctx context.Context, userID string, monitorID string, enabled bool) (*Monitor, error)
	Delete(ctx context.Context, userID string, monitorID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID string, input CreateRequest) (*Monitor, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	if err := validateInput(input.Name, input.URL, input.IntervalSeconds, input.TimeoutSeconds, input.ExpectedStatus); err != nil {
		return nil, err
	}
	expectedStatus := input.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	m := Monitor{
		ID:              uuid.New(),
		UserID:          uid,
		Name:            strings.TrimSpace(input.Name),
		URL:             strings.TrimSpace(input.URL),
		IntervalSeconds: input.IntervalSeconds,
		TimeoutSeconds:  input.TimeoutSeconds,
		ExpectedStatus:  expectedStatus,
		Enabled:         true,
		CurrentStatus:   "unknown",
		NextCheckAt:     time.Now(),
	}

	return s.repo.Create(ctx, m)
}

func (s *service) List(ctx context.Context, userID string) ([]Monitor, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	return s.repo.List(ctx, uid)
}

func (s *service) Get(ctx context.Context, userID string, monitorID string) (*Monitor, error) {
	uid, mid, err := parseIDs(userID, monitorID)
	if err != nil {
		return nil, err
	}

	return s.repo.Get(ctx, uid, mid)
}

func (s *service) Update(ctx context.Context, userID string, monitorID string, input UpdateRequest) (*Monitor, error) {
	uid, mid, err := parseIDs(userID, monitorID)
	if err != nil {
		return nil, err
	}
	if err := validateInput(input.Name, input.URL, input.IntervalSeconds, input.TimeoutSeconds, input.ExpectedStatus); err != nil {
		return nil, err
	}
	expectedStatus := input.ExpectedStatus
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	m := Monitor{
		ID:              mid,
		UserID:          uid,
		Name:            strings.TrimSpace(input.Name),
		URL:             strings.TrimSpace(input.URL),
		IntervalSeconds: input.IntervalSeconds,
		TimeoutSeconds:  input.TimeoutSeconds,
		ExpectedStatus:  expectedStatus,
	}

	return s.repo.Update(ctx, m)
}

func (s *service) SetEnabled(ctx context.Context, userID string, monitorID string, enabled bool) (*Monitor, error) {
	uid, mid, err := parseIDs(userID, monitorID)
	if err != nil {
		return nil, err
	}

	return s.repo.SetEnabled(ctx, uid, mid, enabled)
}

func (s *service) Delete(ctx context.Context, userID string, monitorID string) error {
	uid, mid, err := parseIDs(userID, monitorID)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, uid, mid)
}

func parseIDs(userID string, monitorID string) (uuid.UUID, uuid.UUID, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid user id")
	}

	mid, err := uuid.Parse(monitorID)
	if err != nil {
		return uuid.Nil, uuid.Nil, errors.New("invalid monitor id")
	}

	return uid, mid, nil
}

func validateInput(name string, rawURL string, intervalSeconds int, timeoutSeconds int, expectedStatus int) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	if !isMonitorURLAllowed(rawURL) {
		return errors.New("valid public http or https url is required")
	}
	if intervalSeconds < 30 {
		return errors.New("interval_seconds must be at least 30")
	}
	if timeoutSeconds < 1 {
		return errors.New("timeout_seconds must be at least 1")
	}
	if expectedStatus == 0 {
		expectedStatus = 200
	}
	if expectedStatus < 100 || expectedStatus > 599 {
		return errors.New("expected_status must be a valid HTTP status code")
	}

	return nil
}

func isMonitorURLAllowed(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return false
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "" {
		return false
	}

	if ip := net.ParseIP(host); ip != nil {
		return !ip.IsLoopback() &&
			!ip.IsPrivate() &&
			!ip.IsLinkLocalUnicast() &&
			!ip.IsLinkLocalMulticast() &&
			!ip.IsMulticast() &&
			!ip.IsUnspecified()
	}

	return host != "localhost" &&
		host != "127.0.0.1" &&
		host != "0.0.0.0" &&
		host != "::1" &&
		!strings.HasSuffix(host, ".svc") &&
		!strings.HasSuffix(host, ".cluster.local")
}
