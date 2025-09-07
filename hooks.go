package logx

import (
        "bufio"
        "bytes"
        "encoding/json"
        "fmt"
        "net/http"
        "os"
        "strings"
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

// DataDogHook sends logs to DataDog
type DataDogHook struct {
        APIKey   string
        Region   string // us, eu, us3, us5, ap1, gov
        Source   string
        Service  string
        Hostname string
        Tags     string
        Client   *http.Client
}

func NewDataDogHook(apiKey, region string) *DataDogHook {
        return &DataDogHook{
                APIKey: apiKey,
                Region: region,
                Client: &http.Client{Timeout: 10 * time.Second},
        }
}

func (h *DataDogHook) getEndpoint() string {
        switch h.Region {
        case "eu":
                return "https://http-intake.logs.datadoghq.eu/v1/input"
        case "us3":
                return "https://http-intake.logs.us3.datadoghq.com/v1/input"
        case "us5":
                return "https://http-intake.logs.us5.datadoghq.com/v1/input"
        case "ap1":
                return "https://http-intake.logs.ap1.datadoghq.com/v1/input"
        case "gov":
                return "https://http-intake.logs.ddog-gov.com/v1/input"
        default: // us
                return "https://http-intake.logs.datadoghq.com/v1/input"
        }
}

func (h *DataDogHook) Fire(e *Entry) {
        payload := map[string]interface{}{
                "message": e.Msg,
                "level":   strings.ToLower(e.Level),
        }

        if h.Source != "" {
                payload["ddsource"] = h.Source
        }
        if h.Service != "" {
                payload["service"] = h.Service
        }
        if h.Hostname != "" {
                payload["hostname"] = h.Hostname
        }
        if h.Tags != "" {
                payload["ddtags"] = h.Tags
        }

        // Add custom fields
        if len(e.Fields) > 0 {
                for k, v := range e.Fields {
                        payload[k] = v
                }
        }

        b, err := json.Marshal(payload)
        if err != nil {
                fmt.Fprintf(os.Stderr, "datadoghook encode err: %v\n", err)
                return
        }

        req, _ := http.NewRequest("POST", h.getEndpoint(), bytes.NewReader(b))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("DD-API-KEY", h.APIKey)

        resp, err := h.Client.Do(req)
        if err != nil {
                fmt.Fprintf(os.Stderr, "datadoghook send err: %v\n", err)
                return
        }
        resp.Body.Close()
}

// LogglyHook sends logs to Loggly
type LogglyHook struct {
        Token  string
        Tag    string
        Client *http.Client
}

func NewLogglyHook(token, tag string) *LogglyHook {
        return &LogglyHook{
                Token:  token,
                Tag:    tag,
                Client: &http.Client{Timeout: 10 * time.Second},
        }
}

func (h *LogglyHook) Fire(e *Entry) {
        payload := map[string]interface{}{
                "timestamp": e.Time.Format(time.RFC3339Nano),
                "level":     e.Level,
                "message":   e.Msg,
        }

        // Add custom fields
        if len(e.Fields) > 0 {
                for k, v := range e.Fields {
                        payload[k] = v
                }
        }

        if e.Caller != "" {
                payload["caller"] = e.Caller
        }

        b, err := json.Marshal(payload)
        if err != nil {
                fmt.Fprintf(os.Stderr, "logglyhook encode err: %v\n", err)
                return
        }

        endpoint := fmt.Sprintf("https://logs-01.loggly.com/inputs/%s/tag/%s/", h.Token, h.Tag)
        req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(b))
        req.Header.Set("Content-Type", "application/json")

        resp, err := h.Client.Do(req)
        if err != nil {
                fmt.Fprintf(os.Stderr, "logglyhook send err: %v\n", err)
                return
        }
        resp.Body.Close()
}

// NewRelicHook sends logs to New Relic
type NewRelicHook struct {
        LicenseKey string
        Region     string // us, eu, gov
        Service    string
        Hostname   string
        Client     *http.Client
}

func NewNewRelicHook(licenseKey, region string) *NewRelicHook {
        return &NewRelicHook{
                LicenseKey: licenseKey,
                Region:     region,
                Client:     &http.Client{Timeout: 10 * time.Second},
        }
}

func (h *NewRelicHook) getEndpoint() string {
        switch h.Region {
        case "eu":
                return "https://log-api.eu.newrelic.com/log/v1"
        case "gov":
                return "https://gov-log-api.newrelic.com/log/v1"
        default: // us
                return "https://log-api.newrelic.com/log/v1"
        }
}

func (h *NewRelicHook) Fire(e *Entry) {
        payload := map[string]interface{}{
                "timestamp": e.Time.UnixMilli(),
                "message":   e.Msg,
                "level":     strings.ToLower(e.Level),
        }

        if h.Service != "" {
                payload["service"] = h.Service
        }
        if h.Hostname != "" {
                payload["hostname"] = h.Hostname
        }

        // Add custom fields as attributes
        if len(e.Fields) > 0 {
                attributes := make(map[string]interface{})
                for k, v := range e.Fields {
                        attributes[k] = v
                }
                payload["attributes"] = attributes
        }

        if e.Caller != "" {
                if payload["attributes"] == nil {
                        payload["attributes"] = make(map[string]interface{})
                }
                payload["attributes"].(map[string]interface{})["caller"] = e.Caller
        }

        b, err := json.Marshal(payload)
        if err != nil {
                fmt.Fprintf(os.Stderr, "newrelichook encode err: %v\n", err)
                return
        }

        req, _ := http.NewRequest("POST", h.getEndpoint(), bytes.NewReader(b))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Api-Key", h.LicenseKey)

        resp, err := h.Client.Do(req)
        if err != nil {
                fmt.Fprintf(os.Stderr, "newrelichook send err: %v\n", err)
                return
        }
        resp.Body.Close()
}

// AtatusHook sends logs to Atatus
type AtatusHook struct {
        LicenseKey string
        AppName    string
        Client     *http.Client
}

func NewAtatusHook(licenseKey, appName string) *AtatusHook {
        return &AtatusHook{
                LicenseKey: licenseKey,
                AppName:    appName,
                Client:     &http.Client{Timeout: 10 * time.Second},
        }
}

func (h *AtatusHook) Fire(e *Entry) {
        payload := map[string]interface{}{
                "timestamp": e.Time.Format(time.RFC3339Nano),
                "level":     strings.ToLower(e.Level),
                "message":   e.Msg,
                "app":       h.AppName,
        }

        // Add custom fields
        if len(e.Fields) > 0 {
                for k, v := range e.Fields {
                        payload[k] = v
                }
        }

        if e.Caller != "" {
                payload["caller"] = e.Caller
        }

        b, err := json.Marshal(payload)
        if err != nil {
                fmt.Fprintf(os.Stderr, "atatushook encode err: %v\n", err)
                return
        }

        // Note: Atatus API endpoint may need to be updated based on their actual API
        endpoint := "https://api.atatus.com/v2/logs"
        req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(b))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+h.LicenseKey)

        resp, err := h.Client.Do(req)
        if err != nil {
                fmt.Fprintf(os.Stderr, "atatushook send err: %v\n", err)
                return
        }
        resp.Body.Close()
}