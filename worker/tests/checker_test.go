package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NirajDonga/pingpong/worker/internal/checker"
)

func TestCheckerSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	result := checker.New().Check(context.Background(), checker.Request{
		URL:            server.URL,
		TimeoutSeconds: 1,
		ExpectedStatus: http.StatusNoContent,
	})

	if !result.Success {
		t.Fatalf("expected success, got error %q", result.Error)
	}
	if result.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, result.StatusCode)
	}
	if result.CheckedAt.IsZero() {
		t.Fatal("expected checked_at to be set")
	}
}

func TestCheckerExpectedStatusFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	result := checker.New().Check(context.Background(), checker.Request{
		URL:            server.URL,
		TimeoutSeconds: 1,
		ExpectedStatus: http.StatusOK,
	})

	if result.Success {
		t.Fatal("expected failed result")
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, result.StatusCode)
	}
	if !strings.Contains(result.Error, "expected status 200") {
		t.Fatalf("expected status error, got %q", result.Error)
	}
}

func TestCheckerTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result := checker.New().Check(ctx, checker.Request{
		URL:            server.URL,
		TimeoutSeconds: 1,
		ExpectedStatus: http.StatusOK,
	})

	if result.Success {
		t.Fatal("expected timeout failure")
	}
	if result.Error != "request timed out" {
		t.Fatalf("expected timeout error, got %q", result.Error)
	}
}
