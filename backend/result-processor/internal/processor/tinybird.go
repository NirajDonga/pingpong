package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TinybirdRepository struct {
	host        string
	appendToken string
	client      *http.Client
}

func NewTinybirdRepository(host, appendToken string) *TinybirdRepository {
	return &TinybirdRepository{
		host:        strings.TrimRight(host, "/"),
		appendToken: appendToken,
		client:      &http.Client{},
	}
}

func (r *TinybirdRepository) Insert(ctx context.Context, result CheckResult) error {
	row := map[string]interface{}{
		"monitor_id":       result.MonitorID,
		"checked_at":       result.CheckedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		"success":          boolToUint8(result.Success),
		"status_code":      result.StatusCode,
		"response_time_ms": result.ResponseTimeMS,
		"dns_ms":           result.DNSMS,
		"tcp_ms":           result.TCPMS,
		"tls_ms":           result.TLSMS,
		"ttfb_ms":          result.TTFBMS,
		"error":            result.Error,
		"worker_name":      result.WorkerName,
	}

	data, err := json.Marshal(row)
	if err != nil {
		return fmt.Errorf("tinybird marshal: %w", err)
	}
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

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
