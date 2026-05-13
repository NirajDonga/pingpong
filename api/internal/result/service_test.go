package result

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProcessStoresValidResult(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewService(repo)

	err := svc.Process(context.Background(), CheckResult{
		MonitorID:      uuid.NewString(),
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
}

func TestProcessRejectsInvalidResult(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewService(repo)

	err := svc.Process(context.Background(), CheckResult{
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
}

type fakeRepository struct {
	inserted []CheckResult
}

func (r *fakeRepository) Insert(ctx context.Context, result CheckResult) error {
	r.inserted = append(r.inserted, result)
	return nil
}

func (r *fakeRepository) History(ctx context.Context, monitorID uuid.UUID, limit int) ([]CheckResult, error) {
	return nil, nil
}
