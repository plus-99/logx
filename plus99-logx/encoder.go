package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// JSONFormatter implements Encoder
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

// ConsoleFormatter produces human-friendly lines
type ConsoleFormatter struct {
	FullTimestamp bool
	WithColors    bool
}

func (f ConsoleFormatter) Encode(e *Entry) ([]byte, error) {
	var buf bytes.Buffer
	if f.FullTimestamp {
		buf.WriteString(e.Time.Format(time.RFC3339))
	} else {
		buf.WriteString(e.Time.Format("15:04:05"))
	}
	buf.WriteString(" ")
	buf.WriteString("[")
	buf.WriteString(e.Level)
	buf.WriteString("] ")
	buf.WriteString(e.Msg)
	if len(e.Fields) > 0 {
		// deterministic order
		keys := make([]string, 0, len(e.Fields))
		for k := range e.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, fmt.Sprintf("%s=%v", k, e.Fields[k]))
		}
		buf.WriteString(" ")
		buf.WriteString(strings.Join(pairs, " "))
	}
	if e.Caller != "" {
		buf.WriteString(" ")
		buf.WriteString("(" + e.Caller + ")")
	}
	return buf.Bytes(), nil
}