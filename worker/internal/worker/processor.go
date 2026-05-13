package worker

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/NirajDonga/pingpong/worker/internal/checker"
)

type ResultPublisher interface {
	PublishCheckResult(result CheckResult) error
}

type Processor struct {
	checker    *checker.Checker
	publisher  ResultPublisher
	workerName string
}

func NewProcessor(checker *checker.Checker, publisher ResultPublisher, workerName string) *Processor {
	return &Processor{
		checker:    checker,
		publisher:  publisher,
		workerName: workerName,
	}
}

func (p *Processor) Process(ctx context.Context, job CheckJob) {
	if err := validateJob(job); err != nil {
		log.Printf("invalid check job: %v", err)
		return
	}

	checkResult := p.checker.Check(ctx, checker.Request{
		URL:            job.URL,
		TimeoutSeconds: job.TimeoutSeconds,
		ExpectedStatus: job.ExpectedStatus,
	})
	result := CheckResult{
		MonitorID:      job.MonitorID,
		CheckedAt:      checkResult.CheckedAt,
		Success:        checkResult.Success,
		StatusCode:     checkResult.StatusCode,
		ResponseTimeMS: checkResult.ResponseTimeMS,
		DNSMS:          checkResult.DNSMS,
		TCPMS:          checkResult.TCPMS,
		TLSMS:          checkResult.TLSMS,
		TTFBMS:         checkResult.TTFBMS,
		Error:          checkResult.Error,
		WorkerName:     p.workerName,
	}

	if err := p.publisher.PublishCheckResult(result); err != nil {
		log.Printf("failed to publish check result for monitor %s: %v", job.MonitorID, err)
		return
	}

	log.Printf("published check result for monitor %s success=%t status=%d", job.MonitorID, result.Success, result.StatusCode)
}

func validateJob(job CheckJob) error {
	if strings.TrimSpace(job.MonitorID) == "" {
		return errors.New("monitorId is required")
	}
	if strings.TrimSpace(job.URL) == "" {
		return errors.New("url is required")
	}
	if job.ExpectedStatus <= 0 {
		return errors.New("expectedStatus must be positive")
	}
	return nil
}
