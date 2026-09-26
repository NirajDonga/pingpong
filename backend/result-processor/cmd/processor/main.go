package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NirajDonga/pingpong/backend/result-processor/internal/config"
	"github.com/NirajDonga/pingpong/backend/result-processor/internal/database"
	"github.com/NirajDonga/pingpong/backend/result-processor/internal/nats"
	"github.com/NirajDonga/pingpong/backend/result-processor/internal/processor"
)

func main() {
	cfg := config.Load()

	log.Println("connecting to postgres")
	db, err := database.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("postgres connect: %v", err)
	}
	defer db.Close()
	log.Println("connected to postgres")

	log.Println("connecting to nats")
	natsClient, err := nats.NewClient(cfg.NATSURL)
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer natsClient.Close()
	log.Println("connected to nats")

	tbRepo := processor.NewTinybirdRepository(cfg.TinybirdHost, cfg.TinybirdAppendToken)
	pgRepo := processor.NewPostgresRepository(db)
	svc := processor.NewService(tbRepo, pgRepo)

	sub, err := natsClient.SubscribeCheckResults(func(checkResult processor.CheckResult) {
		if err := svc.Process(context.Background(), checkResult); err != nil {
			log.Printf("failed to process check result for monitor %s: %v", checkResult.MonitorID, err)
			return
		}
		log.Printf("processed check result for monitor %s success=%t status=%d", checkResult.MonitorID, checkResult.Success, checkResult.StatusCode)
	})
	if err != nil {
		log.Fatalf("check result subscription: %v", err)
	}
	defer sub.Unsubscribe()

	log.Println("result-processor is running. Press Ctrl+C to stop")
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down result-processor...")
}
