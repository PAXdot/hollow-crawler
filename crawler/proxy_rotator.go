package crawler

import (
	"context"
	"net/url"
	"sync"
)

// Proxy defines a single proxy endpoint with optional authentication.
type Proxy struct {
	URL      *url.URL
	Username string
	Password string
}

// ProxyRotator is a thread-safe iterator over a list of proxies.
// It supports round‑robin cycling and optional per‑proxy health tracking.
type ProxyRotator struct {
	mu      sync.Mutex
	proxies []Proxy
	idx     int
	healthy map[int]bool // nil means no health tracking
}

// NewProxyRotator creates a rotator from a list of proxy URLs.
// If checkHealth is true, all proxies are initially marked healthy.
func NewProxyRotator(proxies []Proxy, checkHealth bool) *ProxyRotator {
	r := &ProxyRotator{
		proxies: make([]Proxy, len(proxies)),
		idx:     0,
	}
	copy(r.proxies, proxies)
	if checkHealth && len(proxies) > 0 {
		r.healthy = make(map[int]bool, len(proxies))
		for i := range r.proxies {
			r.healthy[i] = true
		}
	}
	return r
}

// Next returns the next proxy in round‑robin order.
// If health tracking is enabled, only healthy proxies are returned.
// It blocks until a healthy proxy is available or the context is cancelled.
func (r *ProxyRotator) Next(ctx context.Context) (Proxy, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.proxies) == 0 {
		return Proxy{}, false
	}

	// Fast path: no health tracking, simple round‑robin.
	if r.healthy == nil {
		p := r.proxies[r.idx]
		r.idx = (r.idx + 1) % len(r.proxies)
		return p, true
	}

	// Slow path: skip unhealthy proxies.
	start := r.idx
	for {
		select {
		case <-ctx.Done():
			return Proxy{}, false
		default:
		}

		if r.healthy[r.idx] {
			p := r.proxies[r.idx]
			r.idx = (r.idx + 1) % len(r.proxies)
			return p, true
		}

		r.idx = (r.idx + 1) % len(r.proxies)
		if r.idx == start {
			// All proxies are unhealthy; block until one is marked healthy.
			r.mu.Unlock()
			select {
			case <-ctx.Done():
				r.mu.Lock()
				return Proxy{}, false
			}
		}
	}
}

// MarkHealthy sets the most recently returned proxy as healthy.
func (r *ProxyRotator) MarkHealthy() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.healthy != nil {
		r.healthy[(r.idx-1+len(r.proxies))%len(r.proxies)] = true
	}
}

// MarkUnhealthy sets the most recently returned proxy as unhealthy.
func (r *ProxyRotator) MarkUnhealthy() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.healthy != nil {
		r.healthy[(r.idx-1+len(r.proxies))%len(r.proxies)] = false
	}
}

// Len returns the total number of proxies (including unhealthy ones).
func (r *ProxyRotator) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.proxies)
}
