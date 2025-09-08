package encoding

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ConsoleFormatter formats log entries in human-readable format
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