package crawler

import (
	"fmt"
	"math"
	"time"
)

// RetryBackoffConfig defines the configuration for the retry/backoff strategy.
type RetryBackoffConfig struct {
	MaxRetries      int
	MinBackoffTime  time.Duration
	MaxBackoffTime  time.Duration
	BackoffFactor   float64
}

// RetryBackoffError is the error type for retry/backoff related issues.
type RetryBackoffError struct {
	Err     error
	Attempt int
}

func (e *RetryBackoffError) Error() string {
	return fmt.Sprintf("retry/backoff failed after %d attempts: %v", e.Attempt, e.Err)
}

// RetryBackoff retries the provided function with exponential backoff and up to maxRetries.
func RetryBackoff[R any, E error](f func() (R, E), config RetryBackoffConfig) (R, *RetryBackoffError) {
	var result R
	var err E
	for attempt := 1; attempt <= config.MaxRetries; attempt++ {
		result, err = f()
		if err == nil {
			return result, nil
		}
		backoffTime := calculateBackoffTime(config, attempt)
		time.Sleep(backoffTime)
	}
	return result, &RetryBackoffError{Err: err, Attempt: config.MaxRetries}
}

func calculateBackoffTime(config RetryBackoffConfig, attempt int) time.Duration {
	return time.Duration(float64(config.MinBackoffTime) * math.Pow(float64(config.BackoffFactor), float64(attempt-1)))
}
