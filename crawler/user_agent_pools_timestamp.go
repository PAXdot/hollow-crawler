FILENAME: crawler/user_agent_pools_timestamp.go
COMMIT: refactor(userAgentPoolsTimestamp): Use time.Now().UnixNano for better precision

package crawler

import (
	"sync"
	"time"
)

type userAgentPoolsTimestamp struct {
	mutex    sync.Mutex
	timestamps map[string]int64 // Use Unix nano seconds for better precision
}

func NewUserAgentPoolsTimestamp() *userAgentPoolsTimestamp {
	return &userAgentPoolsTimestamp{
		timestamps: make(map[string]int64),
	}
}

// SetTimestamp updates the timestamp for the given user agent pool using Unix nano seconds.
func (t *userAgentPoolsTimestamp) SetTimestamp(poolID string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.timestamps[poolID] = time.Now().UnixNano()
}

// GetTimestamp returns the timestamp for the given user agent pool as a Unix nano second value.
func (t *userAgentPoolsTimestamp) GetTimestamp(poolID string) int64 {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.timestamps[poolID]
}
