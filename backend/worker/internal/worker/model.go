package worker

import "time"

type CheckJob struct {
	MonitorID      string `json:"monitorId"`
	URL            string `json:"url"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
	ExpectedStatus int    `json:"expectedStatus"`
}

type CheckResult struct {
	MonitorID      string    `json:"monitorId"`
	CheckedAt      time.Time `json:"checkedAt"`
	Success        bool      `json:"success"`
	StatusCode     int       `json:"statusCode"`
	ResponseTimeMS int64     `json:"responseTimeMs"`
	DNSMS          int64     `json:"dnsMs"`
	TCPMS          int64     `json:"tcpMs"`
	TLSMS          int64     `json:"tlsMs"`
	TTFBMS         int64     `json:"ttfbMs"`
	Error          string    `json:"error"`
	WorkerName     string    `json:"workerName"`
}
