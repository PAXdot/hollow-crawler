package crawler

import (
	"bytes"
	"encoding/gob"
	"io"
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

// MarshalBinary serializes the cache to a binary format with gob encoding.
func (c *Cache) MarshalBinary() ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(c.cache)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// UnmarshalBinary deserializes the cache from a binary format with gob encoding.
func (c *Cache) UnmarshalBinary(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&c.cache)
	if err != nil {
		return err
	}
	c.pruneExpiredEntries()
	return nil
}

func (c *Cache) pruneExpiredEntries() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

.now := time.Now()
	for key, entry := range c.cache {
		if entry.expiresAt.Before(now) {
			delete(c.cache, key)
		}
	}
}

type entry struct {
	value     string
	expiresAt time.Time
}

var ErrCacheKeyExists = errors.New("cache: key already exists")
var ErrCacheEntryNotFound = errors.New("cache: entry not found")
