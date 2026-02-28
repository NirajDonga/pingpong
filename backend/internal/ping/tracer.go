package ping

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/NirajDonga/pingpong/backend/internal/shared"
)

func Measure(targetURL string) (shared.Metrics, error) {
	var m shared.Metrics
	var start, dns, tcp, tlsTime time.Time

	req, err := http.NewRequest("HEAD", targetURL, nil)
	if err != nil {
		return m, err
	}

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { m.DNSMs = int(time.Since(dns).Milliseconds()) },
		GetConn:  func(_ string) { tcp = time.Now() },
		GotConn: func(info httptrace.GotConnInfo) {
			// If the connection was reused, we won't get true TCP time
			m.TCPMs = int(time.Since(tcp).Milliseconds())
		},
		TLSHandshakeStart:    func() { tlsTime = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { m.TLSMs = int(time.Since(tlsTime).Milliseconds()) },
		GotFirstResponseByte: func() { m.TTFBMs = int(time.Since(start).Milliseconds()) },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	// Create a transport that DOES NOT reuse connections
	transport := &http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	start = time.Now()
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}

	m.TotalMs = int(time.Since(start).Milliseconds())
	return m, err
}
