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
        lg.WithFields(logx.Fields{"animal": "walrus", "size": 10}).Info("A group of walrus emerges from the ocean")

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