package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/NirajDonga/pingpong/backend/internal/config"
	"github.com/NirajDonga/pingpong/backend/internal/ping"
	"github.com/NirajDonga/pingpong/backend/internal/shared"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg := config.LoadWorkerConfig()

	nc, err := nats.Connect(cfg.NATSUrl)
	if err != nil {
		log.Fatalf("Worker failed to connect to NATS: %v", err)
	}

	defer nc.Close()
	workerName := cfg.WorkerName

	_, err = nc.Subscribe("start", func(msg *nats.Msg) {
		var content shared.RequstTopic
		if err := json.Unmarshal(msg.Data, &content); err != nil {
			log.Printf("invalid ping.start payload: %v", err)
			return
		}
		go runPinger(nc, content, workerName)
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	select {}
}

func runPinger(nc *nats.Conn, content shared.RequstTopic, workerName string) {
	// 1. Ticker generates a signal every 1 second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 2. Timeout generates exactly ONE signal after DurationSeconds (e.g. 60s)
	duration := time.Duration(content.DurationSeconds) * time.Second
	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			return

		case <-ticker.C:
			trace := ping.Measure(content.TargetURL, content.SessionID, workerName)

			resultBytes, err := json.Marshal(trace)
			if err != nil {
				log.Printf("Worker failed to marshal live result: %v", err)
				continue
			}

			nc.Publish(content.SessionID, resultBytes)
		}
	}
}
