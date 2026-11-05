package crawler

import (
	"container/list"
	"context"
	"sync"
)

// ... (Keep existing types and functions)

type absenteeKeyError struct {
	key string
}

func (e *absenteeKeyError) Error() string {
	return "key not present in cache: " + e.key
}

// Delete removes key from the cache. It returns an absenteeKeyError for absent keys.
func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.cache[key]; ok {
		c.removeLocked(el)
	} else {
		return &absenteeKeyError{key: key}
	}
	return nil
}

// ... (Keep existing functions)
