package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/NirajDonga/pingpong/backend/internal/config"
	"github.com/NirajDonga/pingpong/backend/internal/ping"
	"github.com/NirajDonga/pingpong/backend/internal/shared"
	"github.com/nats-io/nats.go"
)

func main() {

	cfg := config.LoadWorkerConfig()
	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		log.Fatalf("Worker failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	workerName := cfg.WorkerName
	log.Printf("Worker [%s] connected to NATS at %s", workerName, cfg.NATSURL)

	_, err = nc.Subscribe("ping.start", func(msg *nats.Msg) {
		var cmd shared.PingCommand
		if err := json.Unmarshal(msg.Data, &cmd); err != nil {
			log.Printf("invalid ping.start payload: %v", err)
			return
		}
		log.Printf("Testing %s for %ds", cmd.TargetURL, cmd.DurationSeconds)

		go runTest(nc, cmd, workerName)
	})
	if err != nil {
		log.Fatalf("failed to subscribe to ping.start: %v", err)
	}

	select {}
}

func runTest(nc *nats.Conn, cmd shared.PingCommand, workerName string) {
	topic := fmt.Sprintf("ping.result.%s", cmd.SessionID)

	for i := 0; i < cmd.DurationSeconds; i++ {
		metrics, err := ping.Measure(cmd.TargetURL)

		result := shared.PingResult{
			SessionID: cmd.SessionID,
			WorkerID:  workerName,
			Timestamp: time.Now(),
			Success:   err == nil,
			Metrics:   metrics,
		}
		if err != nil {
			result.Error = err.Error()
		}

		bytes, _ := json.Marshal(result)
		nc.Publish(topic, bytes)

		time.Sleep(1 * time.Second)
	}
}
