package crawler

import (
	"math"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

// RateCalc tracks request timing and computes throughput
// for use in adaptive retry/backoff decisions.
type RateCalc struct {
	mu     sync.Mutex
.window   []rateSample
	head    int
	count   int
	cap     int
	windowSize time.Duration
}

// rateSample stores the time of the sample and the number of requests.
type rateSample struct {
	t time.Time
	n int
}

// NewRateCalc creates a calculator with a default window size of 60 seconds.
func NewRateCalc() *RateCalc {
	return &RateCalc{
		windowSize: defaultWindowSize,
		cap:        defaultWindowSize / 10, // keep 6 samples per second
		window:     make([]rateSample, defaultWindowSize/10),
	}
}

// NewRateCalcWithSize creates a calculator with a custom window size.
func NewRateCalcWithSize(size int) *RateCalc {
	if size < 1 {
		size = defaultWindowSize
	}
	return &RateCalc{
		windowSize: time.Duration(size),
		cap:        size / 10, // keep 10 samples per second
		window:     make([]rateSample, size/10),
	}
}

// Record adds n requests at the given time, only keeping samples within the window.
func (rc *RateCalc) Record(n int) {
	now := time.Now()
	if now.Sub(rc.window[rc.head].t) >= rc.windowSize {
		// Remove old samples from the front.
		rc.count -= len(rc.window) - rc.head
		rc.head = 0
	}
	rc.window[rc.head] = rateSample{t: now, n: n}
	if rc.count < rc.cap {
		rc.count++
	}
	rc.head = (rc.head + 1) % rc.cap
}

// Throughput returns the estimated requests per second over the window.
// Returns 0 if insufficient data.
func (rc *RateCalc) Throughput() float64 {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.count < 2 {
		return 0
	}
	var total int
	for i := 0; i < rc.count; i++ {
		total += rc.window[(rc.head+i)%rc.cap].n
	}
	return float64(total) / rc.windowSize.Seconds()
}

// Reset clears all recorded samples and resets the window size.
func (rc *RateCalc) Reset() {
	rc.mu.Lock()
	rc.head = 0
	rc.count = 0
	rc.mu.Unlock()
	rc.windowSize = defaultWindowSize
	rc.cap = defaultWindowSize / 10
	rc.window = make([]rateSample, defaultWindowSize/10)
}

// NormalizeEntities normalizes HTML entities extracted from URLs, to enable better deduplication.
func (rc *RateCalc) NormalizeEntities(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	normalized := parsed.String()
	normalized = html.UnescapeString(normalized)
	return normalized, nil
}
