package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// AtatusHook sends logs to Atatus
type AtatusHook struct {
	LicenseKey string
	AppName    string
	Client     *http.Client
	endpoint   string
}

// NewAtatusHook creates a new Atatus hook
func NewAtatusHook(licenseKey, appName string) *AtatusHook {
	endpoint := "https://api.atatus.com/api/v1/logs"

	return &AtatusHook{
		LicenseKey: licenseKey,
		AppName:    appName,
		Client:     &http.Client{Timeout: 10 * time.Second},
		endpoint:   endpoint,
	}
}

// Fire sends the log entry to Atatus
func (h *AtatusHook) Fire(e *Entry) {
	// Convert to Atatus format
	atatusLog := map[string]interface{}{
		"timestamp": e.Time.Format(time.RFC3339Nano),
		"level":     e.Level,
		"message":   e.Msg,
		"app":       h.AppName,
		"language":  "go",
	}

	// Add fields if they exist
	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			atatusLog[k] = v
		}
	}

	// Add trace information if available
	if e.TraceID != "" {
		atatusLog["trace_id"] = e.TraceID
	}
	if e.SpanID != "" {
		atatusLog["span_id"] = e.SpanID
	}
	if e.Caller != "" {
		atatusLog["caller"] = e.Caller
	}

	payload, err := json.Marshal(atatusLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "atatus hook encode error: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "atatus hook request error: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.LicenseKey)

	go func() {
		resp, err := h.Client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "atatus hook send error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			fmt.Fprintf(os.Stderr, "atatus hook response error: %d\n", resp.StatusCode)
		}
	}()
}