package crawler

import (
	"sync"
	"sync/timeoutsync"
	"time"
)

// DedupCache is a concurrency-safe cache that wraps deduplication logic
// and export pipelines. It stores recently seen items to avoid redundant
// processing and exports them in batches.
type DedupCache struct {
	mu       sync.RWMutex
	seen     map[string]time.Time
	ttl      time.Duration
	maxSize  int
	exportFn func(item string) error
	buffer   []string
	bufSize  int
	bufTick  *time.Ticker
	stopCh   chan struct{}
}

// NewDedupCache creates a DedupCache with the given TTL for stored items,
// maximum cache size, export function, and buffer flush interval.
// The export function is called for each unique item or in batches.
func NewDedupCache(ttl time.Duration, maxSize int, exportFn func(item string) error, flushInterval time.Duration) *DedupCache {
	dc := &DedupCache{
		seen:     make(map[string]time.Time),
		ttl:      ttl,
		maxSize:  maxSize,
		exportFn: exportFn,
		buffer:   make([]string, 0, 100),
		bufSize:  100,
		bufTick:  time.NewTicker(flushInterval),
		stopCh:   make(chan struct{}),
	}
	go dc.flushLoop()
	return dc
}

// flushLoop periodically flushes the buffer to the export function.
func (dc *DedupCache) flushLoop() {
	for {
		select {
		case <-dc.bufTick.C:
			dc.flush()
		case <-dc.stopCh:
			dc.flush()
			return
		}
	}
}

// flush sends all buffered items to the export function.
func (dc *DedupCache) flush() {
	dc.mu.Lock()
	if len(dc.buffer) == 0 {
		dc.mu.Unlock()
		return
	}
	batch := make([]string, len(dc.buffer))
	copy(batch, dc.buffer)
	dc.buffer = dc.buffer[:0]
	dc.mu.Unlock()

	for _, item := range batch {
		if err := dc.exportFn(item); err != nil {
			// Optionally log or handle export error; for simplicity, we skip.
		}
	}
}

// Add attempts to insert an item into the cache. Returns true if the item
// is new (not seen within TTL), false if it is a duplicate. On first seen,
// the item is buffered for export.
func (dc *DedupCache) Add(item string) bool {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	now := time.Now()
	if t, ok := dc.seen[item]; ok && now.Sub(t) < dc.ttl {
		return false
	}

	if len(dc.seen) >= dc.maxSize {
		dc.evictExpired(now)
	}

	dc.seen[item] = now
	dc.buffer = append(dc.buffer, item)
	if len(dc.buffer) >= dc.bufSize {
		batch := make([]string, len(dc.buffer))
		copy(batch, dc.buffer)
		dc.buffer = dc.buffer[:0]
		// Spawn a goroutine to export batch without blocking.
		go func(b []string) {
			for _, it := range b {
				dc.exportFn(it)
			}
		}(batch)
	}
	return true
}

// evictExpired removes items that have exceeded their TTL.
func (dc *DedupCache) evictExpired(now time.Time) {
	for k, v := range dc.seen {
		if now.Sub(v) >= dc.ttl {
			delete(dc.seen, k)
		}
	}
}

// Close stops the flush loop and flushes remaining items.
func (dc *DedupCache) Close() {
	close(dc.stopCh)
	dc.bufTick.Stop()
}

// Size returns the current number of items in the cache.
func (dc *DedupCache) Size() int {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return len(dc.seen)
}
