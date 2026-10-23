package crawler

import (
	"math"
	"sync"
	"time"
)

// RateCalc tracks request timing and computes throughput
// for use in adaptive retry/backoff decisions. It also emits events when proxy rotation occurs.
type RateCalc struct {
	mu     sync.Mutex
	window []rateSample
	head   int
	count  int
	cap    int
	rotateCh chan struct{}
}

// rateSample holds the timestamp and number of requests for a sample.
type rateSample struct {
	t time.Time
	n int
}

// NewRateCalc creates a calculator that keeps up to size samples and emits ProxyRotate events.
func NewRateCalc(size int) *RateCalc {
	if size < 1 {
		size = 100
	}
	rc := &RateCalc{
		window: make([]rateSample, size),
		cap:    size,
		rotateCh: make(chan struct{}),
	}
	go func() {
		for range rc.rotateCh {
			emitEvent(ProxyRotate)
		}
	}()
	return rc
}

// Record adds n requests at the given time. It emits a ProxyRotate event if the head pointer wraps around.
func (rc *RateCalc) Record(n int, t time.Time) {
	if t.IsZero() {
		t = time.Now()
	}
	if n < 1 {
		n = 1
	}
	rc.mu.Lock()
	rc.window[rc.head] = rateSample{t: t, n: n}
	rc.head = (rc.head + 1) % rc.cap
	if rc.count < rc.cap {
		rc.count++
		if rc.head == 0 {
			rc.rotateCh <- struct{}{}
		}
	}
	rc.mu.Unlock()
}

// Throughput returns the estimated requests per second over the window.
// Returns 0 if insufficient data.
func (rc *RateCalc) Throughput() float64 {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.count < 2 {
		return 0
	}
	oldest := rc.window[(rc.head-rc.count+rc.cap)%rc.cap]
	newest := rc.window[(rc.head-1+rc.cap)%rc.cap]
	elapsed := newest.t.Sub(oldest.t)
	if elapsed <= 0 {
		return 0
	}
	var total int
	for i := 0; i < rc.count; i++ {
		idx := (rc.head - rc.count + i + rc.cap) % rc.cap
		total += rc.window[idx].n
	}
	return float64(total) / elapsed.Seconds()
}

// WeightedThroughput computes throughput with exponential weighting,
// giving more weight to recent samples. Returns 0 if insufficient data.
func (rc *RateCalc) WeightedThroughput(alpha float64) float64 {
	if alpha <= 0 || alpha >= 1 {
		alpha = 0.3
	}
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.count < 2 {
		return 0
	}
	var weightSum float64
	var valueSum float64
	for i := 0; i < rc.count; i++ {
		idx := (rc.head - rc.count + i + rc.cap) % rc.cap
		weight := math.Pow(1-alpha, float64(rc.count-1-i))
		weightSum += weight
		valueSum += weight * float64(rc.window[idx].n)
	}
	// estimate total time from oldest to newest sample
	oldest := rc.window[(rc.head-rc.count+rc.cap)%rc.cap]
	newest := rc.window[(rc.head-1+rc.cap)%rc.cap]
	elapsed := newest.t.Sub(oldest.t).Seconds()
	if elapsed <= 0 {
		return 0
	}
	avgWeighted := valueSum / weightSum
	return avgWeighted / elapsed
}

// Reset clears all recorded samples and stops the proxy rotation event emitter.
func (rc *RateCalc) Reset() {
	rc.mu.Lock()
	rc.head = 0
	rc.count = 0
	close(rc.rotateCh)
	rc.mu.Unlock()
}

// ProxyRotate is an event emitted when a proxy rotation occurs.
type ProxyRotate struct{}
