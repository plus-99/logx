package encoding

import (
	"encoding/json"
	"time"
)

// Entry represents a log entry for internal encoding
type Entry struct {
	Time    time.Time              `json:"time"`
	Level   string                 `json:"level"`
	Msg     string                 `json:"msg"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Caller  string                 `json:"caller,omitempty"`
	TraceID string                 `json:"trace_id,omitempty"`
	SpanID  string                 `json:"span_id,omitempty"`
}

// JSONFormatter formats log entries as JSON
type JSONFormatter struct {
	TimestampFormat string
}

func (f JSONFormatter) Encode(e *Entry) ([]byte, error) {
	// Ensure fields deterministic order for tests
	out := make(map[string]interface{})
	out["time"] = e.Time.Format(f.TimestampFormat)
	out["level"] = e.Level
	out["msg"] = e.Msg
	if len(e.Fields) > 0 {
		// Copy fields
		m := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			m[k] = v
		}
		out["fields"] = m
	}
	if e.Caller != "" {
		out["caller"] = e.Caller
	}
	if e.TraceID != "" {
		out["trace_id"] = e.TraceID
	}
	if e.SpanID != "" {
		out["span_id"] = e.SpanID
	}
	return json.Marshal(out)
}