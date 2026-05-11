package database

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func PingClickHouse(ctx context.Context, baseURL string) error {
	pingURL := strings.TrimRight(baseURL, "/") + "/ping"

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(pingCtx, http.MethodGet, pingURL, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	return nil
}
