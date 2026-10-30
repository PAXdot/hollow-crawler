package crawler

import (
	"sync"
	"time"
)

// ...

// GetBestProxyOrRetry returns the proxy with the highest score that has not been
// used recently. If no such proxy is found, it retries with an exponential backoff.
func (ps *ProxyScorer) GetBestProxyOrRetry(maxIdleTime time.Duration, maxRetries int, backoffBase time.Duration) (string, bool, error) {
	ps.proxyMutex.Lock()
	defer ps.proxyMutex.Unlock()

	var bestProxy string
	var bestScore float64
	var retries int

	for retries < maxRetries {
		for proxyID, score := range ps.proxies {
			if time.Since(score.LastUsed) <= maxIdleTime && score.Score > bestScore {
				bestProxy = proxyID
				bestScore = score.Score
			}
		}

		if bestScore > 0 {
			return bestProxy, bestScore > 0, nil
		}

		time.Sleep(backoffBase.Microseconds() * time.Microsecond)
		backoffBase *= 2
		retries++
	}

	return "", false, ErrNoValidProxy
}

// retryWithBackoff retries a function with an exponential backoff until it succeeds
// or reaches the maximum number of retries.
func (ps *ProxyScorer) retryWithBackoff(maxRetries int, backoffBase time.Duration, fn func() (string, bool, error)) (string, bool, error) {
	for retries := 0; retries < maxRetries; retries++ {
		result, isSuccess, err := fn()
		if isSuccess || err == nil {
			return result, isSuccess, err
		}

		time.Sleep(backoffBase.Microseconds() * time.Microsecond)
		backoffBase *= 2
	}

	return "", false, ErrMaxRetriesExceeded
}

// ErrNoValidProxy is returned when no valid proxy is found after retrying.
var ErrNoValidProxy = &Error{Code: "NO_VALID_PROXY", Message: "No valid proxy found after retrying"}

// ErrMaxRetriesExceeded is returned when the maximum number of retries is exceeded.
var ErrMaxRetriesExceeded = &Error{Code: "MAX_RETRIES_EXCEEDED", Message: "Maximum number of retries exceeded"}

// Error represents an error returned by the crawler.
type Error struct {
	Code    string
	Message string
}
