package main

import (
	"context"
	"log"

	"github.com/NirajDonga/pingpong/worker/internal/checker"
	"github.com/NirajDonga/pingpong/worker/internal/config"
	"github.com/NirajDonga/pingpong/worker/internal/nats"
	"github.com/NirajDonga/pingpong/worker/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	natsClient, err := nats.NewClient(cfg.NATSURL)
	if err != nil {
		log.Fatalf("worker failed to connect to NATS: %v", err)
	}
	defer natsClient.Close()

	processor := worker.NewProcessor(checker.New(), natsClient, cfg.WorkerName)

	_, err = natsClient.SubscribeCheckJobs(func(job worker.CheckJob) {
		go processor.Process(context.Background(), job)
	})
	if err != nil {
		log.Fatalf("worker failed to subscribe to check jobs: %v", err)
	}

	log.Println("worker service started")
	select {}
}
