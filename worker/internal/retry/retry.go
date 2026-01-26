package retry

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Policy defines retry behavior
type Policy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
}

// DefaultPolicy returns the default retry policy with fixed 10s backoff
func DefaultPolicy() *Policy {
	return NewPolicy(3)
}

// NewPolicy creates a retry policy with the specified max retries and fixed 10s backoff
func NewPolicy(maxRetries int) *Policy {
	return &Policy{
		MaxRetries:     maxRetries,
		InitialBackoff: 10 * time.Second,
		MaxBackoff:     10 * time.Second, // Fixed backoff (no exponential growth)
		Multiplier:     1.0,              // No multiplier
	}
}

// ErrorType classifies errors for retry decisions
type ErrorType int

const (
	// ErrorTypeTransient - temporary errors that should be retried
	ErrorTypeTransient ErrorType = iota
	// ErrorTypePermanent - permanent errors that should not be retried
	ErrorTypePermanent
	// ErrorTypeTimeout - timeout errors that should not be retried
	ErrorTypeTimeout
)

// ClassifiableError is an error that can indicate if it's retryable
type ClassifiableError interface {
	error
	ErrorType() ErrorType
}

// TransientError represents a temporary error
type TransientError struct {
	Err error
}

func (e *TransientError) Error() string {
	return fmt.Sprintf("transient error: %v", e.Err)
}

func (e *TransientError) ErrorType() ErrorType {
	return ErrorTypeTransient
}

func (e *TransientError) Unwrap() error {
	return e.Err
}

// PermanentError represents a permanent error
type PermanentError struct {
	Err error
}

func (e *PermanentError) Error() string {
	return fmt.Sprintf("permanent error: %v", e.Err)
}

func (e *PermanentError) ErrorType() ErrorType {
	return ErrorTypePermanent
}

func (e *PermanentError) Unwrap() error {
	return e.Err
}

// Result holds the outcome of a retryable operation
type Result struct {
	Attempts int
	LastErr  error
}

// Do executes the function with retry logic
func (p *Policy) Do(ctx context.Context, operation func() error) *Result {
	result := &Result{}
	backoff := p.InitialBackoff

	for attempt := 1; attempt <= p.MaxRetries+1; attempt++ {
		result.Attempts = attempt

		err := operation()
		if err == nil {
			result.LastErr = nil // Clear any previous error on success
			return result
		}

		result.LastErr = err

		// Check if error is permanent
		if classifiable, ok := err.(ClassifiableError); ok {
			if classifiable.ErrorType() != ErrorTypeTransient {
				slog.Warn("Permanent error, not retrying",
					"attempt", attempt,
					"error", err,
				)
				return result
			}
		}

		// Don't retry if we've exhausted attempts
		if attempt > p.MaxRetries {
			slog.Error("Max retries exceeded",
				"attempts", attempt,
				"error", err,
			)
			return result
		}

		slog.Warn("Operation failed, retrying",
			"attempt", attempt,
			"backoff", backoff,
			"error", err,
		)

		// Wait with backoff
		select {
		case <-ctx.Done():
			result.LastErr = ctx.Err()
			return result
		case <-time.After(backoff):
		}

		// Increase backoff for next attempt
		backoff = time.Duration(float64(backoff) * p.Multiplier)
		if backoff > p.MaxBackoff {
			backoff = p.MaxBackoff
		}
	}

	return result
}

// IsRetryExhausted returns true if all retries were used
func (r *Result) IsRetryExhausted() bool {
	return r.Attempts > DefaultPolicy().MaxRetries && r.LastErr != nil
}

// ClassifyHTTPStatus returns the error type based on HTTP status code
func ClassifyHTTPStatus(statusCode int) ErrorType {
	switch {
	case statusCode == 429:
		// Rate limit - retry after backoff
		return ErrorTypeTransient
	case statusCode >= 500 && statusCode < 600:
		// Server error - retry
		return ErrorTypeTransient
	case statusCode >= 400 && statusCode < 500:
		// Client error (except 429) - don't retry
		return ErrorTypePermanent
	default:
		// Success or unknown - don't retry
		return ErrorTypePermanent
	}
}

// ClassifyGitError classifies git/gh CLI output for retry decisions
func ClassifyGitError(output string) ErrorType {
	// Transient errors (network, rate limit, server issues)
	transientPatterns := []string{
		"Could not resolve host",
		"Connection refused",
		"Connection timed out",
		"rate limit",
		"API rate limit",
		"503",
		"502",
		"500",
		"temporarily unavailable",
		"try again later",
	}

	for _, pattern := range transientPatterns {
		if strings.Contains(strings.ToLower(output), strings.ToLower(pattern)) {
			return ErrorTypeTransient
		}
	}

	// Default to permanent for client errors
	return ErrorTypePermanent
}

// WrapWithClassification wraps an error with appropriate retry classification
func WrapWithClassification(err error, output string) error {
	if err == nil {
		return nil
	}

	errType := ClassifyGitError(output)
	if errType == ErrorTypeTransient {
		return &TransientError{Err: err}
	}
	return &PermanentError{Err: err}
}
