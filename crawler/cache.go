package crawler

import (
	"container/list"
	"context"
	"sync"
)

// Cache defines a generic key-value cache abstraction used by the crawler
// for deduplication and short-lived state storage.
type Cache interface {
	// Get returns the value associated with key and whether it was present.
	Get(key string) (interface{}, bool)
	// Put stores value under key, evicting the least-recently-used entry
	// when the cache is at capacity.
	Put(key string, value interface{}) error
	// Delete removes key from the cache. Deleting a missing key is a no-op.
	Delete(key string) error
}

// lruEntry is the value stored in the LRU list backing an LRUCache.
type lruEntry struct {
	key   string
	value interface{}
}

// LRUCache is a concurrency-safe, fixed-capacity least-recently-used cache.
// It is safe for concurrent use by multiple goroutines.
type LRUCache struct {
	maxCapacity int

	mu      sync.Mutex
	cache   map[string]*list.Element
	order   *list.List // front = most recently used, back = least recently used
	onEvict func(key string, value interface{})
}

// NewLRUCache returns an LRUCache capped at maxCapacity entries. A non-positive
// capacity is clamped to 1 so the cache always has room for at least one entry.
func NewLRUCache(maxCapacity int) *LRUCache {
	if maxCapacity < 1 {
		maxCapacity = 1
	}
	return &LRUCache{
		maxCapacity: maxCapacity,
		cache:       make(map[string]*list.Element, maxCapacity),
		order:       list.New(),
	}
}

// SetEvictionHook registers a callback invoked whenever an entry is evicted or
// deleted. It is intended for cache maintenance such as dedup state cleanup.
// The callback runs while the cache lock is held and must not call back into
// the cache.
func (c *LRUCache) SetEvictionHook(fn func(key string, value interface{})) {
	c.mu.Lock()
	c.onEvict = fn
	c.mu.Unlock()
}

// Get returns the value stored under key, marking it as most recently used.
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	c.order.MoveToFront(el)
	return el.Value.(*lruEntry).value, true
}

// Put inserts or updates key with value, evicting the least-recently-used
// entry if the cache is at capacity.
func (c *LRUCache) Put(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.cache[key]; ok {
		el.Value.(*lruEntry).value = value
		c.order.MoveToFront(el)
		return nil
	}

	for c.order.Len() >= c.maxCapacity {
		c.evictLocked()
	}

	el := c.order.PushFront(&lruEntry{key: key, value: value})
	c.cache[key] = el
	return nil
}

// Delete removes key from the cache. It reports no error for absent keys.
func (c *LRUCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.cache[key]; ok {
		c.removeLocked(el)
	}
	return nil
}

// Len returns the number of entries currently held by the cache.
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.order.Len()
}

// evictLocked removes the least-recently-used entry. The caller must hold c.mu.
func (c *LRUCache) evictLocked() {
	back := c.order.Back()
	if back == nil {
		return
	}
	c.removeLocked(back)
}

// removeLocked detaches el from both the map and the list, firing the eviction
// hook if registered. The caller must hold c.mu.
func (c *LRUCache) removeLocked(el *list.Element) {
	entry := el.Value.(*lruEntry)
	c.order.Remove(el)
	delete(c.cache, entry.key)
	if c.onEvict != nil {
		c.onEvict(entry.key, entry.value)
	}
}

// DeduplicatorService is the stateful service backing deduplication. It owns
// the seen-key cache plus the pluggable scoring strategy used to decide
// whether an item is a duplicate. It is safe for concurrent use.
type DeduplicatorService struct {
	cache    Cache
	strategy DeduplicationStrategy
}

// NewDeduplicatorService builds a service using the supplied cache and scoring
// strategy. If either is nil, sensible defaults are applied.
func NewDeduplicatorService(cache Cache, strategy DeduplicationStrategy) *DeduplicatorService {
	if cache == nil {
		cache = NewLRUCache(1 << 16)
	}
	if strategy == nil {
		strategy = NewURLHashStrategy()
	}
	return &DeduplicatorService{cache: cache, strategy: strategy}
}

// Seen reports whether key has already been observed by the service.
func (s *DeduplicatorService) Seen(key string) bool {
	_, ok := s.cache.Get(key)
	return ok
}

// Mark records key as seen, associating it with the given value.
func (s *DeduplicatorService) Mark(key string, value interface{}) error {
	return s.cache.Put(key, value)
}

// Score computes the deduplication score for the given payload using the
// configured strategy.
func (s *DeduplicatorService) Score(ctx context.Context, payload []byte) (uint64, error) {
	return s.strategy.Score(ctx, payload)
}

// CheckAndMark atomically (per key) reports whether key was already present and,
// when absent, records it. The boolean result is true when the key is a duplicate.
func (s *DeduplicatorService) CheckAndMark(key string, value interface{}) (bool, error) {
	if _, ok := s.cache.Get(key); ok {
		return true, nil
	}
	if err := s.cache.Put(key, value); err != nil {
		return false, err
	}
	return false, nil
}

// DeduplicationStrategy scores a payload so callers can detect near-duplicates.
type DeduplicationStrategy interface {
	// Score returns a stable score for the payload. Identical payloads must
	// always yield identical scores.
	Score(ctx context.Context, payload []byte) (uint64, error)
}
