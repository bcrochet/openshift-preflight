package retry

import (
	"context"
	"time"

	"github.com/go-logr/logr"
)

// WithRetry executes the provided function with exponential backoff retry logic.
// It retries on network errors and temporary failures, with a maximum of 3 attempts.
func WithRetry[T any](ctx context.Context, operation func() (T, error)) (T, error) {
	return WithRetryConfig(ctx, operation, Config{
		MaxAttempts: 3,
		BaseDelay:   time.Second,
		MaxDelay:    5 * time.Second,
	})
}

// WithRetry2 executes the provided function with exponential backoff retry logic for functions returning 2 values plus error.
func WithRetry2[T1, T2 any](ctx context.Context, operation func() (T1, T2, error)) (T1, T2, error) {
	return WithRetryConfig2(ctx, operation, Config{
		MaxAttempts: 3,
		BaseDelay:   time.Second,
		MaxDelay:    5 * time.Second,
	})
}

// Config defines retry behavior
type Config struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// WithRetryConfig executes the provided function with configurable retry logic.
func WithRetryConfig[T any](ctx context.Context, operation func() (T, error), config Config) (T, error) {
	logger := logr.FromContextOrDiscard(ctx)
	
	var result T
	var lastErr error
	
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		result, err := operation()
		if err == nil {
			return result, nil
		}
		
		lastErr = err
		
		// Don't retry on the last attempt
		if attempt == config.MaxAttempts {
			break
		}
		
		// Check if error is retryable (network errors, temporary failures)
		if !isRetryableError(err) {
			logger.V(1).Info("non-retryable error encountered", "error", err, "attempt", attempt)
			return result, err
		}
		
		// Calculate delay with exponential backoff
		delay := time.Duration(attempt) * config.BaseDelay
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
		
		logger.V(1).Info("retrying operation after error", "error", err, "attempt", attempt, "delay", delay)
		
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
	
	return result, lastErr
}

// WithRetryConfig2 executes the provided function with configurable retry logic for functions returning 2 values plus error.
func WithRetryConfig2[T1, T2 any](ctx context.Context, operation func() (T1, T2, error), config Config) (T1, T2, error) {
	logger := logr.FromContextOrDiscard(ctx)
	
	var result1 T1
	var result2 T2
	var lastErr error
	
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		result1, result2, err := operation()
		if err == nil {
			return result1, result2, nil
		}
		
		lastErr = err
		
		// Don't retry on the last attempt
		if attempt == config.MaxAttempts {
			break
		}
		
		// Check if error is retryable (network errors, temporary failures)
		if !isRetryableError(err) {
			logger.V(1).Info("non-retryable error encountered", "error", err, "attempt", attempt)
			return result1, result2, err
		}
		
		// Calculate delay with exponential backoff
		delay := time.Duration(attempt) * config.BaseDelay
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
		
		logger.V(1).Info("retrying operation after error", "error", err, "attempt", attempt, "delay", delay)
		
		select {
		case <-ctx.Done():
			return result1, result2, ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}
	
	return result1, result2, lastErr
}

// isRetryableError determines if an error should trigger a retry
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	
	// Network-related errors that are typically transient
	retryablePatterns := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"temporary failure",
		"network is unreachable",
		"no route to host",
		"i/o timeout",
		"context deadline exceeded",
		"TLS handshake timeout",
		"TLS handshake error",
		"server gave HTTP response to HTTPS client",
		"EOF",
	}
	
	for _, pattern := range retryablePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	
	return false
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      indexOf(s, substr) >= 0)))
}

// indexOf returns the index of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}