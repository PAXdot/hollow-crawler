package crawler

import (
	"sync"
	"sync/atomic"
	"time"
)

// RetryStats tracks cumulative statistics for retry and backoff operations.
// All exported methods are safe for concurrent use.
type RetryStats struct {
	attempts       atomic.Uint64
	successes      atomic.Uint64
	failures       atomic.Uint64
	totalBackoff   atomic.Int64 // nanoseconds
	mu             sync.Mutex
	lastReset      time.Time
	lastAttempt    atomic.Int64 // unix nano
}

// NewRetryStats returns an initialized RetryStats with the reset time set to now.
func NewRetryStats() *RetryStats {
	return &RetryStats{lastReset: time.Now()}
}

// RecordAttempt increments the attempt counter and records the time.
func (s *RetryStats) RecordAttempt() {
	s.attempts.Add(1)
	s.lastAttempt.Store(time.Now().UnixNano())
}

// RecordSuccess increments the success counter.
func (s *RetryStats) RecordSuccess() {
	s.successes.Add(1)
}

// RecordFailure increments the failure counter.
func (s *RetryStats) RecordFailure() {
	s.failures.Add(1)
}

// RecordBackoff adds the given duration to the total accumulated backoff.
func (s *RetryStats) RecordBackoff(d time.Duration) {
	s.totalBackoff.Add(int64(d))
}

// Attempts returns the total number of recorded attempts.
func (s *RetryStats) Attempts() uint64 {
	return s.attempts.Load()
}

// Successes returns the total number of recorded successes.
func (s *RetryStats) Successes() uint64 {
	return s.successes.Load()
}

// Failures returns the total number of recorded failures.
func (s *RetryStats) Failures() uint64 {
	return s.failures.Load()
}

// TotalBackoff returns the total accumulated backoff duration.
func (s *RetryStats) TotalBackoff() time.Duration {
	return time.Duration(s.totalBackoff.Load())
}

// SuccessRate returns the fraction of successful attempts relative to total
// attempts. If no attempts have been made, it returns 0.
func (s *RetryStats) SuccessRate() float64 {
	a := s.attempts.Load()
	if a == 0 {
		return 0
	}
	return float64(s.successes.Load()) / float64(a)
}

// FailureRate returns the fraction of failed attempts relative to total
// attempts. If no attempts have been made, it returns 0.
func (s *RetryStats) FailureRate() float64 {
	a := s.attempts.Load()
	if a == 0 {
		return 0
	}
	return float64(s.failures.Load()) / float64(a)
}

// AverageBackoff returns the mean backoff duration per recorded attempt.
// Returns 0 if no attempts have been made.
func (s *RetryStats) AverageBackoff() time.Duration {
	a := s.attempts.Load()
	if a == 0 {
		return 0
	}
	return time.Duration(s.totalBackoff.Load() / int64(a))
}

// LastAttempt returns the time of the last recorded attempt, or the zero time
// if none has occurred.
func (s *RetryStats) LastAttempt() time.Time {
	nano := s.lastAttempt.Load()
	if nano == 0 {
		return time.Time{}
	}
	return time.Unix(0, nano)
}

// Reset clears all counters and sets the last reset time to now.
func (s *RetryStats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.attempts.Store(0)
	s.successes.Store(0)
	s.failures.Store(0)
	s.totalBackoff.Store(0)
	s.lastAttempt.Store(0)
	s.lastReset = time.Now()
}

// SinceReset returns the duration elapsed since the last reset or creation.
func (s *RetryStats) SinceReset() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return time.Since(s.lastReset)
}

// Snapshot returns a consistent point-in-time copy of all statistics.
func (s *RetryStats) Snapshot() RetryStatsSnapshot {
	return RetryStatsSnapshot{
		Attempts:     s.attempts.Load(),
		Successes:    s.successes.Load(),
		Failures:     s.failures.Load(),
		TotalBackoff: time.Duration(s.totalBackoff.Load()),
		LastAttempt:  s.LastAttempt(),
	}
}

// RetryStatsSnapshot is a frozen copy of RetryStats counters.
type RetryStatsSnapshot struct {
	Attempts     uint64
	Successes    uint64
	Failures     uint64
	TotalBackoff time.Duration
	LastAttempt  time.Time
}

// SuccessRate returns the success ratio from the snapshot.
func (s RetryStatsSnapshot) SuccessRate() float64 {
	if s.Attempts == 0 {
		return 0
	}
	return float64(s.Successes) / float64(s.Attempts)
}

// FailureRate returns the failure ratio from the snapshot.
func (s RetryStatsSnapshot) FailureRate() float64 {
	if s.Attempts == 0 {
		return 0
	}
	return float64(s.Failures) / float64(s.Attempts)
}

// AverageBackoff returns the mean backoff from the snapshot.
func (s RetryStatsSnapshot) AverageBackoff() time.Duration {
	if s.Attempts == 0 {
		return 0
	}
	return s.TotalBackoff / time.Duration(s.Attempts)
}
