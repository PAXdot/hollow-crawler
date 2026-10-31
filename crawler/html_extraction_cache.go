package crawler

import (
	"sync"
	"time"
)

// ExtractFunc represents an HTML extraction function that takes raw HTML
// and returns extracted data or an error.
type ExtractFunc func(html []byte) (interface{}, error)

// ExtractionCacheConfig holds configuration for the extraction cache.
type ExtractionCacheConfig struct {
	// MaxEntries limits the number of cached entries.
	MaxEntries int
	// TTL is the time-to-live for cached entries.
	TTL time.Duration
}

// DefaultExtractionCacheConfig returns a sensible default configuration.
func DefaultExtractionCacheConfig() ExtractionCacheConfig {
	return ExtractionCacheConfig{
		MaxEntries: 10000,
		TTL:        5 * time.Minute,
	}
}

type extractionCacheEntry struct {
	result    interface{}
	err       error
	expiresAt time.Time
}

// CachedExtractor wraps an extraction function with a thread-safe cache.
// It is self-contained and does not depend on other cache implementations
// in the package.
type CachedExtractor struct {
	extract ExtractFunc
	config  ExtractionCacheConfig
	mu      sync.Mutex
	cache   map[string]extractionCacheEntry
	order   []string
}

// NewCachedExtractor creates a new CachedExtractor with the given extraction
// function and configuration. If config has zero values, defaults are applied.
func NewCachedExtractor(extract ExtractFunc, config ExtractionCacheConfig) *CachedExtractor {
	if extract == nil {
		panic("crawler: extraction function cannot be nil")
	}
	if config.MaxEntries <= 0 {
		config.MaxEntries = DefaultExtractionCacheConfig().MaxEntries
	}
	if config.TTL <= 0 {
		config.TTL = DefaultExtractionCacheConfig().TTL
	}
	return &CachedExtractor{
		extract: extract,
		config:  config,
		cache:   make(map[string]extractionCacheEntry),
		order:   make([]string, 0, config.MaxEntries),
	}
}

// Extract runs the wrapped extraction function on the given HTML, caching the
// result keyed by the provided key. Subsequent calls with the same key within
// the TTL return the cached result or error.
func (c *CachedExtractor) Extract(key string, html []byte) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	if entry, ok := c.cache[key]; ok {
		if now.Before(entry.expiresAt) {
			return entry.result, entry.err
		}
		// Remove expired entry and its key from order.
		delete(c.cache, key)
		c.removeOrder(key)
	}

	// Run the extraction outside the lock to avoid holding it during potentially
	// expensive work. However, we must unlock before calling c.extract to avoid
	// deadlock on reentrant calls. Since we hold the lock, we use a temporary
	// unlock.
	c.mu.Unlock()
	result, err := c.extract(html)
	c.mu.Lock()

	// Evict oldest entry if at capacity.
	if len(c.cache) >= c.config.MaxEntries && len(c.order) > 0 {
		oldest := c.order[0]
		delete(c.cache, oldest)
		c.order = c.order[1:]
	}

	// Store new entry.
	c.cache[key] = extractionCacheEntry{
		result:    result,
		err:       err,
		expiresAt: now.Add(c.config.TTL),
	}
	c.order = append(c.order, key)

	return result, err
}

func (c *CachedExtractor) removeOrder(key string) {
	for i, k := range c.order {
		if k == key {
			c.order = append(c.order[:i], c.order[i+1:]...)
			return
		}
	}
}

// Clear removes all cached entries.
func (c *CachedExtractor) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]extractionCacheEntry)
	c.order = c.order[:0]
}

// Len returns the number of entries currently in the cache.
func (c *CachedExtractor) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.cache)
}
