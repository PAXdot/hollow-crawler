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
	retries     uint32
	backoff     time.Duration
	config      *RetryBackoffConfig
	shouldRetry bool
}

// NewRetryBackoffState creates a new RetryBackoffState with the given config.
func NewRetryBackoffState(config *RetryBackoffConfig) *RetryBackoffState {
	return &RetryBackoffState{
		config:      config,
		backoff:     config.InitialBackoff,
		shouldRetry: true,
	}
}

// Retry should be called when a task needs to be retried.
func (s *RetryBackoffState) Retry() {
	s.retries++
	s.backoff = calculateBackoff(s.backoff, s.config.BackoffFactor)
	s.shouldRetry = s.retries < uint32(s.config.MaxRetries)
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
	return time.Duration(current * factor)
}
