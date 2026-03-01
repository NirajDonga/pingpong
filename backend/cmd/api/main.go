package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"

	"github.com/NirajDonga/pingpong/backend/internal/config"
	"github.com/NirajDonga/pingpong/backend/internal/shared"
)

func main() {
	cfg := config.LoadAPIConfig()

	natsp, err := strconv.Atoi(cfg.NATSPort)
	if err != nil {
		log.Fatalf("Invalid NATS_PORT: %v", err)
	}

	opts := &server.Options{Host: cfg.NATSHost, Port: natsp}
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

		target := r.URL.Query().Get("target")
		if target == "" {
			http.Error(w, "Bad Request: 'target' parameter is required", http.StatusBadRequest)
			return
		}

		cmd, err := json.Marshal(shared.PingCommand{
			SessionID:       sessionID,
			TargetURL:       target,
			DurationSeconds: 15,
		})
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if err := nc.Publish("ping.start", cmd); err != nil {
			http.Error(w, "Failed to publish start command", http.StatusInternalServerError)
			return
		}

		msgChan := make(chan *nats.Msg, 100)
		sub, err := nc.ChanSubscribe(fmt.Sprintf("ping.result.%s", sessionID), msgChan)
		if err != nil {
			http.Error(w, "Failed to subscribe to results", http.StatusInternalServerError)
			return
		}
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

	apiPort := cfg.Port
	log.Printf("API listening on :%s", apiPort)
	if err := http.ListenAndServe(":"+apiPort, nil); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
