package main

import (
	"context"
	"log"

	"github.com/NirajDonga/pingpong/scheduler/internal/config"
	"github.com/NirajDonga/pingpong/scheduler/internal/database"
	"github.com/NirajDonga/pingpong/scheduler/internal/nats"
	"github.com/NirajDonga/pingpong/scheduler/internal/scheduler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load: %v", err)
	}

	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("scheduler failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	publisher, err := nats.NewPublisher(cfg.NATSURL)
	if err != nil {
		log.Fatalf("scheduler failed to connect to NATS: %v", err)
	}
	defer publisher.Close()

	store := scheduler.NewStore(db)
	runner := scheduler.NewRunner(store, publisher, cfg.BatchSize, cfg.TickEvery)

	log.Println("scheduler service started")
	runner.Run(context.Background())
}
