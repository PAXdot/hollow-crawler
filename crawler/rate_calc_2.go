package crawler

import (
	"time"
)

// RateCalc calculates HTML extraction rates and implements retry/backoff.
type RateCalc struct {
	extractStarts int
	extractEnds   int
	startTime     time.Time
	retryCount    int       // New field to track retries
	retryLimit    int       // Max retries before giving up
	backoffDelay  time.Duration // Initial delay between retries
}

// NewRateCalc returns a new RateCalc instance with default retry limit and backoff delay.
func NewRateCalc() *RateCalc {
	return &RateCalc{
		startTime:    time.Now(),
		retryLimit:   3, // Default retry limit
		backoffDelay: 5 * time.Second, // Default backoff delay
	}
}

// OnExtractStart increments extractStarts and updates the start time.
func (r *RateCalc) OnExtractStart() {
	r.extractStarts++
	r.startTime = time.Now()
}

// OnExtractEnd increments extractEnds and resets retry count.
func (r *RateCalc) OnExtractEnd() {
	r.extractEnds++
	r.retryCount = 0 // Reset retry count
}

// RetryBackoff waits for the calculated backoff duration before retrying.
func (r *RateCalc) RetryBackoff() {
	backoff := r.backoffDelay * time.Duration(1<<uint(r.retryCount))
	time.Sleep(backoff)
	r.backoffDelay *= 2 // Double the backoff delay for next retry
	r.retryCount++       // Increment retry count
	if r.retryCount > r.retryLimit {
		log.Warningf("Gave up after %d retries for extraction", r.retryLimit)
		r.backoffDelay = 5 * time.Second // Reset backoff delay if max retries reached
	}
}

// Rate returns the current extraction rate (extracted per second).
func (r *RateCalc) Rate() float64 {
	elapsed := time.Since(r.startTime).Seconds()
	return float64(r.extractEnds) / elapsed
}

// TotalExtracted returns the total number of extracted HTMLs.
func (r *RateCalc) TotalExtracted() int {
	return r.extractEnds
}
