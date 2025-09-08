package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// NewRelicHook sends logs to New Relic
type NewRelicHook struct {
	LicenseKey string
	Region     string
	Client     *http.Client
	endpoint   string
}

// NewNewRelicHook creates a new New Relic hook
func NewNewRelicHook(licenseKey, region string) *NewRelicHook {
	var endpoint string
	if region == "eu" {
		endpoint = "https://log-api.eu.newrelic.com/log/v1"
	} else {
		endpoint = "https://log-api.newrelic.com/log/v1"
	}

	return &NewRelicHook{
		LicenseKey: licenseKey,
		Region:     region,
		Client:     &http.Client{Timeout: 10 * time.Second},
		endpoint:   endpoint,
	}
}

// Fire sends the log entry to New Relic
func (h *NewRelicHook) Fire(e *Entry) {
	// Convert to New Relic format
	nrLog := map[string]interface{}{
		"timestamp": e.Time.UnixMilli(),
		"level":     e.Level,
		"message":   e.Msg,
		"service":   "logx-app",
		"language":  "go",
	}

	// Add fields if they exist
	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			nrLog[k] = v
		}
	}

	// Add trace information if available
	if e.TraceID != "" {
		nrLog["trace.id"] = e.TraceID
	}
	if e.SpanID != "" {
		nrLog["span.id"] = e.SpanID
	}
	if e.Caller != "" {
		nrLog["caller"] = e.Caller
	}

	// New Relic expects an array of log objects
	payload := []interface{}{nrLog}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "newrelic hook encode error: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "newrelic hook request error: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", h.LicenseKey)

	go func() {
		resp, err := h.Client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "newrelic hook send error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			fmt.Fprintf(os.Stderr, "newrelic hook response error: %d\n", resp.StatusCode)
		}
	}()
}