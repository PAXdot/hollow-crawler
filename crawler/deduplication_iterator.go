package crawler

import (
	"context"
	"hash/fnv"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/time/rate"
)

// DeduplicationIterator is a self-contained iterator utility for deduplication pipeline.
type DeduplicationIterator struct {
	dedupCache *DedupCache
	ch         chan *ExtractedContent
	ctx        context.Context
	mu         sync.Mutex
	rateLimiter *rate.Limiter
	backoff    time.Duration
	closed     bool
}

// NewDeduplicationIterator creates a new DeduplicationIterator instance with a rate limiter and backoff duration.
func NewDeduplicationIterator(dedupCache *DedupCache, ctx context.Context, rateLimit, backoff time.Duration) *DeduplicationIterator {
	return &DeduplicationIterator{
		dedupCache: dedupCache,
		ch:         make(chan *ExtractedContent),
		ctx:        ctx,
		rateLimiter: rate.NewLimiter(rate.Every(rateLimit), 1),
		backoff:    backoff,
	}
}

// Next returns the next extracted content if available, nil if exhausted or rate limited.
func (di *DeduplicationIterator) Next() *ExtractedContent {
	select {
	case <-di.ctx.Done():
		return nil
	default:
		if !di.rateLimiter.Allow() {
			time.Sleep(di.backoff)
			return nil
		}
		select {
		case di.ch <- struct{}{}:
			select {
			case res := <-di.ch:
				return res
			case <-di.ctx.Done():
				return nil
			}
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

// ExtractContent takes an extracted content, adds it to the deduplication cache, and limits the rate of addition.
func (di *DeduplicationIterator) ExtractContent(content *ExtractedContent) {
	hash := fnv.New32()
	hash.Write([]byte(content.URL))
	cacheKey := hash.Sum32()

	di.mu.Lock()
	defer di.mu.Unlock()

	if !di.rateLimiter.Allow() {
		time.Sleep(di.backoff)
		return
	}

	di.dedupCache.Set(cacheKey, content)

	di.ch <- content
}
