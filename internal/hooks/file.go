package hooks

import (
	"encoding/json"
	"os"
	"time"
)

// Entry represents a log entry for internal hook processing
type Entry struct {
	Time    time.Time              `json:"time"`
	Level   string                 `json:"level"`
	Msg     string                 `json:"msg"`
	Fields  map[string]interface{} `json:"fields,omitempty"`
	Caller  string                 `json:"caller,omitempty"`
	TraceID string                 `json:"trace_id,omitempty"`
	SpanID  string                 `json:"span_id,omitempty"`
}

// FileHook writes log entries to a file
type FileHook struct {
	filename string
	file     *os.File
}

// NewFileHook creates a new file hook
func NewFileHook(path string) (*FileHook, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FileHook{
		filename: path,
		file:     file,
	}, nil
}

// Fire writes the log entry to the file
func (h *FileHook) Fire(e *Entry) {
	if h.file != nil {
		data, err := json.Marshal(e)
		if err == nil {
			h.file.Write(data)
			h.file.WriteString("\n")
		}
	}
}

// Close closes the file
func (h *FileHook) Close() error {
	if h.file != nil {
		return h.file.Close()
	}
	return nil
}