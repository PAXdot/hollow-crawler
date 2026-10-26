package crawler

import (
	"sync"
	"time"
)

// RobotsState tracks per-host politeness state for robots.txt compliance.
type RobotsState struct {
	mu            sync.Mutex
	lastFetch     map[string]time.Time
	crawlDelay    map[string]time.Duration
}

// NewRobotsState creates a new RobotsState.
func NewRobotsState() *RobotsState {
	return &RobotsState{
		lastFetch:  make(map[string]time.Time),
		crawlDelay: make(map[string]time.Duration),
	}
}

// SetCrawlDelay sets the crawl delay for a given host.
func (rs *RobotsState) SetCrawlDelay(host string, delay time.Duration) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.crawlDelay[host] = delay
}

// CrawlDelay returns the crawl delay for a host, defaulting to zero.
func (rs *RobotsState) CrawlDelay(host string) time.Duration {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.crawlDelay[host]
}

// RecordFetch records that a fetch happened for host at time t.
func (rs *RobotsState) RecordFetch(host string, t time.Time) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.lastFetch[host] = t
}

// LastFetch returns the last fetch time for host, or the zero time.
func (rs *RobotsState) LastFetch(host string) time.Time {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	return rs.lastFetch[host]
}

// WaitDuration calculates how long to wait before the next allowed fetch for host.
// Returns zero if immediate fetch is allowed.
func (rs *RobotsState) WaitDuration(host string) time.Duration {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	delay, ok := rs.crawlDelay[host]
	if !ok {
		return 0
	}
	last, ok2 := rs.lastFetch[host]
	if !ok2 {
		return 0
	}
	elapsed := time.Since(last)
	if elapsed >= delay {
		return 0
	}
	return delay - elapsed
}
