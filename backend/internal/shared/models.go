package shared

import "time"

type PingCommand struct {
	SessionID       string `json:"sessionId"`
	TargetURL       string `json:"targetUrl"`
	DurationSeconds int    `json:"durationSeconds"`
}

type PingResult struct {
	SessionID  string    `json:"sessionId"`
	WorkerID   string    `json:"workerId"`
	Timestamp  time.Time `json:"timestamp"`
	Success    bool      `json:"success"`
	Metrics    Metrics   `json:"metrics"`
	Error      string    `json:"error,omitempty"`
}

type Metrics struct {
	DNSMs   int `json:"dnsMs"`
	TCPMs   int `json:"tcpMs"`
	TLSMs   int `json:"tlsMs"`
	TTFBMs  int `json:"ttfbMs"`
	TotalMs int `json:"totalMs"`
}