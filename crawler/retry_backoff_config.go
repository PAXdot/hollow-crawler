// crawler/retry_backoff_config.go

package crawler

import "time"

// RetryBackoffConfig represents the retry and backoff configuration.
type RetryBackoffConfig struct {
	MaxRetries      int
	MinBackoff      time.Duration
	MaxBackoff      time.Duration
	BackoffFactor   float64
	BackoffJitter   float64
	NonRetryStatusCodes []int
}

// DefaultRetryBackoffConfig returns the default retry and backoff configuration.
func DefaultRetryBackoffConfig() RetryBackoffConfig {
	return RetryBackoffConfig{
		MaxRetries:      5,
		MinBackoff:      1 * time.Second,
		MaxBackoff:      60 * time.Second,
		BackoffFactor:   2,
		BackoffJitter:   0.2,
		NonRetryStatusCodes: []int{404, 410, 500, 503},
	}
}
