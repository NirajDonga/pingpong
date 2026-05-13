package tests

import (
	"context"
	"testing"
	"time"

	"github.com/NirajDonga/pingpong/api/internal/result"
	"github.com/google/uuid"
)

func TestProcessStoresValidResult(t *testing.T) {
	repo := &fakeRepository{}
	updater := &fakeStatusUpdater{}
	svc := result.NewService(repo, updater)

	monitorID := uuid.New()

	err := svc.Process(context.Background(), result.CheckResult{
		MonitorID:      monitorID.String(),
		CheckedAt:      time.Now(),
		Success:        true,
		StatusCode:     200,
		ResponseTimeMS: 42,
		WorkerName:     "worker-test",
	})
	if err != nil {
		t.Fatalf("expected valid result, got %v", err)
	}

	if len(repo.inserted) != 1 {
		t.Fatalf("expected 1 inserted result, got %d", len(repo.inserted))
	}
	if len(updater.applied) != 1 {
		t.Fatalf("expected 1 status update, got %d", len(updater.applied))
	}
	if updater.applied[0].monitorID != monitorID {
		t.Fatalf("expected monitor id %s, got %s", monitorID, updater.applied[0].monitorID)
	}
	if !updater.applied[0].success {
		t.Fatal("expected success to be applied")
	}
}

func TestProcessRejectsInvalidResult(t *testing.T) {
	repo := &fakeRepository{}
	updater := &fakeStatusUpdater{}
	svc := result.NewService(repo, updater)

	err := svc.Process(context.Background(), result.CheckResult{
		MonitorID:  "not-a-uuid",
		CheckedAt:  time.Now(),
		StatusCode: 200,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(repo.inserted) != 0 {
		t.Fatalf("expected no inserted results, got %d", len(repo.inserted))
	}
	if len(updater.applied) != 0 {
		t.Fatalf("expected no status updates, got %d", len(updater.applied))
	}
}

type fakeRepository struct {
	inserted []result.CheckResult
}

func (r *fakeRepository) Insert(ctx context.Context, result result.CheckResult) error {
	r.inserted = append(r.inserted, result)
	return nil
}

func (r *fakeRepository) History(ctx context.Context, monitorID uuid.UUID, limit int) ([]result.CheckResult, error) {
	return nil, nil
}

type appliedResult struct {
	monitorID uuid.UUID
	success   bool
}

type fakeStatusUpdater struct {
	applied []appliedResult
}

func (u *fakeStatusUpdater) ApplyCheckResult(ctx context.Context, monitorID uuid.UUID, success bool) error {
	u.applied = append(u.applied, appliedResult{
		monitorID: monitorID,
		success:   success,
	})
	return nil
}
