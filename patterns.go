package logx

// BuiltInPatterns contains regex pattern strings for detecting common sensitive data
// This array can be easily updated to add new patterns or modify existing ones
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