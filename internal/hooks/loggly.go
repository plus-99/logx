package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// LogglyHook sends logs to Loggly
type LogglyHook struct {
	Token    string
	Tag      string
	Client   *http.Client
	endpoint string
}

// NewLogglyHook creates a new Loggly hook
func NewLogglyHook(token, tag string) *LogglyHook {
	endpoint := fmt.Sprintf("https://logs-01.loggly.com/inputs/%s/tag/%s/", token, tag)

	return &LogglyHook{
		Token:    token,
		Tag:      tag,
		Client:   &http.Client{Timeout: 10 * time.Second},
		endpoint: endpoint,
	}
}

// Fire sends the log entry to Loggly
func (h *LogglyHook) Fire(e *Entry) {
	// Convert to Loggly format
	logglyLog := map[string]interface{}{
		"timestamp": e.Time.Format(time.RFC3339),
		"level":     e.Level,
		"message":   e.Msg,
		"service":   "logx-app",
		"language":  "go",
	}

	// Add fields if they exist
	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			logglyLog[k] = v
		}
	}

	// Add trace information if available
	if e.TraceID != "" {
		logglyLog["trace_id"] = e.TraceID
	}
	if e.SpanID != "" {
		logglyLog["span_id"] = e.SpanID
	}
	if e.Caller != "" {
		logglyLog["caller"] = e.Caller
	}

	payload, err := json.Marshal(logglyLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loggly hook encode error: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "loggly hook request error: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	go func() {
		resp, err := h.Client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "loggly hook send error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			fmt.Fprintf(os.Stderr, "loggly hook response error: %d\n", resp.StatusCode)
		}
	}()
}