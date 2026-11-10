package crawler

import (
	"sync/atomic"
	"time"
)

// RetryBackoffConfig represents the configuration for retry/backoff.
type RetryBackoffConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	BackoffFactor  float64
}

// RetryBackoffState holds the state for retry/backoff.
type RetryBackoffState struct {
	retries       uint32
	backoff       time.Duration
	config        *RetryBackoffConfig
	shouldRetry   bool
	attemptsLeft atomic.Value // Use atomic.Value for safe access/modification
}

// NewRetryBackoffState creates a new RetryBackoffState with the given config.
func NewRetryBackoffState(config *RetryBackoffConfig) *RetryBackoffState {
	return &RetryBackoffState{
		config:        config,
		backoff:       config.InitialBackoff,
		shouldRetry:   true,
		attemptsLeft:  atomic.Value{},
	}
}

// Retry should be called when a task needs to be retried.
func (s *RetryBackoffState) Retry() {
	s.retries++
	s.backoff = calculateBackoff(s.backoff, s.config.BackoffFactor)

	// Update attemptsleft using atomic.Value
	s.attemptsLeft.Store(s.config.MaxRetries - int(s.retries))

	s.shouldRetry = s.attemptsLeft.Load().(int) > 0
}

// AttemptsLeft returns the number of attempts left for retries.
func (s *RetryBackoffState) AttemptsLeft() int {
	return s.attemptsLeft.Load().(int)
}

// WaitingTime returns the time to wait before the next retry.
func (s *RetryBackoffState) WaitingTime() time.Duration {
	if !s.shouldRetry {
		return 0
	}
	return s.backoff
}

// calculateBackoff calculates the next backoff delay.
func calculateBackoff(current, factor float64) time.Duration {
	return time.Duration(min(current*factor, float64(s.config.MaxBackoff)))
}

// min returns the minimum of two values.
func min(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}
