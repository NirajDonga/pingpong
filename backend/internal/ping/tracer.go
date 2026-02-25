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

	req, _ := http.NewRequest("HEAD", targetURL, nil)

	// Attach hooks to record exact network timings
	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { m.DNSMs = int(time.Since(dns).Milliseconds()) },
		GetConn:              func(_ string) { tcp = time.Now() },
		GotConn:              func(_ httptrace.GotConnInfo) { m.TCPMs = int(time.Since(tcp).Milliseconds()) },
		TLSHandshakeStart:    func() { tlsTime = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { m.TLSMs = int(time.Since(tlsTime).Milliseconds()) },
		GotFirstResponseByte: func() { m.TTFBMs = int(time.Since(start).Milliseconds()) },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	client := &http.Client{Timeout: 5 * time.Second}

	start = time.Now()
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	
	m.TotalMs = int(time.Since(start).Milliseconds())
	return m, err
}