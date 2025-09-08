package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// DataDogHook sends logs to DataDog
type DataDogHook struct {
	APIKey   string
	Region   string
	Client   *http.Client
	endpoint string
}

// NewDataDogHook creates a new DataDog hook
func NewDataDogHook(apiKey, region string) *DataDogHook {
	var endpoint string
	if region == "eu" {
		endpoint = "https://http-intake.logs.eu.datadoghq.com/v1/input/" + apiKey
	} else {
		endpoint = "https://http-intake.logs.datadoghq.com/v1/input/" + apiKey
	}

	return &DataDogHook{
		APIKey:   apiKey,
		Region:   region,
		Client:   &http.Client{Timeout: 10 * time.Second},
		endpoint: endpoint,
	}
}

// Fire sends the log entry to DataDog
func (h *DataDogHook) Fire(e *Entry) {
	// Convert to DataDog format
	ddLog := map[string]interface{}{
		"timestamp": e.Time.Unix() * 1000, // DataDog expects milliseconds
		"level":     e.Level,
		"message":   e.Msg,
		"service":   "logx-app",
		"source":    "go",
	}

	// Add fields if they exist
	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			ddLog[k] = v
		}
	}

	// Add trace information if available
	if e.TraceID != "" {
		ddLog["dd.trace_id"] = e.TraceID
	}
	if e.SpanID != "" {
		ddLog["dd.span_id"] = e.SpanID
	}
	if e.Caller != "" {
		ddLog["caller"] = e.Caller
	}

	payload, err := json.Marshal(ddLog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "datadog hook encode error: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", h.endpoint, bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintf(os.Stderr, "datadog hook request error: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	go func() {
		resp, err := h.Client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "datadog hook send error: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			fmt.Fprintf(os.Stderr, "datadog hook response error: %d\n", resp.StatusCode)
		}
	}()
}