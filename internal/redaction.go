package internal

import (
	"regexp"
	"strings"
	"sync"

	"github.com/plus-99/logx/internal/redaction"
)

// SecretString wraps a string that should always be redacted in logs
type SecretString = redaction.SecretString

// NewSecretString creates a new secret string wrapper
func NewSecretString(s string) SecretString {
	return redaction.NewSecretString(s)
}

// SecretBytes wraps bytes that should always be redacted in logs
type SecretBytes = redaction.SecretBytes

// NewSecretBytes creates a new secret bytes wrapper
func NewSecretBytes(b []byte) SecretBytes {
	return redaction.NewSecretBytes(b)
}

// RedactorFunc is a function that redacts sensitive data
type RedactorFunc = redaction.RedactorFunc

// MessageRedactorFunc is a function that redacts sensitive data from messages
type MessageRedactorFunc = redaction.MessageRedactorFunc

// Global redaction configuration
var (
	redactionMutex     sync.RWMutex
	redactionEnabled   = true
	messageRedactionOn = true
	keyRedactors       = make(map[string]bool)
	regexRedactors     = make([]*regexp.Regexp, 0)
	customRedactors    = make([]RedactorFunc, 0)
	messageRedactors   = make([]MessageRedactorFunc, 0)
)

// Built-in compiled regex patterns for common sensitive data
// Patterns are loaded from internal/redaction for easier maintenance
var builtInPatterns []*regexp.Regexp

func init() {
	// Initialize built-in patterns
	redaction.InitPatterns()
	builtInPatterns = redaction.CompiledPatterns

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

// EnableMessageRedaction enables or disables message redaction specifically
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
		if compiled, err := regexp.Compile(pattern); err == nil {
			regexRedactors = append(regexRedactors, compiled)
		}
	}
}

// AddCustomRedactor adds a custom redaction function
func AddCustomRedactor(fn RedactorFunc) {
	redactionMutex.Lock()
	defer redactionMutex.Unlock()
	customRedactors = append(customRedactors, fn)
}

// AddMessageRedactor adds a custom message redaction function
func AddMessageRedactor(fn MessageRedactorFunc) {
	redactionMutex.Lock()
	defer redactionMutex.Unlock()
	messageRedactors = append(messageRedactors, fn)
}

// Mask selectively redacts sensitive keys from a data map
func Mask(data map[string]interface{}, sensitiveKeys []string) map[string]interface{} {
	return redaction.MaskSensitiveData(data, sensitiveKeys)
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

// applyMessageRedaction applies redaction patterns to log messages
func applyMessageRedaction(msg string) string {
	redactionMutex.RLock()
	enabled := messageRedactionOn
	patterns := append([]*regexp.Regexp(nil), regexRedactors...)
	msgRedactors := append([]MessageRedactorFunc(nil), messageRedactors...)
	redactionMutex.RUnlock()

	if !enabled {
		return msg
	}

	return redaction.ApplyMessageRedaction(msg, patterns, msgRedactors)
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
		switch value.(type) {
		case SecretString:
			redacted[key] = "[REDACTED]"
			continue
		case SecretBytes:
			redacted[key] = "[REDACTED]"
			continue
		}

		// Apply redaction using internal package
		tempFields := map[string]interface{}{key: value}
		redactedFields := redaction.ApplyFieldRedaction(tempFields, keyRedactors, regexRedactors, customRedactors)
		redacted[key] = redactedFields[key]
	}

	return redacted
}
