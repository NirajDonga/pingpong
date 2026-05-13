package result

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ClickHouseRepository struct {
	baseURL string
	client  *http.Client
}

func NewClickHouseRepository(baseURL string) *ClickHouseRepository {
	return &ClickHouseRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{},
	}
}

func (r *ClickHouseRepository) Insert(ctx context.Context, result CheckResult) error {
	row := clickHouseResultRow{
		MonitorID:      result.MonitorID,
		CheckedAt:      result.CheckedAt.UTC().Format("2006-01-02 15:04:05.000"),
		Success:        result.Success,
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
		return err
	}

	query := "INSERT INTO check_results FORMAT JSONEachRow"
	return r.exec(ctx, query, append(data, '\n'))
}

func (r *ClickHouseRepository) History(ctx context.Context, monitorID uuid.UUID, limit int) ([]CheckResult, error) {
	query := fmt.Sprintf(`
		SELECT
			toString(monitor_id) AS monitorId,
			toUnixTimestamp64Milli(checked_at) AS checkedAtMs,
			success,
			status_code AS statusCode,
			response_time_ms AS responseTimeMs,
			dns_ms AS dnsMs,
			tcp_ms AS tcpMs,
			tls_ms AS tlsMs,
			ttfb_ms AS ttfbMs,
			error,
			worker_name AS workerName
		FROM check_results
		WHERE monitor_id = toUUID('%s')
		ORDER BY checked_at DESC
		LIMIT %d
		FORMAT JSONEachRow
	`, monitorID.String(), limit)

	body, err := r.query(ctx, query)
	if err != nil {
		return nil, err
	}

	results := []CheckResult{}
	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var row historyRow
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			return nil, err
		}

		results = append(results, CheckResult{
			MonitorID:      row.MonitorID,
			CheckedAt:      time.UnixMilli(row.CheckedAtMS).UTC(),
			Success:        row.Success,
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
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *ClickHouseRepository) exec(ctx context.Context, query string, body []byte) error {
	_, err := r.do(ctx, query, body)
	return err
}

func (r *ClickHouseRepository) query(ctx context.Context, query string) ([]byte, error) {
	return r.do(ctx, query, nil)
}

func (r *ClickHouseRepository) do(ctx context.Context, query string, body []byte) ([]byte, error) {
	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(opCtx, http.MethodPost, r.baseURL+"/?query="+url.QueryEscape(query), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out bytes.Buffer
	if _, err := out.ReadFrom(res.Body); err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("clickhouse query failed: %s: %s", res.Status, strings.TrimSpace(out.String()))
	}

	return out.Bytes(), nil
}

type clickHouseResultRow struct {
	MonitorID      string `json:"monitor_id"`
	CheckedAt      string `json:"checked_at"`
	Success        bool   `json:"success"`
	StatusCode     int    `json:"status_code"`
	ResponseTimeMS int64  `json:"response_time_ms"`
	DNSMS          int64  `json:"dns_ms"`
	TCPMS          int64  `json:"tcp_ms"`
	TLSMS          int64  `json:"tls_ms"`
	TTFBMS         int64  `json:"ttfb_ms"`
	Error          string `json:"error"`
	WorkerName     string `json:"worker_name"`
}

type historyRow struct {
	MonitorID      string `json:"monitorId"`
	CheckedAtMS    int64  `json:"checkedAtMs"`
	Success        bool   `json:"success"`
	StatusCode     int    `json:"statusCode"`
	ResponseTimeMS int64  `json:"responseTimeMs"`
	DNSMS          int64  `json:"dnsMs"`
	TCPMS          int64  `json:"tcpMs"`
	TLSMS          int64  `json:"tlsMs"`
	TTFBMS         int64  `json:"ttfbMs"`
	Error          string `json:"error"`
	WorkerName     string `json:"workerName"`
}
