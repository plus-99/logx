package main

import (
	"context"
	"os"
	"time"

	"github.com/plus99/logx"
)

func main() {
	// Configure logger for demo
	logger := logx.New()
	logger.SetEncoder(logx.JSONFormatter{TimestampFormat: time.RFC3339})
	logger.SetLevel(logx.InfoLevel)

	println("=== LogX Redaction Demo ===")
	println()

	// 1. Secret Wrappers Demo
	println("1. Secret Wrappers Demo:")
	email := logx.NewSecretString("user@company.com")
	password := logx.NewSecretString("supersecret123")
	
	logger.WithFields(logx.Fields{
		"email":    email,
		"password": password,
	}).Info("User registration")
	println()

	// 2. Field Redaction Demo
	println("2. Field Redaction Demo:")
	logx.AddKeyRedactor("api_key", "credit_card")
	
	logger.WithFields(logx.Fields{
		"username":    "johndoe",
		"api_key":     "sk_12345abcdef",
		"credit_card": "4111-1111-1111-1111",
	}).Info("Payment processing")
	println()

	// 3. Message Redaction Demo
	println("3. Message Redaction Demo:")
	logx.EnableMessageRedaction(true)
	
	logger.Info("Login attempt with password=mypassword123")
	logger.Warn("API call failed with key=SECRET_API_KEY_HERE")
	println()

	// 4. Custom Redactor Demo
	println("4. Custom Redactor Demo:")
	logx.AddCustomRedactor(func(key string, val interface{}) interface{} {
		if key == "ssn" {
			return "[MASKED-SSN]"
		}
		return val
	})
	
	logger.WithFields(logx.Fields{
		"ssn":  "123-45-6789",
		"name": "John Doe",
	}).Info("User profile update")
	println()

	// 5. Regex Redactor Demo
	println("5. Regex Redactor Demo:")
	logx.AddRegexRedactor(`(?i)token=[A-Za-z0-9-_]+`)
	
	logger.Error("Authentication failed: token=abc123def456")
	println()

	// 6. Mask Helper Demo
	println("6. Mask Helper Demo:")
	userData := map[string]interface{}{
		"username":    "alice",
		"password":    "secret123",
		"email":       "alice@example.com",
		"phone":       "555-123-4567",
	}
	
	maskedData := logx.Mask(userData, []string{"password", "email"})
	logger.WithFields(logx.Fields{
		"user_data": maskedData,
	}).Info("User profile retrieved")
	println()

	// 7. Per-Logger Redaction Control
	println("7. Per-Logger Redaction Control:")
	
	// Create a logger with redaction disabled (for local testing)
	devLogger := logger.WithRedaction(false)
	devLogger.WithFields(logx.Fields{
		"debug_token": "debug_12345",
		"test_key":    "test_secret",
	}).Info("Development debug info (redaction disabled)")
	
	// Regular logger still has redaction enabled
	logger.WithFields(logx.Fields{
		"prod_token": "prod_67890",
		"secret":     "production_secret",
	}).Info("Production logging (redaction enabled)")
	println()

	// 8. Built-in Patterns Demo
	println("8. Built-in Patterns Demo (automatic detection):")
	logger.Info("Processing credit card 4111-1111-1111-1111")
	logger.Info("SSN validation: 123-45-6789")
	logger.Info("Email notification: user@example.com")
	logger.Info("AWS key found: AKIAIOSFODNN7EXAMPLE")
	println()

	// 9. Environment-based Redaction
	println("9. Environment-based Redaction:")
	if os.Getenv("ENV") == "production" {
		logx.EnableRedaction(true)
		logger.Info("Production mode: redaction enabled")
	} else {
		// For demo, we'll show both modes
		logx.EnableRedaction(false)
		logger.Info("Development mode: redaction disabled - token=dev_token_123")
		
		logx.EnableRedaction(true)
		logger.Info("Production mode: redaction enabled - token=prod_token_456")
	}
	println()

	// 10. Context with Trace IDs + Redaction
	println("10. Context with Trace IDs + Redaction:")
	ctx := logx.ContextWithTraceSpan(context.Background(), "trace-789", "span-123")
	contextLogger := logger.WithContext(ctx)
	
	contextLogger.WithFields(logx.Fields{
		"operation": "user_auth",
		"token":     "auth_token_secret",
	}).Info("Authentication successful")
	println()

	println("=== Demo Complete ===")
}