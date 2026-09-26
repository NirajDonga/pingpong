package processor

import "time"

type CheckResult struct {
	MonitorID      string    `json:"monitor_id"`
	CheckedAt      time.Time `json:"checked_at"`
	Success        bool      `json:"success"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMS int       `json:"response_time_ms"`
	DNSMS          int       `json:"dns_ms"`
	TCPMS          int       `json:"tcp_ms"`
	TLSMS          int       `json:"tls_ms"`
	TTFBMS         int       `json:"ttfb_ms"`
	Error          string    `json:"error,omitempty"`
	WorkerName     string    `json:"worker_name"`
}
