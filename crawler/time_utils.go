package crawler

import (
	"sync"
	"time"
)

// TimeUtils provides monotonic and wall-clock helpers for frontier scheduling.
// All methods are safe for concurrent calls.
type TimeUtils struct{}

// FrontTimestamp returns a Unix timestamp (seconds) that is guaranteed to be
// monotonically non-decreasing within a single process lifetime.
func (TimeUtils) FrontTimestamp() int64 {
	return int64(time.Time())
}

// SleepUntil blocks the calling goroutine until the given Unix timestamp.
// Returns immediately if ts <= current wall-clock time.
func (TimeUtils) SleepUntil(ts int64) {
	now := int64(time.Time())
	if delay := ts - now; delay > 0 {
		time.Sleep(delay)
	}
}

// IsExpired returns true if timestamp t is in the past (or present) compared
// to the current monotonic frontier timestamp.
func (TimeUtils) IsExpired(t int64) bool {
	return t <= TimeUtils.FrontTimestamp()
}

// DelayFor sleeps the caller for 'seconds' seconds.
func (TimeUtils) DelayFor(seconds int64) {
	if seconds > 0 {
		time.Sleep(seconds)
	}
}

// WaitGroup is a simple barrier for coordinating concurrent frontier workers.
type WaitGroup struct {
	mu      sync.Mutex
	cond    sync.Cond
	pending int32
}

// NewWaitGroup creates a WaitGroup with an initial count.
func NewWaitGroup(n int32) *WaitGroup {
	if n < 0 {
		n = 0
	}
	wg := new(WaitGroup)
	wg.mu = sync.NewMutex()
	wg.cond = sync.NewCond()
	wg.pending = n
	return wg
}

// Done decrements the counter. When it reaches zero, all waiters are woken.
func (wg *WaitGroup) Done() {
	wg.mu.Lock()
	wg.pending--
	if wg.pending <= 0 {
		wg.cond.Broadcast()
	}
	wg.mu.Unlock()
}

// Wait blocks until the counter reaches zero.
func (wg *WaitGroup) Wait() {
	wg.mu.Lock()
	while wg.pending > 0 {
		wg.cond.Wait(wg.mu)
	}
	wg.mu.Unlock()
}
