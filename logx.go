package logx

import (
	"github.com/plus-99/logx/internal"
)

// Re-export all public types from internal package

// Level represents logging severity
type Level = internal.Level

// Log levels
const (
	TraceLevel = internal.TraceLevel
	DebugLevel = internal.DebugLevel
	InfoLevel  = internal.InfoLevel
	WarnLevel  = internal.WarnLevel
	ErrorLevel = internal.ErrorLevel
	PanicLevel = internal.PanicLevel
	FatalLevel = internal.FatalLevel
)

// Fields for structured logging
type Fields = internal.Fields

// Entry represents a log entry
type Entry = internal.Entry

// Encoder interface for log formatting
type Encoder = internal.Encoder

// Hook interface for extending logging
type Hook = internal.Hook

// Logger is the main logging struct
type Logger = internal.Logger

// Formatter types
type JSONFormatter = internal.JSONFormatter
type ConsoleFormatter = internal.ConsoleFormatter

// Hook types
type FileHook = internal.FileHook
type HTTPHook = internal.HTTPHook
type RotationHook = internal.RotationHook
type DataDogHook = internal.DataDogHook
type LogglyHook = internal.LogglyHook
type NewRelicHook = internal.NewRelicHook
type AtatusHook = internal.AtatusHook

// Redaction types
type SecretString = internal.SecretString
type SecretBytes = internal.SecretBytes
type RedactorFunc = internal.RedactorFunc
type MessageRedactorFunc = internal.MessageRedactorFunc

// Re-export all public functions

// Logger functions
var New = internal.New
var SetOutput = internal.SetOutput
var SetEncoder = internal.SetEncoder
var SetLevel = internal.SetLevel
var SetReportCaller = internal.SetReportCaller
var AddHook = internal.AddHook

// Global logging functions
var WithFields = internal.WithFields
var WithContext = internal.WithContext
var Info = internal.Info
var Warn = internal.Warn
var Error = internal.Error
var Debug = internal.Debug
var Fatal = internal.Fatal

// Hook constructor functions
var NewFileHook = internal.NewFileHook
var NewHTTPHook = internal.NewHTTPHook
var NewRotationHook = internal.NewRotationHook
var NewDataDogHook = internal.NewDataDogHook
var NewLogglyHook = internal.NewLogglyHook
var NewNewRelicHook = internal.NewNewRelicHook
var NewAtatusHook = internal.NewAtatusHook

// Redaction functions
var NewSecretString = internal.NewSecretString
var NewSecretBytes = internal.NewSecretBytes
var EnableRedaction = internal.EnableRedaction
var EnableMessageRedaction = internal.EnableMessageRedaction
var AddKeyRedactor = internal.AddKeyRedactor
var AddRegexRedactor = internal.AddRegexRedactor
var AddCustomRedactor = internal.AddCustomRedactor
var AddMessageRedactor = internal.AddMessageRedactor
var Mask = internal.Mask

// Context functions
var ContextWithTraceSpan = internal.ContextWithTraceSpan
