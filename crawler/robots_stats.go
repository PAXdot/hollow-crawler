package crawler

import (
	"sync"
	"sync/atomic"
	"time"
)

// RobotsStats tracks counts and timing of robots.txt allow/block decisions.
type RobotsStats struct {
	mu         sync.RWMutex
	allowed    atomic.Uint64
	blocked    atomic.Uint64
	cached     atomic.Uint64
	fetched    atomic.Uint64
	lastAccess atomic.Int64 // unix nano
	durations  sync.Map     // map[domain]time.Duration
}

// RecordAllowed increments the allowed counter and updates last access time.
func (rs *RobotsStats) RecordAllowed() {
	rs.allowed.Add(1)
	rs.lastAccess.Store(time.Now().UnixNano())
}

// RecordBlocked increments the blocked counter.
func (rs *RobotsStats) RecordBlocked() {
	rs.blocked.Add(1)
	rs.lastAccess.Store(time.Now().UnixNano())
}

// RecordCacheHit increments the cached counter.
func (rs *RobotsStats) RecordCacheHit() {
	rs.cached.Add(1)
}

// RecordFetch increments the fetched counter and records fetch duration for a domain.
func (rs *RobotsStats) RecordFetch(domain string, d time.Duration) {
	rs.fetched.Add(1)
	rs.durations.Store(domain, d)
}

// Allowed returns the number of allowed decisions.
func (rs *RobotsStats) Allowed() uint64 {
	return rs.allowed.Load()
}

// Blocked returns the number of blocked decisions.
func (rs *RobotsStats) Blocked() uint64 {
	return rs.blocked.Load()
}

// Cached returns the number of cache hits.
func (rs *RobotsStats) Cached() uint64 {
	return rs.cached.Load()
}

// Fetched returns the number of fetched robots.txt files.
func (rs *RobotsStats) Fetched() uint64 {
	return rs.fetched.Load()
}

// LastAccess returns the time of the last recorded access.
func (rs *RobotsStats) LastAccess() time.Time {
	nanos := rs.lastAccess.Load()
	if nanos == 0 {
		return time.Time{}
	}
	return time.Unix(0, nanos)
}

// DomainDuration returns the recorded fetch duration for a domain, if any.
func (rs *RobotsStats) DomainDuration(domain string) (time.Duration, bool) {
	v, ok := rs.durations.Load(domain)
	if !ok {
		return 0, false
	}
	return v.(time.Duration), true
}

// Snapshot returns a point-in-time copy of all stats.
func (rs *RobotsStats) Snapshot() RobotsStatsSnapshot {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	snap := RobotsStatsSnapshot{
		Allowed:    rs.allowed.Load(),
		Blocked:    rs.blocked.Load(),
		Cached:     rs.cached.Load(),
		Fetched:    rs.fetched.Load(),
		LastAccess: rs.LastAccess(),
	}
	rs.durations.Range(func(key, value any) bool {
		if snap.Durations == nil {
			snap.Durations = make(map[string]time.Duration)
		}
		snap.Durations[key.(string)] = value.(time.Duration)
		return true
	})
	return snap
}

// RobotsStatsSnapshot is a read-only copy of RobotsStats.
type RobotsStatsSnapshot struct {
	Allowed    uint64
	Blocked    uint64
	Cached     uint64
	Fetched    uint64
	LastAccess time.Time
	Durations  map[string]time.Duration
}
