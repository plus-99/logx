package redaction

// SecretString wraps a string that should always be redacted in logs
type SecretString struct {
	value string
}

// NewSecretString creates a new secret string wrapper
func NewSecretString(s string) SecretString {
	return SecretString{value: s}
}

// String returns [REDACTED] for logging
func (s SecretString) String() string {
	return "[REDACTED]"
}

// GetValue returns the actual value (for internal use only)
func (s SecretString) GetValue() string {
	return s.value
}

// SecretBytes wraps bytes that should always be redacted in logs
type SecretBytes struct {
	value []byte
}

// NewSecretBytes creates a new secret bytes wrapper
func NewSecretBytes(b []byte) SecretBytes {
	return SecretBytes{value: b}
}

// String returns [REDACTED] for logging
func (s SecretBytes) String() string {
	return "[REDACTED]"
}

// GetValue returns the actual value (for internal use only)
func (s SecretBytes) GetValue() []byte {
	return s.value
}