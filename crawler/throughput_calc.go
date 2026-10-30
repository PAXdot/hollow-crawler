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

// RequestsPerSecond returns the average request rate since start.
func (tc *ThroughputCalculator) RequestsPerSecond() float64 {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	elapsed := time.Since(tc.start).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return float64(tc.requests) / elapsed
}

// BytesPerSecond returns the average byte throughput since start.
func (tc *ThroughputCalculator) BytesPerSecond() float64 {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	elapsed := time.Since(tc.start).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return float64(tc.bytes) / elapsed
}

// WindowRequestsPerSecond returns the request rate in the current sliding window
// (since last Reset or since creation).
func (tc *ThroughputCalculator) WindowRequestsPerSecond() float64 {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	elapsed := time.Since(tc.windowStart).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return float64(tc.windowReqs) / elapsed
}

// WindowBytesPerSecond returns the byte rate in the current sliding window.
func (tc *ThroughputCalculator) WindowBytesPerSecond() float64 {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	elapsed := time.Since(tc.windowStart).Seconds()
	if elapsed <= 0 {
		return 0
	}
	return float64(tc.windowBytes) / elapsed
}

// Totals returns cumulative request and byte counts.
func (tc *ThroughputCalculator) Totals() (requests int64, bytes int64) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return tc.requests, tc.bytes
}

// Elapsed returns the duration since the calculator was created or reset.
func (tc *ThroughputCalculator) Elapsed() time.Duration {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return time.Since(tc.start)
}
