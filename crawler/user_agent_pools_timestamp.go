package crawler

import (
	"sync"
	"time"
)

// userAgentPoolsTimestamp keeps track of timestamps for user-agent pools.
type userAgentPoolsTimestamp struct {
	mutex    sync.Mutex
	timestamps map[string]time.Time
}

// NewUserAgentPoolsTimestamp returns a new userAgentPoolsTimestamp.
func NewUserAgentPoolsTimestamp() *userAgentPoolsTimestamp {
	return &userAgentPoolsTimestamp{
		timestamps: make(map[string]time.Time),
	}
}

// SetTimestamp updates the timestamp for the given user agent pool.
func (t *userAgentPoolsTimestamp) SetTimestamp(poolID string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.timestamps[poolID] = time.Now()
}

// GetTimestamp returns the timestamp for the given user agent pool.
func (t *userAgentPoolsTimestamp) GetTimestamp(poolID string) time.Time {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.timestamps[poolID]
}
