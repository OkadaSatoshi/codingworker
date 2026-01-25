package retry

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Policy defines retry behavior
type Policy struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Multiplier     float64
}

// DefaultPolicy returns the default exponential backoff policy
// 30s -> 60s -> 120s (max 3 retries)
func DefaultPolicy() *Policy {
	return &Policy{
		MaxRetries:     3,
		InitialBackoff: 30 * time.Second,
		MaxBackoff:     120 * time.Second,
		Multiplier:     2.0,
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
