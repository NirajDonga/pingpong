package result

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TinybirdRepository stores and retrieves check results via Tinybird's HTTP APIs.
type TinybirdRepository struct {
	host        string
	appendToken string
	readToken   string
	client      *http.Client
}

func NewTinybirdRepository(host, appendToken, readToken string) *TinybirdRepository {
	return &TinybirdRepository{
		host:        strings.TrimRight(host, "/"),
		appendToken: appendToken,
		readToken:   readToken,
		client:      &http.Client{},
	}
}

// Insert sends a single check result to Tinybird via the Events API (NDJSON).
func (r *TinybirdRepository) Insert(ctx context.Context, result CheckResult) error {
	row := tinybirdEventRow{
		MonitorID:      result.MonitorID,
		CheckedAt:      result.CheckedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		Success:        boolToUint8(result.Success),
		StatusCode:     result.StatusCode,
		ResponseTimeMS: result.ResponseTimeMS,
		DNSMS:          result.DNSMS,
		TCPMS:          result.TCPMS,
		TLSMS:          result.TLSMS,
		TTFBMS:         result.TTFBMS,
		Error:          result.Error,
		WorkerName:     result.WorkerName,
	}

	data, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("tinybird marshal: %w", err)
	}
	// NDJSON: single line terminated by newline
	data = append(data, '\n')

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := r.host + "/v0/events?name=check_results"
	req, err := http.NewRequestWithContext(opCtx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("tinybird insert request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.appendToken)
	req.Header.Set("Content-Type", "application/x-ndjson")

	res, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("tinybird insert: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("tinybird insert failed (%s): %s", res.Status, strings.TrimSpace(string(body)))
	}

	return nil
}

// History queries Tinybird's SQL API and returns check results for a monitor.
func (r *TinybirdRepository) History(ctx context.Context, monitorID uuid.UUID, limit int) ([]CheckResult, error) {
	query := fmt.Sprintf(
		`SELECT monitor_id, checked_at, success, status_code, response_time_ms, dns_ms, tcp_ms, tls_ms, ttfb_ms, error, worker_name FROM check_results WHERE monitor_id = '%s' ORDER BY checked_at DESC LIMIT %d FORMAT JSON`,
		monitorID.String(), limit,
	)

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	reqURL := r.host + "/v0/sql?q=" + url.QueryEscape(query)
	req, err := http.NewRequestWithContext(opCtx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("tinybird query request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.readToken)

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tinybird query: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("tinybird read body: %w", err)
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("tinybird query failed (%s): %s", res.Status, strings.TrimSpace(string(body)))
	}

	var response tinybirdQueryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("tinybird unmarshal: %w", err)
	}

	results := make([]CheckResult, 0, len(response.Data))
	for _, row := range response.Data {
		checkedAt, err := time.Parse("2006-01-02 15:04:05.000", row.CheckedAt)
		if err != nil {
			// Try ISO format as fallback
			checkedAt, err = time.Parse("2006-01-02T15:04:05.000Z", row.CheckedAt)
			if err != nil {
				return nil, fmt.Errorf("tinybird parse checked_at %q: %w", row.CheckedAt, err)
			}
		}

		results = append(results, CheckResult{
			MonitorID:      row.MonitorID,
			CheckedAt:      checkedAt.UTC(),
			Success:        row.Success != 0,
			StatusCode:     row.StatusCode,
			ResponseTimeMS: row.ResponseTimeMS,
			DNSMS:          row.DNSMS,
			TCPMS:          row.TCPMS,
			TLSMS:          row.TLSMS,
			TTFBMS:         row.TTFBMS,
			Error:          row.Error,
			WorkerName:     row.WorkerName,
		})
	}

	return results, nil
}


func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// tinybirdEventRow matches the Tinybird check_results datasource schema for writes.
type tinybirdEventRow struct {
	MonitorID      string `json:"monitor_id"`
	CheckedAt      string `json:"checked_at"`
	Success        uint8  `json:"success"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMS int64  `json:"response_time_ms"`
	DNSMS          int64  `json:"dns_ms"`
	TCPMS          int64  `json:"tcp_ms"`
	TLSMS          int64  `json:"tls_ms"`
	TTFBMS         int64  `json:"ttfb_ms"`
	Error          string `json:"error"`
	WorkerName     string `json:"worker_name"`
}

// tinybirdQueryRow matches a row from the Tinybird SQL query response for reads.
type tinybirdQueryRow struct {
	MonitorID      string `json:"monitor_id"`
	CheckedAt      string `json:"checked_at"`
	Success        uint8  `json:"success"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMS int64  `json:"response_time_ms"`
	DNSMS          int64  `json:"dns_ms"`
	TCPMS          int64  `json:"tcp_ms"`
	TLSMS          int64  `json:"tls_ms"`
	TTFBMS         int64  `json:"ttfb_ms"`
	Error          string `json:"error"`
	WorkerName     string `json:"worker_name"`
}

// tinybirdQueryResponse wraps the JSON envelope returned by Tinybird's SQL API.
type tinybirdQueryResponse struct {
	Data []tinybirdQueryRow `json:"data"`
}
