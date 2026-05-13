package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NirajDonga/pingpong/worker/internal/checker"
	workerpkg "github.com/NirajDonga/pingpong/worker/internal/worker"
)

func TestProcessorPublishesResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	publisher := &fakePublisher{}
	processor := workerpkg.NewProcessor(checker.New(), publisher, "worker-test")

	processor.Process(context.Background(), workerpkg.CheckJob{
		MonitorID:      "monitor-1",
		URL:            server.URL,
		TimeoutSeconds: 1,
		ExpectedStatus: http.StatusOK,
	})

	if len(publisher.results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(publisher.results))
	}
	result := publisher.results[0]
	if !result.Success {
		t.Fatalf("expected success, got error %q", result.Error)
	}
	if result.WorkerName != "worker-test" {
		t.Fatalf("expected worker name to be set, got %q", result.WorkerName)
	}
}

func TestProcessorSkipsInvalidJob(t *testing.T) {
	publisher := &fakePublisher{}
	processor := workerpkg.NewProcessor(checker.New(), publisher, "worker-test")

	processor.Process(context.Background(), workerpkg.CheckJob{})

	if len(publisher.results) != 0 {
		t.Fatalf("expected no results, got %d", len(publisher.results))
	}
}

type fakePublisher struct {
	results []workerpkg.CheckResult
}

func (p *fakePublisher) PublishCheckResult(result workerpkg.CheckResult) error {
	p.results = append(p.results, result)
	return nil
}
