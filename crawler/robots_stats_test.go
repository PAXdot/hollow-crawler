package crawler

import (
	"testing"
	"time"
)

func TestRobotsStats(t *testing.T) {
	rs := &RobotsStats{}

	// Test initial values
	if rs.Allowed() != 0 {
		t.Error("Allowed counter should be 0")
	}
	if rs.Blocked() != 0 {
		t.Error("Blocked counter should be 0")
	}
	if rs.Cached() != 0 {
		t.Error("Cached counter should be 0")
	}
	if rs.Fetched() != 0 {
		t.Error("Fetched counter should be 0")
	}
	if !rs.LastAccess().IsZero() {
		t.Error("LastAccess should be zero time")
	}
	_, ok := rs.DomainDuration("any")
	if ok {
		t.Error("DomainDuration should return false for non-existing domain")
	}

	// Record some decisions and test counters
	rs.RecordAllowed()
	rs.RecordBlocked()
	rs.RecordCacheHit()
	rs.RecordFetch("example.com", 150*time.Millisecond)
	if rs.Allowed() != 1 {
		t.Error("Allowed counter should be 1")
	}
	if rs.Blocked() != 1 {
		t.Error("Blocked counter should be 1")
	}
	if rs.Cached() != 1 {
		t.Error("Cached counter should be 1")
	}
	if rs.Fetched() != 1 {
		t.Error("Fetched counter should be 1")
	}
	d, ok := rs.DomainDuration("example.com")
	if !ok || d != 150*time.Millisecond {
		t.Errorf("DomainDuration should return %v, %t", 150*time.Millisecond, true)
	}

	// Test Snapshot
	snap := rs.Snapshot()
	if snap.Allowed != 1 {
		t.Error("Allowed in snapshot should be 1")
	}
	if snap.Blocked != 1 {
		t.Error("Blocked in snapshot should be 1")
	}
	if snap.Cached != 1 {
		t.Error("Cache in snapshot should be 1")
	}
	if snap.Fetched != 1 {
		t.Error("Fetched in snapshot should be 1")
	}
	if !snap.LastAccess.IsZero() {
		t.Error("LastAccess in snapshot should be zero time")
	}
	if d, ok := snap.Durations["example.com"]; !ok || d != 150*time.Millisecond {
		t.Errorf("DomainDuration in snapshot should return %v, true", 150*time.Millisecond)
	}

	// Test LastAccess updates
	time.Sleep(10 * time.Millisecond)
	rs.RecordAllowed()
	if rs.LastAccess().Before(time.Now()) {
		t.Error("LastAccess should be after current time")
	}
}
