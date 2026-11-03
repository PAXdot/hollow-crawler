package crawler

import (
	"math"
	"time"
)

// RetryBackoff represents a retry/backoff strategy with exponential backoff.
type RetryBackoff struct {
	maxRetries int
	maxDuration time.Duration
 CurrentRetry int
	LastSuccess time.Time
}

// NewRetryBackoff creates a new RetryBackoff instance with the given maxRetries and maxDuration.
func NewRetryBackoff(maxRetries int, maxDuration time.Duration) *RetryBackoff {
	return &RetryBackoff{
		maxRetries: maxRetries,
		maxDuration: maxDuration,
		CurrentRetry: 0,
		LastSuccess: time.Time{},
	}
}

// ShouldRetry checks if a retry is necessary based on the current retry count and elapsed time since last success.
func (r *RetryBackoff) ShouldRetry(elapsedTime time.Duration) bool {
	if r.CurrentRetry >= r.maxRetries {
		return false
	}
	if r.LastSuccess.IsZero() || elapsedTime > r.maxDuration {
		return true
	}
	return false
}

// NextRetryDelay calculates the delay for the next retry using exponential backoff.
func (r *RetryBackoff) NextRetryDelay() time.Duration {
	if r.CurrentRetry == 0 {
		return 0
	}
	baseDelay := time.Duration(math.Pow(float64(r.CurrentRetry), 2))
	return time.Duration(baseDelay) * time.Second
}

// OnSuccess updates the last success time and resets the current retry count.
func (r *RetryBackoff) OnSuccess() {
	r.LastSuccess = time.Now()
	r.CurrentRetry = 0
}

// OnRetry increments the current retry count.
func (r *RetryBackoff) OnRetry() {
	r.CurrentRetry++
}
