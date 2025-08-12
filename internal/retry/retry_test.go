package retry

import (
	"context"
	"errors"
	"testing"
)

func TestWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	result, err := WithRetry(ctx, func() (string, error) {
		return "success", nil
	})
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}
}

func TestWithRetry_RetryableError(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	
	result, err := WithRetry(ctx, func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("connection refused")
		}
		return "success", nil
	})
	
	if err != nil {
		t.Errorf("Expected no error after retries, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %v", result)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestWithRetry_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	
	result, err := WithRetry(ctx, func() (string, error) {
		attempts++
		return "", errors.New("invalid input")
	})
	
	if err == nil {
		t.Error("Expected error for non-retryable error")
	}
	if result != "" {
		t.Errorf("Expected empty result, got %v", result)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt for non-retryable error, got %d", attempts)
	}
}

func TestWithRetry_MaxAttemptsReached(t *testing.T) {
	ctx := context.Background()
	attempts := 0
	
	result, err := WithRetry(ctx, func() (string, error) {
		attempts++
		return "", errors.New("timeout")
	})
	
	if err == nil {
		t.Error("Expected error after max attempts")
	}
	if result != "" {
		t.Errorf("Expected empty result, got %v", result)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestIsRetryableError(t *testing.T) {
	testCases := []struct {
		err        error
		retryable  bool
		name       string
	}{
		{errors.New("connection refused"), true, "connection refused"},
		{errors.New("TLS handshake timeout"), true, "TLS handshake timeout"},
		{errors.New("server gave HTTP response to HTTPS client"), true, "HTTP to HTTPS error"},
		{errors.New("invalid input"), false, "invalid input"},
		{errors.New("unauthorized"), false, "unauthorized"},
		{nil, false, "nil error"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isRetryableError(tc.err)
			if result != tc.retryable {
				t.Errorf("Expected %v for error '%v', got %v", tc.retryable, tc.err, result)
			}
		})
	}
}