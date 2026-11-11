package crawler

import (
	"sync"
	"time"
)

// UserAgentThroughputCalculator tracks the request throughput per user-agent
// and reports the requests per second for each agent in the pool.
type UserAgentThroughputCalculator struct {
	mu             sync.Mutex
	agentsThroughput map[string]throughputRecord
	interval       time.Duration // sliding window duration
}

type throughputRecord struct {
	total        int
	lastUpdated  int64 // unix second of last update
	throughput   float64
}

// NewUserAgentThroughputCalculator creates a calculator with the given sliding window interval.
// interval must be > 0; if zero, the default of 1 second is used.
func NewUserAgentThroughputCalculator(interval time.Duration) *UserAgentThroughputCalculator {
	if interval <= 0 {
		interval = time.Second
	}
	return &UserAgentThroughputCalculator{
		agentsThroughput: make(map[string]throughputRecord),
		interval:         interval,
	}
}

// IncrementRequest records one request for the given user-agent.
func (c *UserAgentThroughputCalculator) IncrementRequest(agent string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().Unix()
	record, ok := c.agentsThroughput[agent]
	if !ok {
		c.agentsThroughput[agent] = throughputRecord{total: 1, lastUpdated: now}
		return
	}
	record.total++
	if now != record.lastUpdated {
		record.lastUpdated = now
		record.throughput = float64(record.total) / float64(now-record.lastUpdated)
	}
}

// Throughput returns the requests per second for the given user-agent in the last interval.
// If no data exists, it returns 0.
func (c *UserAgentThroughputCalculator) Throughput(agent string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	record, ok := c.agentsThroughput[agent]
	if !ok {
		return 0
	}
	return record.throughput
}

// AllThroughputs returns a map of agent -> requests per second for all tracked agents.
func (c *UserAgentThroughputCalculator) AllThroughputs() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make(map[string]float64, len(c.agentsThroughput))
	for agent, record := range c.agentsThroughput {
		out[agent] = record.throughput
	}
	return out
}

// Reset clears all tracked data.
func (c *UserAgentThroughputCalculator) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.agentsThroughput = make(map[string]throughputRecord)
}
