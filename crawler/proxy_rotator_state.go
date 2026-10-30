package crawler

import (
	"net/http"
	"sync/atomic"
	"time"
)

// ProxyRotatorState represents the internal state of the proxy rotator.
type ProxyRotatorState struct {
	currentProxyIndex int32
	lastUpdated       time.Time
 Proxies             []Proxy
}

// NewProxyRotatorState creates a new ProxyRotatorState instance with the given proxies.
func NewProxyRotatorState(proxies []Proxy) *ProxyRotatorState {
	return &ProxyRotatorState{
		Proxies:             proxies,
		lastUpdated:       time.Now(),
	}
}

// GetCurrentProxy returns the current proxy being used.
func (s *ProxyRotatorState) GetCurrentProxy() Proxy {
	index := atomic.LoadInt32(&s.currentProxyIndex)
	return s.Proxies[index]
}

// UpdateCurrentProxy moves to the next proxy in the rotation.
func (s *ProxyRotatorState) UpdateCurrentProxy() {
	index := atomic.AddInt32(&s.currentProxyIndex, 1)
	if index >= int32(len(s.Proxies)) {
		index = 0
		atomic.StoreInt32(&s.currentProxyIndex, index)
	}
	s.lastUpdated = time.Now()
}

// GetLastUpdated returns the timestamp of the last proxy update.
func (s *ProxyRotatorState) GetLastUpdated() time.Time {
	return s.lastUpdated
}

// Proxy represents a proxy configuration.
type Proxy struct {
	URL          string
	LastFailed   time.Time
	ConsecutiveFailures int
}

// NewProxy creates a new Proxy instance with the given URL.
func NewProxy(url string) Proxy {
	return Proxy{
		URL:          url,
		LastFailed:   time.Time{},
		ConsecutiveFailures: 0,
	}
}

// SetLastFailed updates the last failed timestamp and increments consecutive failures.
func (p *Proxy) SetLastFailed() {
	p.LastFailed = time.Now()
	p.ConsecutiveFailures++
}

// ResetLastFailed resets the last failed timestamp and consecutive failures.
func (p *Proxy) ResetLastFailed() {
	p.LastFailed = time.Time{}
	p.ConsecutiveFailures = 0
}

// ShouldRetry checks if the proxy should be retried based on its failure history.
func (p *Proxy) ShouldRetry(maxFailures int) bool {
	return p.ConsecutiveFailures < maxFailures
}

// IsHealthy checks if the proxy is healthy based on its failure history.
func (p *Proxy) IsHealthy(maxFailures int) bool {
	return p.ConsecutiveFailures == 0
}

// RequestError represents an error encountered during a request using a proxy.
type RequestError struct {
	Err    error
	Proxy  Proxy
}

// NewRequestError creates a new RequestError instance with the given error and proxy.
func NewRequestError(err error, proxy Proxy) RequestError {
	return RequestError{
		Err:    err,
		Proxy:  proxy,
	}
}
