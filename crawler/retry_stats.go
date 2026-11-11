package crawler

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// NewRetryStats returns an initialized RetryStats with the reset time set to now.
func NewRetryStats() *RetryStats {
	return &RetryStats{lastReset: time.Now()}
}

// ExponentialBackoff calculates the next backoff duration based on the number of attempts.
func (s *RetryStats) ExponentialBackoff(attempt uint64) time.Duration {
	// Maximum backoff duration is 60 seconds.
	maxBackoff := time.Minute
	base := time.Second
	return min(maxBackoff, base*uint64(math.Pow(2, float64(attempt))))
}

// RecordAttempt increments the attempt counter and records the time, then calculates the backoff.
func (s *RetryStats) RecordAttempt(attempt uint64) {
	s.attempts.Add(1)
	s.lastAttempt.Store(time.Now().UnixNano())
	backoff := s.ExponentialBackoff(attempt)
	s.RecordBackoff(backoff)
	time.Sleep(backoff)
}

// ... Rest of the file remains unchanged
