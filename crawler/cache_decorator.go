package crawler

import (
	"sync"
	"time"
)

// CacheDecorator is a wrapper for crawler functions that adds caching functionality.
type CacheDecorator[I any, O any] struct {
	function func(I) (O, error)
	cache    sync.Map
	maxAge   time.Duration
}

// NewCacheDecorator returns a new CacheDecorator with the given function and maximum cache age.
func NewCacheDecorator[I any, O any](function func(I) (O, error), maxAge time.Duration) *CacheDecorator[I, O] {
	return &CacheDecorator[I, O]{
		function: function,
		maxAge:   maxAge,
	}
}

// Call uses the underlying function to compute the result if it's not cached, or retrieves it from the cache.
func (cd *CacheDecorator[I, O]) Call(input I) (O, error) {
	if cached, ok := cd.cache.Load(input); ok {
		return cached.(O), nil
	}

	result, err := cd.function(input)
	if err != nil {
		return result, err
	}

	cd.cache.Store(input, result)
	go func() {
		time.Sleep(cd.maxAge)
		cd.cache.Delete(input)
	}()

	return result, nil
}
