package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	"github.com/NirajDonga/pingpong/backend/internal/shared"
)

func main() {
	// 1. Start Embedded NATS
	opts := &server.Options{Host: "127.0.0.1", Port: 4222}
	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("Error creating NATS server: %v", err)
	}
	go ns.Start()
	ns.ReadyForConnections(5 * time.Second)

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		log.Fatalf("Error connecting NATS client: %v", err)
	}
	defer nc.Close()

	// 2. HTTP Server & SSE Endpoint
	http.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		sessionID := fmt.Sprintf("req_%d", time.Now().UnixMilli())

		// Broadcast Start Command
		cmd, _ := json.Marshal(shared.PingCommand{
			SessionID:       sessionID,
			TargetURL:       "https://example.com",
			DurationSeconds: 15,
		})
		nc.Publish("ping.start", cmd)

		// Stream Results Back
		msgChan := make(chan *nats.Msg, 100)
		sub, _ := nc.ChanSubscribe(fmt.Sprintf("ping.result.%s", sessionID), msgChan)
		defer sub.Unsubscribe()

		timeout := time.After(15 * time.Second)

		for {
			select {
			case msg := <-msgChan:
				fmt.Fprintf(w, "data: %s\n\n", string(msg.Data))
				flusher.Flush()
			case <-timeout:
				fmt.Fprintf(w, "data: {\"status\": \"completed\"}\n\n")
				flusher.Flush()
				return
			case <-r.Context().Done():
				return
			}
		}
	})

	log.Println("API listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
