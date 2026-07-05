package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	mu      sync.RWMutex
	clients map[string]map[*websocket.Conn]bool
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (m *Manager) AddClient(monitorID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.clients[monitorID] == nil {
		m.clients[monitorID] = make(map[*websocket.Conn]bool)
	}
	m.clients[monitorID][conn] = true
	log.Printf("ws: client connected for monitor %s (%d clients)", monitorID, len(m.clients[monitorID]))
}

func (m *Manager) RemoveClient(monitorID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	conns, ok := m.clients[monitorID]
	if !ok {
		return
	}

	delete(conns, conn)
	log.Printf("ws: client disconnected for monitor %s (%d clients remaining)", monitorID, len(conns))

	if len(conns) == 0 {
		delete(m.clients, monitorID)
	}
}

func (m *Manager) Broadcast(monitorID string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ws: failed to marshal payload: %v", err)
		return
	}

	m.mu.RLock()
	conns, ok := m.clients[monitorID]
	if !ok || len(conns) == 0 {
		m.mu.RUnlock()
		return
	}

	targets := make([]*websocket.Conn, 0, len(conns))
	for conn := range conns {
		targets = append(targets, conn)
	}
	m.mu.RUnlock()

	var failed []*websocket.Conn
	for _, conn := range targets {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("ws: write error, removing client: %v", err)
			conn.Close()
			failed = append(failed, conn)
		}
	}

	if len(failed) > 0 {
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, conn := range failed {
			delete(m.clients[monitorID], conn)
		}
		if len(m.clients[monitorID]) == 0 {
			delete(m.clients, monitorID)
		}
	}
}
