package crawler

import (
	"errors"
	"fmt"
	"time"
)

// RetryError wraps an underlying error with metadata about retry attempts
// and backoff schedule.
type RetryError struct {
	Err         error
	Attempt     int
	MaxAttempts int
	NextBackoff time.Duration
}

func (e *RetryError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("retry error (attempt %d/%d, next backoff: %v)",
			e.Attempt, e.MaxAttempts, e.NextBackoff)
	}
	return fmt.Sprintf("retry error (attempt %d/%d, next backoff: %v): %v",
		e.Attempt, e.MaxAttempts, e.NextBackoff, e.Err)
}

func (e *RetryError) Unwrap() error {
	return e.Err
}

// IsRetryExhausted returns true if the maximum number of retry attempts
// has been reached.
func IsRetryExhausted(err error) bool {
	var re *RetryError
	if errors.As(err, &re) {
		return re.Attempt >= re.MaxAttempts
	}
	return false
}

// RetryAfter returns the suggested backoff duration from a RetryError, or
// zero if the error is not a RetryError.
func RetryAfter(err error) time.Duration {
	var re *RetryError
	if errors.As(err, &re) {
		return re.NextBackoff
	}
	return 0
}

// RetryAttempt returns the current attempt number from a RetryError, or 0.
func RetryAttempt(err error) int {
	var re *RetryError
	if errors.As(err, &re) {
		return re.Attempt
	}
	return 0
}

// NewRetryError creates a RetryError with the given underlying error, the
// current attempt (1-based), max attempts, and the next backoff duration.
func NewRetryError(err error, attempt, maxAttempts int, nextBackoff time.Duration) *RetryError {
	return &RetryError{
		Err:         err,
		Attempt:     attempt,
		MaxAttempts: maxAttempts,
		NextBackoff: nextBackoff,
	}
}
