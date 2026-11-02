package crawler

import (
	"math"
	"time"
)

// RetryBackoffScore computes a rank for a retry attempt based on the
// current backoff duration, attempt number, and elapsed time. It returns
// a float64 in the range [0, 1], where higher values indicate a more
// favorable time to retry. The score balances waiting long enough to
// respect backoff with not stalling progress unnecessarily.
func RetryBackoffScore(attempt int, baseDelay, maxDelay time.Duration, elapsed time.Duration) float64 {
	if attempt < 0 || baseDelay <= 0 || maxDelay < baseDelay || elapsed < 0 {
		return 0
	}

	// Compute the exponential backoff for the current attempt.
	backoff := time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))
	if backoff > maxDelay {
		backoff = maxDelay
	}

	// If the minimum backoff hasn't elapsed, the score is 0.
	if elapsed < baseDelay {
		return 0
	}

	// If we've waited at least the full backoff, the score is 1.
	if elapsed >= backoff {
		return 1
	}

	// Otherwise, interpolate linearly between baseDelay and backoff.
	progress := float64(elapsed-baseDelay) / float64(backoff-baseDelay)
	return progress
}
