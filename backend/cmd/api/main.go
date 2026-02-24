package main

import (
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {

	opts := &server.Options{
		Host: "127.0.0.1",
		Port: 4222,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("Error creating NATS server: %v", err)
	}

	go ns.Start()

	if !ns.ReadyForConnections(5 * time.Second) {
		log.Fatal("NATS server failed to start within 5 seconds")
	}

	log.Printf("Embedded NATS server running on %s:%d", opts.Host, opts.Port)

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		log.Fatalf("Error connecting NATS client: %v", err)
	}
	defer nc.Close()

	log.Println("Central API connected to internal NATS broker")

	select {}
}
