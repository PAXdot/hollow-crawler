package crawler

import (
	"sync"
	"time"
)

// ThroughputCalculator tracks request rates and cumulative throughput
// for a crawling session. It is safe for concurrent use.
type ThroughputCalculator struct {
	mu          sync.Mutex
	start       time.Time
	lastReset   time.Time
	requests    int64
	bytes       int64
	// The sliding window tracks rates in a rolling timeframe
	// (since last Reset or since creation) which is useful
	// for observing performance trends and detecting spikes.
	windowStart time.Time
	windowReqs  int64
	windowBytes int64
}

// NewThroughputCalculator creates a calculator with the current time as epoch.
func NewThroughputCalculator() *ThroughputCalculator {
	now := time.Now()
	return &ThroughputCalculator{
		start:       now,
		lastReset:   now,
		windowStart: now,
	}
}

// RecordRequest increments the request counter, optionally with byte size.
func (tc *ThroughputCalculator) RecordRequest(size int64) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.requests++
	tc.bytes += size
	tc.windowReqs++
	tc.windowBytes += size
}

// Reset clears all counters and restarts the measurement window.
func (tc *ThroughputCalculator) Reset() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	now := time.Now()
	tc.start = now
	tc.lastReset = now
	tc.windowStart = now
	tc.requests = 0
	tc.bytes = 0
	tc.windowReqs = 0
	tc.windowBytes = 0
}

// ... (rest of the file remains the same)
