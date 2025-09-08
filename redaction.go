package logx

import (
        "encoding/json"
        "regexp"
        "strings"
        "sync"
)

// SecretString wraps a string that should always be redacted in logs
type SecretString struct {
        value string
}

// SecretBytes wraps bytes that should always be redacted in logs
type SecretBytes struct {
        value []byte
}

// NewSecretString creates a new SecretString
func NewSecretString(s string) SecretString {
        return SecretString{value: s}
}

// NewSecretBytes creates a new SecretBytes
func NewSecretBytes(b []byte) SecretBytes {
        return SecretBytes{value: b}
}

// String implements fmt.Stringer - always returns [REDACTED]
func (s SecretString) String() string {
        return "[REDACTED]"
}

// MarshalJSON implements json.Marshaler - always returns [REDACTED]
func (s SecretString) MarshalJSON() ([]byte, error) {
        return json.Marshal("[REDACTED]")
}

// String implements fmt.Stringer - always returns [REDACTED]
func (s SecretBytes) String() string {
        return "[REDACTED]"
}

// MarshalJSON implements json.Marshaler - always returns [REDACTED]
func (s SecretBytes) MarshalJSON() ([]byte, error) {
        return json.Marshal("[REDACTED]")
}

// GetValue returns the actual value (for internal use only)
func (s SecretString) GetValue() string {
        return s.value
}

// GetValue returns the actual value (for internal use only)
func (s SecretBytes) GetValue() []byte {
        return s.value
}

// RedactorFunc is a function that redacts sensitive data
type RedactorFunc func(key string, val interface{}) interface{}

// MessageRedactorFunc is a function that redacts sensitive data from messages
type MessageRedactorFunc func(msg string) string

// Global redaction configuration
var (
        redactionMutex        sync.RWMutex
        redactionEnabled      = true
        messageRedactionOn    = true
        keyRedactors         = make(map[string]bool)
        regexRedactors       = make([]*regexp.Regexp, 0)
        customRedactors      = make([]RedactorFunc, 0)
        messageRedactors     = make([]MessageRedactorFunc, 0)
)

// Built-in regex patterns for common sensitive data
var builtInPatterns = []*regexp.Regexp{
        // Credit Card Numbers (basic patterns)
        regexp.MustCompile(`\b(?:\d{4}[-\s]?){3}\d{4}\b`),
        // Social Security Numbers
        regexp.MustCompile(`\b\d{3}-?\d{2}-?\d{4}\b`),
        // API Keys (common patterns)
        regexp.MustCompile(`(?i)(?:api[_-]?key|apikey|token|secret)=[\w\-_]{8,}`),
        regexp.MustCompile(`(?i)bearer\s+[\w\-_\.]{20,}`),
        // Email addresses
        regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
        // AWS Access Keys
        regexp.MustCompile(`(?i)AKIA[0-9A-Z]{16}`),
        // Private Keys (PEM format headers)
        regexp.MustCompile(`-----BEGIN\s+(?:RSA\s+)?PRIVATE\s+KEY-----`),
        // Common password patterns
        regexp.MustCompile(`(?i)password=[\w\-_!@#$%^&*()]{4,}`),
}

func init() {
        // Initialize built-in redactors
        regexRedactors = append(regexRedactors, builtInPatterns...)
        
        // Add default key redactors
        keyRedactors["password"] = true
        keyRedactors["passwd"] = true
        keyRedactors["secret"] = true
        keyRedactors["token"] = true
        keyRedactors["apikey"] = true
        keyRedactors["api_key"] = true
        keyRedactors["access_token"] = true
        keyRedactors["refresh_token"] = true
        keyRedactors["private_key"] = true
        keyRedactors["ssn"] = true
        keyRedactors["social_security"] = true
        keyRedactors["credit_card"] = true
        keyRedactors["cc_number"] = true
}

// EnableRedaction enables or disables global redaction
func EnableRedaction(enabled bool) {
        redactionMutex.Lock()
        defer redactionMutex.Unlock()
        redactionEnabled = enabled
}

// EnableMessageRedaction enables or disables message redaction
func EnableMessageRedaction(enabled bool) {
        redactionMutex.Lock()
        defer redactionMutex.Unlock()
        messageRedactionOn = enabled
}

// AddKeyRedactor adds keys that should be redacted
func AddKeyRedactor(keys ...string) {
        redactionMutex.Lock()
        defer redactionMutex.Unlock()
        for _, key := range keys {
                keyRedactors[strings.ToLower(key)] = true
        }
}

// AddRegexRedactor adds regex patterns for redaction
func AddRegexRedactor(patterns ...string) {
        redactionMutex.Lock()
        defer redactionMutex.Unlock()
        for _, pattern := range patterns {
                if regex, err := regexp.Compile(pattern); err == nil {
                        regexRedactors = append(regexRedactors, regex)
                }
        }
}

// AddCustomRedactor adds a custom redaction function
func AddCustomRedactor(fn RedactorFunc) {
        redactionMutex.Lock()
        defer redactionMutex.Unlock()
        customRedactors = append(customRedactors, fn)
}

// AddMessageRedactor adds a message redaction function
func AddMessageRedactor(fn MessageRedactorFunc) {
        redactionMutex.Lock()
        defer redactionMutex.Unlock()
        messageRedactors = append(messageRedactors, fn)
}

// Mask recursively masks sensitive fields in maps/structs
func Mask(data map[string]interface{}, sensitiveKeys []string) map[string]interface{} {
        if data == nil {
                return nil
        }
        
        result := make(map[string]interface{})
        sensitiveSet := make(map[string]bool)
        for _, key := range sensitiveKeys {
                sensitiveSet[strings.ToLower(key)] = true
        }
        
        for key, value := range data {
                if sensitiveSet[strings.ToLower(key)] {
                        result[key] = "[REDACTED]"
                } else {
                        result[key] = maskValue(value, sensitiveSet)
                }
        }
        
        return result
}

func maskValue(value interface{}, sensitiveKeys map[string]bool) interface{} {
        if value == nil {
                return nil
        }
        
        switch v := value.(type) {
        case map[string]interface{}:
                result := make(map[string]interface{})
                for key, val := range v {
                        if sensitiveKeys[strings.ToLower(key)] {
                                result[key] = "[REDACTED]"
                        } else {
                                result[key] = maskValue(val, sensitiveKeys)
                        }
                }
                return result
        case []interface{}:
                result := make([]interface{}, len(v))
                for i, val := range v {
                        result[i] = maskValue(val, sensitiveKeys)
                }
                return result
        default:
                return value
        }
}

// applyRedaction applies all configured redaction rules
func applyRedaction(fields Fields) Fields {
        redactionMutex.RLock()
        enabled := redactionEnabled
        redactionMutex.RUnlock()
        
        if !enabled || fields == nil {
                return fields
        }
        
        redacted := make(Fields)
        
        redactionMutex.RLock()
        defer redactionMutex.RUnlock()
        
        for key, value := range fields {
                // Check for SecretString/SecretBytes (always redacted)
                switch v := value.(type) {
                case SecretString:
                        redacted[key] = "[REDACTED]"
                        continue
                case SecretBytes:
                        redacted[key] = "[REDACTED]"
                        continue
                }
                
                // Apply key-based redaction
                if keyRedactors[strings.ToLower(key)] {
                        redacted[key] = "[REDACTED]"
                        continue
                }
                
                // Apply custom redactors
                redactedValue := value
                for _, redactor := range customRedactors {
                        redactedValue = redactor(key, redactedValue)
                }
                
                // Apply regex redaction to string values
                if strValue, ok := redactedValue.(string); ok {
                        for _, regex := range regexRedactors {
                                if regex.MatchString(strValue) {
                                        redactedValue = "[REDACTED]"
                                        break
                                }
                        }
                }
                
                redacted[key] = redactedValue
        }
        
        return redacted
}

// applyMessageRedaction applies redaction to log messages
func applyMessageRedaction(msg string) string {
        redactionMutex.RLock()
        enabled := redactionEnabled && messageRedactionOn
        redactionMutex.RUnlock()
        
        if !enabled {
                return msg
        }
        
        result := msg
        
        redactionMutex.RLock()
        defer redactionMutex.RUnlock()
        
        // Apply custom message redactors
        for _, redactor := range messageRedactors {
                result = redactor(result)
        }
        
        // Apply regex redaction
        for _, regex := range regexRedactors {
                result = regex.ReplaceAllString(result, "[REDACTED]")
        }
        
        return result
}

// shouldRedact determines if redaction should be applied based on logger and global settings
func shouldRedact(loggerRedaction *bool) bool {
        if loggerRedaction != nil {
                return *loggerRedaction
        }
        redactionMutex.RLock()
        defer redactionMutex.RUnlock()
        return redactionEnabled
}

// isRedactionEnabled returns whether redaction is currently enabled globally
func isRedactionEnabled() bool {
        redactionMutex.RLock()
        defer redactionMutex.RUnlock()
        return redactionEnabled
}