package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// HTTPHook - for Loki-like or generic HTTP ingestion
type HTTPHook struct {
	Endpoint string
	Client   *http.Client
}

// NewHTTPHook creates a new HTTP hook
func NewHTTPHook(endpoint string) *HTTPHook {
	return &HTTPHook{Endpoint: endpoint, Client: &http.Client{Timeout: 5 * time.Second}}
}

// Fire sends the log entry to HTTP endpoint
func (h *HTTPHook) Fire(e *Entry) {
	b, err := json.Marshal(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httphook encode err: %v\n", err)
		return
	}
	req, _ := http.NewRequest("POST", h.Endpoint, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.Client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "httphook send err: %v\n", err)
		return
	}
	resp.Body.Close()
}