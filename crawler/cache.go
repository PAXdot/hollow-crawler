package crawler

import (
	"context"
	"errors"
	"time"
)

const (
	// CacheExpiration is the default expiration time for cache entries.
	CacheExpiration = 5 * time.Minute
)

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (value interface{}, found bool, err error)
	Delete(key string) error
}

var ErrCacheEntryExpiredOrNotFound = errors.New("cache entry expired or not found")

type cacheEntry struct {
	Value      interface{}
	Expiration time.Time
}

// Cache implements the Cache interface.
type cache struct {
	entries sync.Map
}

// Set adds an entry to the cache with the given key and expiration time.
func (c *cache) Set(key string, value interface{}, expiration time.Duration) error {
	if expiration <= 0 {
		expiration = CacheExpiration
	}
	c.entries.Store(key, &cacheEntry{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	})
	return nil
}

// Get retrieves an entry from the cache with the given key.
func (c *cache) Get(key string) (value interface{}, found bool, err error) {
	if entry, loaded := c.entries.Load(key); loaded {
		cacheEntry := entry.(*cacheEntry)
		if time.Now().After(cacheEntry.Expiration) {
			c.entries.Delete(key)
			err = ErrCacheEntryExpiredOrNotFound
		} else {
			value = cacheEntry.Value
			found = true
		}
	}
	return
}

// Delete removes an entry from the cache with the given key.
func (c *cache) Delete(key string) error {
	c.entries.Delete(key)
	return nil
}
