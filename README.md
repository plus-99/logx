# LogX - High-Performance Go Logging Library

LogX is a fast, lightweight, and feature-rich logging library for Go applications. It's designed as a high-performance alternative to popular logging frameworks like logrus and zerolog, with a focus on simplicity, performance, and flexibility.

## Features

- **🚀 High Performance**: Optimized for minimal allocations and maximum throughput
- **📊 Structured Logging**: Support for key-value fields and JSON output
- **🎯 Multiple Log Levels**: TRACE, DEBUG, INFO, WARN, ERROR, PANIC, FATAL
- **🔄 Multiple Output Formats**: JSON and human-readable console formats
- **🪝 Extensible Hooks**: Built-in file, HTTP, and rotation hooks
- **🔗 Context Integration**: Extract trace/span IDs from Go context
- **🔒 Thread-Safe**: Safe for concurrent use across goroutines
- **📦 Log Rotation**: Integrated with lumberjack for automatic log rotation
- **🎨 Customizable**: Flexible encoders and output destinations

## Installation

```bash
go get github.com/plus99/logx
```

## Quick Start

```go
package main

import (
    "github.com/plus99/logx"
)

func main() {
    // Simple logging
    logx.Info("Hello, LogX!")
    logx.Warn("This is a warning")
    logx.Error("Something went wrong")

    // Structured logging with fields
    logx.WithFields(logx.Fields{
        "user": "john",
        "action": "login",
        "ip": "192.168.1.1",
    }).Info("User logged in")
}
```

## Usage Examples

### Basic Configuration

```go
package main

import (
    "github.com/plus99/logx"
)

func main() {
    // Create a new logger
    logger := logx.New()
    
    // Set log level
    logger.SetLevel(logx.InfoLevel)
    
    // Use console formatter for human-readable output
    logger.SetEncoder(logx.ConsoleFormatter{
        FullTimestamp: true,
        WithColors: false,
    })
    
    // Enable caller information
    logger.SetReportCaller(true)
    
    logger.Info("Logger configured successfully")
}
```

### Structured Logging

```go
// Log with structured fields
logx.WithFields(logx.Fields{
    "component": "database",
    "operation": "connect",
    "duration": "150ms",
}).Info("Database connection established")

// Chain multiple field sets
logger := logx.WithFields(logx.Fields{"service": "api"})
logger.WithFields(logx.Fields{"endpoint": "/users"}).Info("Handling request")
```

### Context Integration

```go
import (
    "context"
    "github.com/plus99/logx"
)

// Add trace/span IDs to context
ctx := logx.ContextWithTraceSpan(context.Background(), "trace-123", "span-456")

// Logger will automatically extract and include trace/span IDs
logger := logx.WithContext(ctx)
logger.Info("Processing request")
// Output: {"level":"INFO","msg":"Processing request","time":"...","fields":{"trace_id":"trace-123","span_id":"span-456"}}
```

### File Logging with Rotation

```go
import (
    "github.com/plus99/logx"
)

func main() {
    logger := logx.New()
    
    // Add file rotation hook
    rotationHook := logx.NewRotationHook(
        "app.log",  // filename
        10,         // max size in MB
        5,          // max number of backups
        30,         // max age in days
    )
    logger.AddHook(rotationHook)
    
    logger.Info("This will be written to app.log with rotation")
}
```

### HTTP Hook for Remote Logging

```go
// Send logs to remote endpoint (e.g., Loki, ELK)
httpHook := logx.NewHTTPHook("https://logs.example.com/api/logs")
logger.AddHook(httpHook)

logger.Error("This log will be sent to remote endpoint")
```

### Custom Formatters

```go
// JSON formatter with custom timestamp
logger.SetEncoder(logx.JSONFormatter{
    TimestampFormat: "2006-01-02 15:04:05.000",
})

// Console formatter with full timestamps
logger.SetEncoder(logx.ConsoleFormatter{
    FullTimestamp: true,
    WithColors: true,
})
```

## Log Levels

LogX supports the following log levels (in order of severity):

- `TRACE` - Very detailed information, typically only of interest when diagnosing problems
- `DEBUG` - Detailed information, typically only of interest when diagnosing problems
- `INFO` - General information about program execution
- `WARN` - Warning messages for potentially harmful situations
- `ERROR` - Error events that might still allow the application to continue
- `PANIC` - Severe error events that will cause the program to panic
- `FATAL` - Very severe error events that will cause the program to exit

```go
logger.SetLevel(logx.DebugLevel) // Only logs DEBUG and above will be output

logx.Trace("This won't be shown")
logx.Debug("This will be shown")
logx.Info("This will be shown")
logx.Error("This will be shown")
```

## Performance

LogX is designed for high performance with:

- **Object Pooling**: Reuses log entry objects to reduce garbage collection
- **Lazy Formatting**: Messages are only formatted if they meet the log level threshold
- **Minimal Allocations**: Optimized to reduce memory allocations during logging operations
- **Concurrent Safe**: Uses read-write mutexes for optimal concurrent performance

### Benchmarks

Run the included benchmarks to compare performance:

```bash
go test -bench=. -benchmem
```

Example results comparing LogX with popular alternatives:

```
BenchmarkLogxInfo-8      2000000    750 ns/op    120 B/op    3 allocs/op
BenchmarkLogrusInfo-8    1000000   1500 ns/op    280 B/op    8 allocs/op
BenchmarkZerologInfo-8   3000000    550 ns/op     96 B/op    2 allocs/op
```

## API Reference

### Logger Methods

```go
// Logger creation
logger := logx.New()

// Configuration
logger.SetLevel(level Level)
logger.SetEncoder(encoder Encoder)
logger.SetOutput(w io.Writer)
logger.SetReportCaller(enabled bool)
logger.AddHook(hook Hook)

// Logging methods
logger.Trace(msg string)
logger.Debug(msg string)
logger.Info(msg string)
logger.Warn(msg string)
logger.Error(msg string)
logger.Panic(msg string)  // Calls panic() after logging
logger.Fatal(msg string)  // Calls os.Exit(1) after logging

// Formatted logging
logger.Infof(format string, args ...interface{})

// Structured logging
logger.WithFields(fields Fields) *Logger
logger.WithContext(ctx context.Context) *Logger
```

### Global Functions

```go
// Global logger functions (use default logger)
logx.SetLevel(level Level)
logx.SetEncoder(encoder Encoder)
logx.SetOutput(w io.Writer)
logx.SetReportCaller(enabled bool)
logx.AddHook(hook Hook)

logx.Info(msg string)
logx.WithFields(fields Fields) *Logger
logx.WithContext(ctx context.Context) *Logger
```

### Fields Type

```go
type Fields map[string]interface{}

// Usage
fields := logx.Fields{
    "user_id": 12345,
    "action": "login",
    "success": true,
    "duration": 150.5,
}
```

## Built-in Hooks

### File Hook
```go
fileHook, err := logx.NewFileHook("/var/log/app.log")
if err != nil {
    log.Fatal(err)
}
logger.AddHook(fileHook)
```

### Rotation Hook (with Lumberjack)
```go
rotationHook := logx.NewRotationHook(
    "/var/log/app.log", // filename
    100,                // max size in MB
    10,                 // max backups
    30,                 // max age in days
)
logger.AddHook(rotationHook)
```

### HTTP Hook
```go
httpHook := logx.NewHTTPHook("https://logs.example.com/api/v1/logs")
logger.AddHook(httpHook)
```

## Requirements

- Go 1.19 or later

## Dependencies

- [lumberjack.v2](https://gopkg.in/natefinch/lumberjack.v2) - Log rotation
- [logrus](https://github.com/sirupsen/logrus) - Benchmarking only
- [zerolog](https://github.com/rs/zerolog) - Benchmarking only

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.