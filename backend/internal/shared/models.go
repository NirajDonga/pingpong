package shared

type RequstTopic struct {
	SessionID       string `json:"sessionId"`
	TargetURL       string `json:"targetUrl"`
	DurationSeconds int    `json:"durationSeconds"`
}

type PingResult struct {
	SessionID      string `json:"sessionId"`
	WorkerName     string `json:"workerName"`
	TotalMs        int64  `json:"totalMs"`
	DNSMs          int64  `json:"dnsMs"`
	ConnMs         int64  `json:"connMs"`
	ServerMs       int64  `json:"serverMs"`
	ResponseMs     int64  `json:"responseMs"`
	TLSHandshakeMs int64  `json:"tlsHandshakeMs"`
	IsConnReused   bool   `json:"isConnReused"`
	Error          string `json:"error,omitempty"`
}
