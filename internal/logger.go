package internal

import (
        "context"
        "fmt"
        "io"
        "os"
        "runtime"
        "sync"
        "time"
)

// Level represents logging severity
type Level int

const (
        TraceLevel Level = iota
        DebugLevel
        InfoLevel
        WarnLevel
        ErrorLevel
        PanicLevel
        FatalLevel
)

func (l Level) String() string {
        switch l {
        case TraceLevel:
                return "TRACE"
        case DebugLevel:
                return "DEBUG"
        case InfoLevel:
                return "INFO"
        case WarnLevel:
                return "WARN"
        case ErrorLevel:
                return "ERROR"
        case PanicLevel:
                return "PANIC"
        case FatalLevel:
                return "FATAL"
        default:
                return "UNKNOWN"
        }
}

// Fields for structured logging
type Fields map[string]interface{}

// Entry is the log record
type Entry struct {
        Time      time.Time       `json:"time"`
        Level     string          `json:"level"`
        Msg       string          `json:"msg"`
        Fields    Fields          `json:"fields,omitempty"`
        Caller    string          `json:"caller,omitempty"`
        TraceID   string          `json:"trace_id,omitempty"`
        SpanID    string          `json:"span_id,omitempty"`
}

// Encoder formats an entry into a byte slice
type Encoder interface {
        Encode(e *Entry) ([]byte, error)
}

// Hook is a function called with an entry for side-effects
type Hook interface {
        Fire(e *Entry)
}

// HookFunc adapter
type HookFunc func(e *Entry)

func (f HookFunc) Fire(e *Entry) { f(e) }

// Logger is the core logger
type Logger struct {
        mu               sync.RWMutex
        out              io.Writer
        encoder          Encoder
        level            Level
        hooks            []Hook
        withFields       Fields
        reportCaller     bool
        redactionEnabled *bool // nil means use global setting
        pool             sync.Pool
}

var std = New()

// New creates a new logger with defaults
func New() *Logger {
        l := &Logger{
                out:     os.Stdout,
                encoder: JSONFormatter{TimestampFormat: time.RFC3339Nano},
                level:   InfoLevel,
                withFields: make(Fields),
                pool: sync.Pool{
                        New: func() interface{} { return new(Entry) },
                },
        }
        return l
}

// Setters for global
func SetOutput(w io.Writer) { std.SetOutput(w) }
func SetEncoder(e Encoder) { std.SetEncoder(e) }
func SetLevel(l Level) { std.SetLevel(l) }
func SetReportCaller(b bool) { std.SetReportCaller(b) }
func AddHook(h Hook) { std.AddHook(h) }

// Methods
func (l *Logger) SetOutput(w io.Writer) {
        l.mu.Lock()
        defer l.mu.Unlock()
        l.out = w
}

func (l *Logger) SetEncoder(e Encoder) {
        l.mu.Lock()
        defer l.mu.Unlock()
        l.encoder = e
}

func (l *Logger) SetLevel(lv Level) {
        l.mu.Lock()
        defer l.mu.Unlock()
        l.level = lv
}

func (l *Logger) SetReportCaller(b bool) {
        l.mu.Lock()
        defer l.mu.Unlock()
        l.reportCaller = b
}

func (l *Logger) AddHook(h Hook) {
        l.mu.Lock()
        defer l.mu.Unlock()
        l.hooks = append(l.hooks, h)
}

// WithFields returns a derived logger with additional fields
func (l *Logger) WithFields(f Fields) *Logger {
        l.mu.RLock()
        defer l.mu.RUnlock()
        newFields := make(Fields, len(l.withFields)+len(f))
        for k, v := range l.withFields {
                newFields[k] = v
        }
        for k, v := range f {
                newFields[k] = v
        }
        return &Logger{
                out:              l.out,
                encoder:          l.encoder,
                level:            l.level,
                hooks:            l.hooks,
                withFields:       newFields,
                reportCaller:     l.reportCaller,
                redactionEnabled: l.redactionEnabled,
                pool:             l.pool,
        }
}

// WithContext extracts trace/span IDs from context and returns derived logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
        flds := Fields{}
        if ctx == nil {
                return l.WithFields(flds)
        }
        if tid := ctx.Value(contextKey("trace_id")); tid != nil {
                flds["trace_id"] = tid
        }
        if sid := ctx.Value(contextKey("span_id")); sid != nil {
                flds["span_id"] = sid
        }
        return l.WithFields(flds)
}

// WithRedaction returns a derived logger with redaction enabled/disabled
func (l *Logger) WithRedaction(enabled bool) *Logger {
        l.mu.RLock()
        defer l.mu.RUnlock()
        return &Logger{
                out:              l.out,
                encoder:          l.encoder,
                level:            l.level,
                hooks:            l.hooks,
                withFields:       l.withFields,
                reportCaller:     l.reportCaller,
                redactionEnabled: &enabled,
                pool:             l.pool,
        }
}

// helper to capture caller
func (l *Logger) caller() string {
        // skip 3 frames to get to the caller of public API
        pc, file, line, ok := runtime.Caller(3)
        if !ok {
                return ""
        }
        fn := runtime.FuncForPC(pc)
        if fn == nil {
                return fmt.Sprintf("%s:%d", file, line)
        }
        return fmt.Sprintf("%s:%d %s", file, line, fn.Name())
}

func (l *Logger) log(level Level, msg string, f Fields) {
        l.mu.RLock()
        if level < l.level {
                l.mu.RUnlock()
                return
        }
        encoder := l.encoder
        out := l.out
        hooks := append([]Hook(nil), l.hooks...)
        reportCaller := l.reportCaller
        baseFields := l.withFields
        redactionEnabled := l.redactionEnabled
        l.mu.RUnlock()

        // build entry
        ent := l.pool.Get().(*Entry)
        ent.Time = time.Now()
        ent.Level = level.String()
        
        // Apply redaction to message if enabled
        redactedMsg := msg
        if shouldRedact(redactionEnabled) {
                redactedMsg = applyMessageRedaction(msg)
        }
        ent.Msg = redactedMsg
        
        // merge fields
        fields := make(Fields, len(baseFields)+len(f))
        for k, v := range baseFields {
                fields[k] = v
        }
        for k, v := range f {
                fields[k] = v
        }
        
        // Apply redaction to fields if enabled
        if shouldRedact(redactionEnabled) {
                fields = applyRedaction(fields)
        }
        ent.Fields = fields
        if reportCaller {
                ent.Caller = l.caller()
        }
        // run hooks (non-blocking best-effort)
        for _, h := range hooks {
                // run sync for now, hooks can dispatch async themselves
                h.Fire(ent)
        }
        b, err := encoder.Encode(ent)
        if err == nil {
                // ensure trailing newline
                if len(b) == 0 || b[len(b)-1] != '\n' {
                        b = append(b, '\n')
                }
                out.Write(b)
        } else {
                fmt.Fprintf(os.Stderr, "logx: encode error: %v\n", err)
        }

        // release entry
        ent.Time = time.Time{}
        ent.Level = ""
        ent.Msg = ""
        ent.Fields = nil
        ent.Caller = ""
        ent.TraceID = ""
        ent.SpanID = ""
        l.pool.Put(ent)
}

// Convenience methods
func (l *Logger) Info(msg string)                     { l.log(InfoLevel, msg, nil) }
func (l *Logger) Infof(format string, args ...any)    { l.log(InfoLevel, fmt.Sprintf(format, args...), nil) }
func (l *Logger) Warn(msg string)                     { l.log(WarnLevel, msg, nil) }
func (l *Logger) Error(msg string)                    { l.log(ErrorLevel, msg, nil) }
func (l *Logger) Debug(msg string)                    { l.log(DebugLevel, msg, nil) }
func (l *Logger) Fatal(msg string)                    { l.log(FatalLevel, msg, nil); os.Exit(1) }
func (l *Logger) Panic(msg string)                    { l.log(PanicLevel, msg, nil); panic(msg) }

// Global wrappers
func WithFields(f Fields) *Logger { return std.WithFields(f) }
func WithContext(ctx context.Context) *Logger { return std.WithContext(ctx) }
func Info(msg string) { std.Info(msg) }
func Warn(msg string) { std.Warn(msg) }
func Error(msg string) { std.Error(msg) }
func Debug(msg string) { std.Debug(msg) }
func Fatal(msg string) { std.Fatal(msg) }