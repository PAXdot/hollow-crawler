package crawler

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Stats represents the statistics about the crawler.
type Stats struct {
	// Number of requests made.
	RequestsMade uint64

	// Number of successes, i.e. status code 200.
	Successes uint64

	// Number of failures, i.e. status code != 200.
	Failures uint64

	// Number of retries made.
	Retries uint64

	// Total duration of the crawler in seconds.
	TotalDuration time.Duration

	// Duration of the last request in seconds.
	LastRequestDuration time.Duration

	// Number of unique URLs found.
	UniqueURLs int

	// Number of pages parsed.
	PagesParsed int

	// Mutex to protect the statistics.
	m sync.Mutex
}

// IncrementRequestsMade increments the number of requests made.
func (s *Stats) IncrementRequestsMade() {
	atomic.AddUint64(&s.RequestsMade, 1)
}

// IncrementSuccesses increments the number of successes.
func (s *Stats) IncrementSuccesses() {
	atomic.AddUint64(&s.Successes, 1)
}

// IncrementFailures increments the number of failures.
func (s *Stats) IncrementFailures() {
	atomic.AddUint64(&s.Failures, 1)
}

// IncrementRetries increments the number of retries made.
func (s *Stats) IncrementRetries() {
	atomic.AddUint64(&s.Retries, 1)
}

// UpdateLastRequestDuration updates the duration of the last request.
func (s *Stats) UpdateLastRequestDuration(d time.Duration) {
	s.LastRequestDuration = d
}

// IncrementUniqueURLs increments the number of unique URLs found.
func (s *Stats) IncrementUniqueURLs() {
	s.UniqueURLs++
}

// IncrementPagesParsed increments the number of pages parsed.
func (s *Stats) IncrementPagesParsed() {
	s.PagesParsed++
}

// GetTotalDuration returns the total duration of the crawler in seconds.
func (s *Stats) GetTotalDuration() time.Duration {
	s.m.Lock()
	defer s.m.Unlock()
	return s.TotalDuration
}

// SetTotalDuration sets the total duration of the crawler in seconds.
func (s *Stats) SetTotalDuration(d time.Duration) {
	s.m.Lock()
	defer s.m.Unlock()
	s.TotalDuration = d
}

// Print prints the current statistics.
func (s *Stats) Print() {
	s.m.Lock()
	defer s.m.Unlock()
	fmt.Printf("Requests Made: %d\n", s.RequestsMade)
	fmt.Printf("Successes: %d\n", s.Successes)
	fmt.Printf("Failures: %d\n", s.Failures)
	fmt.Printf("Retries: %d\n", s.Retries)
	fmt.Printf("Total Duration: %v\n", s.TotalDuration)
	fmt.Printf("Last Request Duration: %v\n", s.LastRequestDuration)
	fmt.Printf("Unique URLs: %d\n", s.UniqueURLs)
	fmt.Printf("Pages Parsed: %d\n", s.PagesParsed)
}
