package internal

import (
	"github.com/plus-99/logx/internal/encoding"
)

// JSONFormatter implements Encoder
type JSONFormatter struct {
	TimestampFormat string
}

func (f JSONFormatter) Encode(e *Entry) ([]byte, error) {
	// Convert to internal entry format
	internalEntry := &encoding.Entry{
		Time:    e.Time,
		Level:   e.Level,
		Msg:     e.Msg,
		Fields:  e.Fields,
		Caller:  e.Caller,
		TraceID: e.TraceID,
		SpanID:  e.SpanID,
	}

	formatter := encoding.JSONFormatter{TimestampFormat: f.TimestampFormat}
	return formatter.Encode(internalEntry)
}

// ConsoleFormatter produces human-friendly lines
type ConsoleFormatter struct {
	FullTimestamp bool
	WithColors    bool
}

func (f ConsoleFormatter) Encode(e *Entry) ([]byte, error) {
	// Convert to internal entry format
	internalEntry := &encoding.Entry{
		Time:    e.Time,
		Level:   e.Level,
		Msg:     e.Msg,
		Fields:  e.Fields,
		Caller:  e.Caller,
		TraceID: e.TraceID,
		SpanID:  e.SpanID,
	}

	formatter := encoding.ConsoleFormatter{FullTimestamp: f.FullTimestamp, WithColors: f.WithColors}
	return formatter.Encode(internalEntry)
}
