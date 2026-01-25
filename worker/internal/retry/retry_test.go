package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestClassifyHTTPStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       ErrorType
	}{
		{"429 rate limit", 429, ErrorTypeTransient},
		{"500 server error", 500, ErrorTypeTransient},
		{"502 bad gateway", 502, ErrorTypeTransient},
		{"503 service unavailable", 503, ErrorTypeTransient},
		{"400 bad request", 400, ErrorTypePermanent},
		{"401 unauthorized", 401, ErrorTypePermanent},
		{"403 forbidden", 403, ErrorTypePermanent},
		{"404 not found", 404, ErrorTypePermanent},
		{"200 ok", 200, ErrorTypePermanent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyHTTPStatus(tt.statusCode)
			if got != tt.want {
				t.Errorf("ClassifyHTTPStatus(%d) = %v, want %v", tt.statusCode, got, tt.want)
			}
		})
	}
}

func TestClassifyGitError(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   ErrorType
	}{
		{"network error", "Could not resolve host: github.com", ErrorTypeTransient},
		{"connection refused", "Connection refused", ErrorTypeTransient},
		{"connection timeout", "Connection timed out", ErrorTypeTransient},
		{"rate limit", "API rate limit exceeded", ErrorTypeTransient},
		{"server 503", "error 503: Service Unavailable", ErrorTypeTransient},
		{"temporarily unavailable", "Repository temporarily unavailable", ErrorTypeTransient},
		{"try again later", "Please try again later", ErrorTypeTransient},
		{"not found", "Repository not found", ErrorTypePermanent},
		{"permission denied", "Permission denied (publickey)", ErrorTypePermanent},
		{"unknown error", "Some unknown error", ErrorTypePermanent},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyGitError(tt.output)
			if got != tt.want {
				t.Errorf("ClassifyGitError(%q) = %v, want %v", tt.output, got, tt.want)
			}
		})
	}
}

func TestWrapWithClassification(t *testing.T) {
	baseErr := errors.New("base error")

	t.Run("nil error returns nil", func(t *testing.T) {
		got := WrapWithClassification(nil, "any output")
		if got != nil {
			t.Errorf("WrapWithClassification(nil, ...) = %v, want nil", got)
		}
	})

	t.Run("transient error wrapping", func(t *testing.T) {
		got := WrapWithClassification(baseErr, "Could not resolve host")
		if _, ok := got.(*TransientError); !ok {
			t.Errorf("expected TransientError, got %T", got)
		}
	})

	t.Run("permanent error wrapping", func(t *testing.T) {
		got := WrapWithClassification(baseErr, "Repository not found")
		if _, ok := got.(*PermanentError); !ok {
			t.Errorf("expected PermanentError, got %T", got)
		}
	})
}

func TestPolicy_Do_Success(t *testing.T) {
	policy := &Policy{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	callCount := 0
	result := policy.Do(context.Background(), func() error {
		callCount++
		return nil
	})

	if result.LastErr != nil {
		t.Errorf("expected no error, got %v", result.LastErr)
	}
	if result.Attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", result.Attempts)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestPolicy_Do_TransientRetry(t *testing.T) {
	policy := &Policy{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	callCount := 0
	result := policy.Do(context.Background(), func() error {
		callCount++
		if callCount < 3 {
			return &TransientError{Err: errors.New("temporary failure")}
		}
		return nil
	})

	if result.LastErr != nil {
		t.Errorf("expected no error after retries, got %v", result.LastErr)
	}
	if result.Attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", result.Attempts)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestPolicy_Do_PermanentNoRetry(t *testing.T) {
	policy := &Policy{
		MaxRetries:     3,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	callCount := 0
	result := policy.Do(context.Background(), func() error {
		callCount++
		return &PermanentError{Err: errors.New("permanent failure")}
	})

	if result.LastErr == nil {
		t.Error("expected error, got nil")
	}
	if result.Attempts != 1 {
		t.Errorf("expected 1 attempt (no retry for permanent), got %d", result.Attempts)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestPolicy_Do_MaxRetriesExceeded(t *testing.T) {
	policy := &Policy{
		MaxRetries:     2,
		InitialBackoff: 10 * time.Millisecond,
		MaxBackoff:     100 * time.Millisecond,
		Multiplier:     2.0,
	}

	callCount := 0
	result := policy.Do(context.Background(), func() error {
		callCount++
		return &TransientError{Err: errors.New("always fails")}
	})

	if result.LastErr == nil {
		t.Error("expected error after max retries, got nil")
	}
	// MaxRetries=2 means 3 total attempts (initial + 2 retries)
	if result.Attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", result.Attempts)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestPolicy_Do_ContextCancelled(t *testing.T) {
	policy := &Policy{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Second,
		MaxBackoff:     10 * time.Second,
		Multiplier:     2.0,
	}

	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	result := policy.Do(ctx, func() error {
		callCount++
		return &TransientError{Err: errors.New("transient")}
	})

	if !errors.Is(result.LastErr, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", result.LastErr)
	}
}

func TestTransientError(t *testing.T) {
	baseErr := errors.New("base")
	err := &TransientError{Err: baseErr}

	if err.ErrorType() != ErrorTypeTransient {
		t.Errorf("expected ErrorTypeTransient, got %v", err.ErrorType())
	}

	if !errors.Is(err, baseErr) {
		t.Error("Unwrap should return base error")
	}

	if err.Error() != "transient error: base" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestPermanentError(t *testing.T) {
	baseErr := errors.New("base")
	err := &PermanentError{Err: baseErr}

	if err.ErrorType() != ErrorTypePermanent {
		t.Errorf("expected ErrorTypePermanent, got %v", err.ErrorType())
	}

	if !errors.Is(err, baseErr) {
		t.Error("Unwrap should return base error")
	}

	if err.Error() != "permanent error: base" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
