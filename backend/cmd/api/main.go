package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NirajDonga/pingpong/backend/internal/config"
	"github.com/NirajDonga/pingpong/backend/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {

	cfg := config.LoadApiConfig()

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
		log.Fatal("NATS server did not become ready within 5s")
	}

	nc, err := nats.Connect(ns.ClientURL())
	if err != nil {
		log.Fatalf("Error connecting NATS client: %v", err)
	}
	defer nc.Close()

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	router := gin.Default()
	router.GET("/api/ping", func(c *gin.Context) {
		target := c.Query("target")
		if target == "" {
			c.String(http.StatusBadRequest, "Bad Request: 'target' parameter is required")
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("websocket upgrade failed: %v", err)
			return
		}
		defer conn.Close()

		sessionID := fmt.Sprintf("req_%d", time.Now().UnixMilli())

		runDuration := 60
		content := shared.RequstTopic{
			SessionID:       sessionID,
			TargetURL:       target,
			DurationSeconds: runDuration,
		}
		contentBytes, err := json.Marshal(content)
		if err != nil {
			log.Printf("Failed to format NATS message: %v", err)
			return
		}

		sub, err := nc.Subscribe(sessionID, func(msg *nats.Msg) {
			if writeErr := conn.WriteMessage(websocket.TextMessage, msg.Data); writeErr != nil {
				log.Printf("Failed to write to websocket for session %s: %v", sessionID, writeErr)
			}
		})
		if err != nil {
			log.Printf("failed to subscribe to NATS results: %v", err)
			return
		}
		defer sub.Unsubscribe()

		nc.Publish("start", contentBytes)

		done := make(chan struct{})
		go func() {
			select {
			case <-time.After(time.Duration(runDuration)*time.Second + 2*time.Second):
				message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "session completed")
				if err := conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second)); err != nil {
					log.Printf("Failed to close websocket for session %s: %v", sessionID, err)
				}
				conn.Close()
			case <-done:
			}
		}()

		for {
			_, _, readErr := conn.ReadMessage()
			if readErr != nil {
				log.Printf("Client disconnected from session %s", sessionID)
				break
			}
		}
		close(done)
	})

	router.Run(":" + cfg.Port)
}
