package crawler

import (
	"sync"
	"time"
)

// UserAgentRateCalc tracks request throughput per user-agent
// and reports items/second for each agent in the pool.
type UserAgentRateCalc struct {
	mu       sync.Mutex
	agents   map[string]*agentRecord
	interval time.Duration // sliding window duration
}

type agentRecord struct {
	count   int
	lastSec int64 // unix second of last update
}

// NewUserAgentRateCalc creates a calculator with the given sliding window interval.
// interval must be > 0; if zero, the default of 1 second is used.
func NewUserAgentRateCalc(interval time.Duration) *UserAgentRateCalc {
	if interval <= 0 {
		interval = time.Second
	}
	return &UserAgentRateCalc{
		agents:   make(map[string]*agentRecord),
		interval: interval,
	}
}

// Increment records one request for the given user-agent.
func (c *UserAgentRateCalc) Increment(agent string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Unix()
	r, ok := c.agents[agent]
	if !ok {
		c.agents[agent] = &agentRecord{count: 1, lastSec: now}
		return
	}
	if now != r.lastSec {
		// new second, reset count
		r.count = 1
		r.lastSec = now
	} else {
		r.count++
	}
}

// Rate returns the requests/second for the given user-agent in the last interval.
// If no data exists, it returns 0.
func (c *UserAgentRateCalc) Rate(agent string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	r, ok := c.agents[agent]
	if !ok {
		return 0
	}
	now := time.Now().Unix()
	elapsed := now - r.lastSec
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed == 0 {
		return float64(r.count)
	}
	return float64(r.count) / float64(elapsed)
}

// AllRates returns a map of agent -> requests/second for all tracked agents.
func (c *UserAgentRateCalc) AllRates() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Unix()
	out := make(map[string]float64, len(c.agents))
	for agent, r := range c.agents {
		elapsed := now - r.lastSec
		if elapsed < 0 {
			elapsed = 0
		}
		if elapsed == 0 {
			out[agent] = float64(r.count)
		} else {
			out[agent] = float64(r.count) / float64(elapsed)
		}
	}
	return out
}

// Reset clears all tracked data.
func (c *UserAgentRateCalc) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.agents = make(map[string]*agentRecord)
}
