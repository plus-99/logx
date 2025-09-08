package main

import (
        "context"
        "time"

        "github.com/plus99/logx"
)

func main() {
        l := logx.New()
        l.SetEncoder(logx.ConsoleFormatter{FullTimestamp: true, WithColors: false})
        l.SetLevel(logx.InfoLevel)
        l.SetReportCaller(true)
        
        // rotation hook
        h := logx.NewRotationHook("app.log", 50, 7, 30)
        l.AddHook(h)

        ctx := logx.ContextWithTraceSpan(context.Background(), "trace-123", "span-1")
        lg := l.WithContext(ctx)
        lg.WithFields(logx.Fields{"service": "database", "connections": 10}).Info("Database connection pool initialized successfully")

        // Test different log levels
        logx.Info("This is an info message")
        logx.Warn("This is a warning message")
        logx.Error("This is an error message")
        logx.Debug("This debug message won't show (level is Info)")

        // Test with fields
        logx.WithFields(logx.Fields{"user": "john", "action": "login"}).Info("User activity")

        // give hooks time to flush if async
        time.Sleep(100 * time.Millisecond)
}