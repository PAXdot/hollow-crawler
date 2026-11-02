package crawler

import (
	"context"
	"hash/fnv"
	"sync"

	"github.com/urfave/cli/v2"
)

// DeduplicationIterator is a self-contained iterator utility for deduplication pipeline.
type DeduplicationIterator struct {
	dedupCache *DedupCache
	ch         chan *ExtractedContent
	ctx        context.Context
	mu         sync.Mutex
	closed     bool
}

// NewDeduplicationIterator creates a new DeduplicationIterator instance.
func NewDeduplicationIterator(dedupCache *DedupCache, ctx context.Context) *DeduplicationIterator {
	return &DeduplicationIterator{
		dedupCache: dedupCache,
		ch:         make(chan *ExtractedContent),
		ctx:        ctx,
	}
}

// Next returns the next extracted content if available, nil if exhausted.
func (di *DeduplicationIterator) Next() *ExtractedContent {
	select {
	case <-di.ctx.Done():
		return nil
	case di.ch <- struct{}{}:
		select {
		case res := <-di.ch:
			return res
		case <-di.ctx.Done():
			return nil
		}
	}
}

// Close stops the iterator and cleans up resources.
func (di *DeduplicationIterator) Close() {
	di.mu.Lock()
	defer di.mu.Unlock()
	if !di.closed {
		close(di.ch)
		di.closed = true
	}
}

// ExtractContent takes an extracted content and adds it to the deduplication cache.
func (di *DeduplicationIterator) ExtractContent(content *ExtractedContent) {
	hash := fnv.New32()
	hash.Write([]byte(content.URL))
	cacheKey := hash.Sum32()

	di.mu.Lock()
	di.dedupCache.Set(cacheKey, content)
	di.mu.Unlock()

	di.ch <- content
}
