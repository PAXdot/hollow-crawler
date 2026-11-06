package crawler

import (
	"sync"
	"time"
)

// Cache is an in-memory cache with time-based expiration.
type Cache struct {
	cache     map[string]entry
	expiration time.Duration
	mutex     sync.Mutex
}

// NewCache creates a new Cache instance with the specified expiration duration.
func NewCache(exp time.Duration) *Cache {
	return &Cache{
		cache:     make(map[string]entry),
		expiration: exp,
	}
}

// Put inserts or updates an entry in the cache. If the key already exists, an error is returned.
func (c *Cache) Put(key, value string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, ok := c.cache[key]; ok {
		return ErrCacheKeyExists
	}

	c.cache[key] = entry{value: value, expiresAt: time.Now().Add(c.expiration)}
	return nil
}

// Get retrieves an entry from the cache by key. If the entry is expired or not found, an error is returned.
func (c *Cache) Get(key string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry, ok := c.cache[key]
	if !ok || entry.expiresAt.Before(time.Now()) {
		return "", ErrCacheEntryNotFound
	}

	return entry.value, nil
}

type entry struct {
	value     string
	expiresAt time.Time
}

var ErrCacheKeyExists = errors.New("cache: key already exists")
var ErrCacheEntryNotFound = errors.New("cache: entry not found")
```
