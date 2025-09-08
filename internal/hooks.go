package logx

import (
        "github.com/plus99/logx/internal/hooks"
)

// FileHook writes logs to a file (append). It is a simple hook; for rotation use RotationHook.
type FileHook struct {
        internal *hooks.FileHook
}

func NewFileHook(path string) (*FileHook, error) {
        internal, err := hooks.NewFileHook(path)
        if err != nil {
                return nil, err
        }
        return &FileHook{internal: internal}, nil
}

func (h *FileHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}

// HTTPHook - for Loki-like or generic HTTP ingestion
type HTTPHook struct {
        internal *hooks.HTTPHook
}

func NewHTTPHook(endpoint string) *HTTPHook {
        internal := hooks.NewHTTPHook(endpoint)
        return &HTTPHook{internal: internal}
}

func (h *HTTPHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}

// RotationHook integrates lumberjack for rotation
type RotationHook struct {
        internal *hooks.RotationHook
}

func NewRotationHook(path string, maxSizeMB, maxBackups int, maxAgeDays int) *RotationHook {
        internal := hooks.NewRotationHook(path, maxSizeMB, maxBackups, maxAgeDays)
        return &RotationHook{internal: internal}
}

func (h *RotationHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}

// DataDogHook sends logs to DataDog
type DataDogHook struct {
        internal *hooks.DataDogHook
}

func NewDataDogHook(apiKey, region string) *DataDogHook {
        internal := hooks.NewDataDogHook(apiKey, region)
        return &DataDogHook{internal: internal}
}

func (h *DataDogHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}

// LogglyHook sends logs to Loggly
type LogglyHook struct {
        internal *hooks.LogglyHook
}

func NewLogglyHook(token, tag string) *LogglyHook {
        internal := hooks.NewLogglyHook(token, tag)
        return &LogglyHook{internal: internal}
}

func (h *LogglyHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}

// NewRelicHook sends logs to New Relic
type NewRelicHook struct {
        internal *hooks.NewRelicHook
}

func NewNewRelicHook(licenseKey, region string) *NewRelicHook {
        internal := hooks.NewNewRelicHook(licenseKey, region)
        return &NewRelicHook{internal: internal}
}

func (h *NewRelicHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}

// AtatusHook sends logs to Atatus
type AtatusHook struct {
        internal *hooks.AtatusHook
}

func NewAtatusHook(licenseKey, appName string) *AtatusHook {
        internal := hooks.NewAtatusHook(licenseKey, appName)
        return &AtatusHook{internal: internal}
}

func (h *AtatusHook) Fire(e *Entry) {
        internalEntry := &hooks.Entry{
                Time:    e.Time,
                Level:   e.Level,
                Msg:     e.Msg,
                Fields:  e.Fields,
                Caller:  e.Caller,
                TraceID: e.TraceID,
                SpanID:  e.SpanID,
        }
        h.internal.Fire(internalEntry)
}