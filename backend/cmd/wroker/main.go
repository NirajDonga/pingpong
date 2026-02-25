package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/NirajDonga/pingpong/backend/internal/ping"
	"github.com/NirajDonga/pingpong/backend/internal/shared"
)

func main() {
	nc, _ := nats.Connect("nats://127.0.0.1:4222")
	defer nc.Close()
	log.Println("✅ Worker connected to NATS")

	nc.Subscribe("ping.start", func(msg *nats.Msg) {
		var cmd shared.PingCommand
		json.Unmarshal(msg.Data, &cmd)
		log.Printf("📥 Testing %s for %ds", cmd.TargetURL, cmd.DurationSeconds)

		go runTest(nc, cmd) 
	})

	select {} 
}

func runTest(nc *nats.Conn, cmd shared.PingCommand) {
	topic := fmt.Sprintf("ping.result.%s", cmd.SessionID)

	for i := 0; i < cmd.DurationSeconds; i++ {
		metrics, err := ping.Measure(cmd.TargetURL)
		
		result := shared.PingResult{
			SessionID: cmd.SessionID,
			WorkerID:  "us-east-worker",
			Timestamp: time.Now(),
			Success:   err == nil,
			Metrics:   metrics,
		}

		bytes, _ := json.Marshal(result)
		nc.Publish(topic, bytes)

		time.Sleep(1 * time.Second)
	}
}