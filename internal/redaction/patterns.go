package redaction

import (
	"regexp"
	"strings"
)

// BuiltInPatterns contains regex pattern strings for detecting common sensitive data
var BuiltInPatterns = []string{
	// Credit Card Numbers (basic patterns)
	`\b(?:\d{4}[-\s]?){3}\d{4}\b`,
	
	// Social Security Numbers
	`\b\d{3}-?\d{2}-?\d{4}\b`,
	
	// API Keys (common patterns)
	`(?i)(?:api[_-]?key|apikey|token|secret)=[\w\-_]{8,}`,
	`(?i)bearer\s+[\w\-_\.]{20,}`,
	
	// Email addresses
	`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`,
	
	// AWS Access Keys
	`(?i)AKIA[0-9A-Z]{16}`,
	
	// Private key headers
	`-----BEGIN\s+(?:RSA\s+)?PRIVATE\s+KEY-----`,
	
	// Password patterns
	`(?i)password=[\w\-_!@#$%^&*()]{4,}`,
}

// CompiledPatterns holds the compiled regex patterns
var CompiledPatterns []*regexp.Regexp

// RedactorFunc is a function that redacts sensitive data
type RedactorFunc func(key string, val interface{}) interface{}

// MessageRedactorFunc is a function that redacts sensitive data from messages
type MessageRedactorFunc func(msg string) string

// InitPatterns compiles all built-in patterns
func InitPatterns() {
	CompiledPatterns = make([]*regexp.Regexp, 0, len(BuiltInPatterns))
	for _, pattern := range BuiltInPatterns {
		compiled := regexp.MustCompile(pattern)
		CompiledPatterns = append(CompiledPatterns, compiled)
	}
}

// ApplyMessageRedaction applies redaction patterns to a message
func ApplyMessageRedaction(msg string, regexPatterns []*regexp.Regexp, messageRedactors []MessageRedactorFunc) string {
	result := msg
	
	// Apply built-in patterns
	for _, pattern := range CompiledPatterns {
		if pattern.MatchString(result) {
			result = pattern.ReplaceAllString(result, "[REDACTED]")
		}
	}
	
	// Apply additional regex patterns
	for _, pattern := range regexPatterns {
		if pattern.MatchString(result) {
			result = pattern.ReplaceAllString(result, "[REDACTED]")
		}
	}
	
	// Apply custom message redactors
	for _, redactor := range messageRedactors {
		result = redactor(result)
	}
	
	return result
}

// ApplyFieldRedaction applies redaction to field values
func ApplyFieldRedaction(fields map[string]interface{}, keyRedactors map[string]bool, regexPatterns []*regexp.Regexp, customRedactors []RedactorFunc) map[string]interface{} {
	if fields == nil {
		return fields
	}
	
	redacted := make(map[string]interface{})
	
	for key, value := range fields {
		// Check for SecretString/SecretBytes (handled by caller)
		
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
			for _, regex := range regexPatterns {
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

// MaskSensitiveData masks sensitive fields in a data structure
func MaskSensitiveData(data map[string]interface{}, sensitiveKeys []string) map[string]interface{} {
	if data == nil {
		return nil
	}
	
	// Create lookup map for sensitive keys
	sensitive := make(map[string]bool)
	for _, key := range sensitiveKeys {
		sensitive[strings.ToLower(key)] = true
	}
	
	return maskValue(data, sensitive).(map[string]interface{})
}

// maskValue recursively masks values in data structures
func maskValue(value interface{}, sensitiveKeys map[string]bool) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, val := range v {
			if sensitiveKeys[strings.ToLower(k)] {
				result[k] = "[REDACTED]"
			} else {
				result[k] = maskValue(val, sensitiveKeys)
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