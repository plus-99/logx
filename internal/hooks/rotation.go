package hooks

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// RotationHook integrates lumberjack for rotation
type RotationHook struct {
	lj *lumberjack.Logger
}

// NewRotationHook creates a new rotation hook
func NewRotationHook(path string, maxSizeMB, maxBackups int, maxAgeDays int) *RotationHook {
	lj := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   true,
	}
	return &RotationHook{lj: lj}
}

// Fire writes the log entry with rotation
func (h *RotationHook) Fire(e *Entry) {
	b, err := json.Marshal(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "rotationhook encode err: %v\n", err)
		return
	}
	h.lj.Write(b)
	h.lj.Write([]byte("\n"))
}