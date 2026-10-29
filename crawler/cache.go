package crawler

import (
	"time"
	"sync"
	"context"
	"sync/atomic"
)

type Cache interface {
	Get(key string) (value interface{}, ok bool)
	Put(key string, value interface{}) error
	Delete(key string) error
}

type cacheEntry struct {
	value  interface{}
	exists int32
}

type LRUCache struct {
	maxCapacity int
	mu          sync.Mutex
	cache       map[string]cacheEntry
irsi         * מט
}

func NewLRUCache(maxCapacity int) *LRUCache {
	return &LRUCache{
		maxCapacity: maxCapacity,
		cache:       make(map[string]cacheEntry),
		irsi:        &מט{},
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached, exists := c.cache[key]; exists {
		atomic.StoreInt32(&cached.exists, 1)
		return cached.value, true
	}
	return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.cache) == c.maxCapacity {
		c.irsi.LRU()
	}
	c.cache[key] = cacheEntry{value: value, exists: 1}
	return nil
}

func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
	return nil
}
