package main

import (
	"context"
	"log"

	"github.com/NirajDonga/pingpong/backend/worker/internal/checker"
	"github.com/NirajDonga/pingpong/backend/worker/internal/config"
	"github.com/NirajDonga/pingpong/backend/worker/internal/nats"
	"github.com/NirajDonga/pingpong/backend/worker/internal/worker"
)

func main() {
	cfg := config.Load()

	natsClient, err := nats.NewClient(cfg.NATSURL)
	if err != nil {
		log.Fatalf("worker failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	processor := worker.NewProcessor(checker.New(), natsClient, cfg.WorkerName)

	cons, err := natsClient.SubscribeCheckJobs(context.Background(), func(ctx context.Context, job worker.CheckJob) error {
		return processor.Process(ctx, job)
	})
	if err != nil {
		log.Fatalf("worker failed to subscribe to check jobs: %v", err)
	}
	defer cons.Stop()

	log.Println("worker service started")
	select {}
}
