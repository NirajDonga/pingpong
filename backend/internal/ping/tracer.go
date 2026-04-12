package ping

import (
	"log"
	"time"

	"github.com/NirajDonga/pingpong/backend/internal/shared"
	"resty.dev/v3"
)

var client = resty.New().SetTimeout(5 * time.Second)

func Measure(target string, sessionID string, workerName string) shared.PingResult {
	res, err := client.R().EnableTrace().Get(target)

	if err != nil {
		log.Println("Error:", err)
		return shared.PingResult{
			SessionID:  sessionID,
			WorkerName: workerName,
			Error:      err.Error(),
		}
	}

	traceInfo := res.Request.TraceInfo()

	result := shared.PingResult{
		SessionID:      sessionID,
		WorkerName:     workerName,
		TotalMs:        traceInfo.TotalTime.Milliseconds(),
		DNSMs:          traceInfo.DNSLookup.Milliseconds(),
		ConnMs:         traceInfo.TCPConnTime.Milliseconds(),
		ServerMs:       traceInfo.ServerTime.Milliseconds(),
		ResponseMs:     traceInfo.ResponseTime.Milliseconds(),
		TLSHandshakeMs: traceInfo.TLSHandshake.Milliseconds(),
		IsConnReused:   traceInfo.IsConnReused,
	}

	return result
}
