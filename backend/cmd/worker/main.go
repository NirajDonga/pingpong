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

	config.Load()
	natsURL := config.GetEnv("NATS_URL", "nats://127.0.0.1:4222")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Worker failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	workerName := config.GetEnv("WORKER_NAME", "local-worker")
	log.Printf("Worker [%s] connected to NATS at %s", workerName, natsURL)

	nc.Subscribe("ping.start", func(msg *nats.Msg) {
		var cmd shared.PingCommand
		json.Unmarshal(msg.Data, &cmd)
		log.Printf("Testing %s for %ds", cmd.TargetURL, cmd.DurationSeconds)

		go runTest(nc, cmd, workerName)
	})

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

		bytes, _ := json.Marshal(result)
		nc.Publish(topic, bytes)

		time.Sleep(1 * time.Second)
	}
}
