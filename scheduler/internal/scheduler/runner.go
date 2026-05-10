package scheduler

import (
	"context"
	"log"
	"time"
)

type Publisher interface {
	PublishCheckJob(job CheckJob) error
}

type Runner struct {
	store     *Store
	publisher Publisher
	batchSize int
	tickEvery time.Duration
}

func NewRunner(store *Store, publisher Publisher, batchSize int, tickEvery time.Duration) *Runner {
	return &Runner{
		store:     store,
		publisher: publisher,
		batchSize: batchSize,
		tickEvery: tickEvery,
	}
}

func (r *Runner) Run(ctx context.Context) {
	log.Println("scheduler loop started")

	ticker := time.NewTicker(r.tickEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("scheduler loop stopped")
			return
		case <-ticker.C:
			r.dispatchDue(ctx)
		}
	}
}

func (r *Runner) dispatchDue(ctx context.Context) {
	monitors, err := r.store.DueMonitors(ctx, r.batchSize)
	if err != nil {
		log.Printf("failed to read due monitors: %v", err)
		return
	}

	for _, monitor := range monitors {
		job := CheckJob{
			MonitorID:      monitor.ID.String(),
			URL:            monitor.URL,
			TimeoutSeconds: monitor.TimeoutSeconds,
			ExpectedStatus: monitor.ExpectedStatus,
		}

		if err := r.publisher.PublishCheckJob(job); err != nil {
			log.Printf("failed to publish check job for monitor %s: %v", monitor.ID, err)
			continue
		}

		if err := r.store.UpdateNextCheck(ctx, monitor.ID, monitor.IntervalSeconds); err != nil {
			log.Printf("failed to update next_check_at for monitor %s: %v", monitor.ID, err)
			continue
		}
	}
}
