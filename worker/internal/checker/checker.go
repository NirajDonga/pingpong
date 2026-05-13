package checker

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Request struct {
	URL            string
	TimeoutSeconds int
	ExpectedStatus int
}

type Result struct {
	CheckedAt      time.Time
	Success        bool
	StatusCode     int
	ResponseTimeMS int64
	DNSMS          int64
	TCPMS          int64
	TLSMS          int64
	TTFBMS         int64
	Error          string
}

type Checker struct {
	client *http.Client
}

func New() *Checker {
	return &Checker{
		client: &http.Client{},
	}
}

func (c *Checker) Check(ctx context.Context, request Request) Result {
	result := Result{
		CheckedAt: time.Now().UTC(),
	}

	timeout := time.Duration(request.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	timings := &requestTimings{}
	reqCtx = httptrace.WithClientTrace(reqCtx, timings.trace())

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, request.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	result.ResponseTimeMS = elapsedMS(start, time.Now())
	result.DNSMS = timings.dnsMS()
	result.TCPMS = timings.tcpMS()
	result.TLSMS = timings.tlsMS()
	result.TTFBMS = timings.ttfbMS()

	if err != nil {
		result.Error = checkError(err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode != request.ExpectedStatus {
		result.Error = fmt.Sprintf("expected status %d, got %d", request.ExpectedStatus, resp.StatusCode)
		return result
	}

	result.Success = true
	return result
}

type requestTimings struct {
	dnsStart      time.Time
	dnsDone       time.Time
	connectStart  time.Time
	connectDone   time.Time
	tlsStart      time.Time
	tlsDone       time.Time
	wroteRequest  time.Time
	firstResponse time.Time
}

func (t *requestTimings) trace() *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) {
			t.dnsStart = time.Now()
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			t.dnsDone = time.Now()
		},
		ConnectStart: func(_, _ string) {
			t.connectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			t.connectDone = time.Now()
		},
		TLSHandshakeStart: func() {
			t.tlsStart = time.Now()
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			t.tlsDone = time.Now()
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			t.wroteRequest = time.Now()
		},
		GotFirstResponseByte: func() {
			t.firstResponse = time.Now()
		},
	}
}

func (t *requestTimings) dnsMS() int64 {
	return elapsedMS(t.dnsStart, t.dnsDone)
}

func (t *requestTimings) tcpMS() int64 {
	return elapsedMS(t.connectStart, t.connectDone)
}

func (t *requestTimings) tlsMS() int64 {
	return elapsedMS(t.tlsStart, t.tlsDone)
}

func (t *requestTimings) ttfbMS() int64 {
	return elapsedMS(t.wroteRequest, t.firstResponse)
}

func elapsedMS(start time.Time, end time.Time) int64 {
	if start.IsZero() || end.IsZero() || end.Before(start) {
		return 0
	}
	return end.Sub(start).Milliseconds()
}

func checkError(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "request timed out"
	}
	return err.Error()
}
