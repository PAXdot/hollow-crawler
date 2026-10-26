package crawler

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter calculates and enforces a moving average rate limit.
type RateLimiter struct {
	mutex          sync.Mutex
	lastSuccess    time.Time
	retryAttempts  int
	rateLimit       int // requests per second
	rateLimitPeriod time.Duration
}

// NewRateLimiter creates a new RateLimiter with the given rate limit per second.
// The rate limit period is the time over which requests are averaged.
func NewRateLimiter(rateLimit int, rateLimitPeriod time.Duration) *RateLimiter {
	return &RateLimiter{
		rateLimit:       rateLimit,
		rateLimitPeriod: rateLimitPeriod,
	}
}

// AllowRequest checks if a request can be made without exceeding the rate limit.
func (rl *RateLimiter) AllowRequest() (bool, string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	elapsedTime := time.Since(rl.lastSuccess)
	if elapsedTime < rl.rateLimitPeriod {
		retryPeriod := (rl.rateLimitPeriod - elapsedTime) / time.Second
		time.Sleep(retryPeriod)
		rl.retryAttempts++
		return false, fmt.Sprintf("Rate limit exceeded. Retrying in %v seconds (%v attempts)", retryPeriod.Seconds(), rl.retryAttempts)
	}

	rl.lastSuccess = time.Now()
	return true, ""
}

// ResetRetryCount resets the retry attempt count.
func (rl *RateLimiter) ResetRetryCount() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.retryAttempts = 0
}
