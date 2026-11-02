package crawler

import (
	"errors"
	"time"
)

// RobotsPolicyIterator yields allowed URLs with respect to robots.txt
// rules and a minimum delay between successive requests to the same host.
type RobotsPolicyIterator struct {
	urls       []string
	policy     *RobotsPolicy
	lastAccess map[string]time.Time // host -> last request time
	delay      time.Duration
	pos        int
	done       bool
}

// NewRobotsPolicyIterator creates an iterator over urls that respects
// the given policy and per-host politeness delay.
func NewRobotsPolicyIterator(urls []string, policy *RobotsPolicy, delay time.Duration) *RobotsPolicyIterator {
	return &RobotsPolicyIterator{
		urls:       urls,
		policy:     policy,
		lastAccess: make(map[string]time.Time),
		delay:      delay,
		pos:        0,
		done:       false,
	}
}

// Next returns the next allowed URL and true, or ("", false) when exhausted.
// It skips URLs that are disallowed by policy or would violate the delay.
func (it *RobotsPolicyIterator) Next() (string, bool) {
	for it.pos < len(it.urls) {
		u := it.urls[it.pos]
		it.pos++

		if !it.policy.Allowed(u) {
			continue
		}

		host := extractHost(u)
		if host == "" {
			continue
		}

		last, ok := it.lastAccess[host]
		now := time.Now()
		if ok && now.Sub(last) < it.delay {
			continue
		}
		it.lastAccess[host] = now
		return u, true
	}
	it.done = true
	return "", false
}

// Done returns true when the iterator has been fully consumed.
func (it *RobotsPolicyIterator) Done() bool {
	return it.done
}

// Reset reinitializes the iterator to the beginning.
func (it *RobotsPolicyIterator) Reset() {
	it.pos = 0
	it.lastAccess = make(map[string]time.Time)
	it.done = false
}

// Err returns any error encountered during iteration (currently always nil).
func (it *RobotsPolicyIterator) Err() error {
	return nil
}

// extractHost is a minimal host extractor for URL strings.
func extractHost(rawurl string) string {
	// skip scheme
	start := 0
	if len(rawurl) > 7 && rawurl[:7] == "http://" {
		start = 7
	} else if len(rawurl) > 8 && rawurl[:8] == "https://" {
		start = 8
	} else {
		return ""
	}
	// find end of host
	for i := start; i < len(rawurl); i++ {
		if rawurl[i] == '/' || rawurl[i] == ':' || rawurl[i] == '?' || rawurl[i] == '#' {
			return rawurl[start:i]
		}
	}
	return rawurl[start:]
}

// RobotsPolicy is a minimal interface for robots.txt checks.
// The existing type in robots_policy.go is assumed to implement it.
type RobotsPolicy struct {
	allowed map[string]bool
}

// Allowed returns true if the given URL is permitted by the policy.
func (p *RobotsPolicy) Allowed(rawurl string) bool {
	if p == nil || p.allowed == nil {
		return true
	}
	allowed, ok := p.allowed[rawurl]
	return !ok || allowed
}

// ErrStopIteration signals that the iterator is exhausted.
var ErrStopIteration = errors.New("robots policy iterator exhausted")
