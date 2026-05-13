package tests

import (
	"context"
	"testing"

	"github.com/NirajDonga/pingpong/api/internal/incident"
	"github.com/google/uuid"
)

func TestIncidentListNormalizesLimit(t *testing.T) {
	repo := &fakeIncidentRepository{}
	svc := incident.NewService(repo)
	userID := uuid.New()

	_, err := svc.List(context.Background(), userID.String(), 0)
	if err != nil {
		t.Fatalf("expected valid list request, got %v", err)
	}

	if repo.listLimit != 100 {
		t.Fatalf("expected default limit 100, got %d", repo.listLimit)
	}
}

func TestIncidentListForMonitorRejectsInvalidIDs(t *testing.T) {
	repo := &fakeIncidentRepository{}
	svc := incident.NewService(repo)

	_, err := svc.ListForMonitor(context.Background(), "bad-user", uuid.NewString(), 100)
	if err == nil {
		t.Fatal("expected invalid user id error")
	}

	_, err = svc.ListForMonitor(context.Background(), uuid.NewString(), "bad-monitor", 100)
	if err == nil {
		t.Fatal("expected invalid monitor id error")
	}
}

func TestIncidentListForMonitorPassesParsedIDs(t *testing.T) {
	repo := &fakeIncidentRepository{}
	svc := incident.NewService(repo)
	userID := uuid.New()
	monitorID := uuid.New()

	_, err := svc.ListForMonitor(context.Background(), userID.String(), monitorID.String(), 25)
	if err != nil {
		t.Fatalf("expected valid monitor incident request, got %v", err)
	}

	if repo.monitorUserID != userID {
		t.Fatalf("expected user id %s, got %s", userID, repo.monitorUserID)
	}
	if repo.monitorID != monitorID {
		t.Fatalf("expected monitor id %s, got %s", monitorID, repo.monitorID)
	}
	if repo.monitorLimit != 25 {
		t.Fatalf("expected limit 25, got %d", repo.monitorLimit)
	}
}

type fakeIncidentRepository struct {
	listUserID uuid.UUID
	listLimit  int

	monitorUserID uuid.UUID
	monitorID     uuid.UUID
	monitorLimit  int
}

func (r *fakeIncidentRepository) List(ctx context.Context, userID uuid.UUID, limit int) ([]incident.Incident, error) {
	r.listUserID = userID
	r.listLimit = limit
	return nil, nil
}

func (r *fakeIncidentRepository) ListForMonitor(ctx context.Context, userID uuid.UUID, monitorID uuid.UUID, limit int) ([]incident.Incident, error) {
	r.monitorUserID = userID
	r.monitorID = monitorID
	r.monitorLimit = limit
	return nil, nil
}
