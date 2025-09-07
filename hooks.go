package logx

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// FileHook writes logs to a file (append). It is a simple hook; for rotation use RotationHook.
type FileHook struct {
	w  *os.File
	bw *bufio.Writer
}

func NewFileHook(path string) (*FileHook, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileHook{w: f, bw: bufio.NewWriter(f)}, nil
}

func (h *FileHook) Fire(e *Entry) {
	b, err := JSONFormatter{TimestampFormat: time.RFC3339Nano}.Encode(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "filehook encode err: %v\n", err)
		return
	}
	h.bw.Write(b)
	h.bw.WriteByte('\n')
	h.bw.Flush()
}

// HTTPHook - for Loki-like or generic HTTP ingestion
type HTTPHook struct {
	Endpoint string
	Client   *http.Client
}

func NewHTTPHook(endpoint string) *HTTPHook {
	return &HTTPHook{Endpoint: endpoint, Client: &http.Client{Timeout: 5 * time.Second}}
}

func (h *HTTPHook) Fire(e *Entry) {
	b, err := JSONFormatter{TimestampFormat: time.RFC3339Nano}.Encode(e)
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

// RotationHook integrates lumberjack for rotation
type RotationHook struct {
	lj *lumberjack.Logger
}

func NewRotationHook(path string, maxSizeMB, maxBackups int, maxAgeDays int) *RotationHook {
	lj := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   false,
	}
	return &RotationHook{lj: lj}
}

func (h *RotationHook) Fire(e *Entry) {
	b, err := JSONFormatter{TimestampFormat: time.RFC3339Nano}.Encode(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rotationhook encode err: %v\n", err)
		return
	}
	h.lj.Write(append(b, '\n'))
}