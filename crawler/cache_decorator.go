package crawler

import (
	"sync"
	"time"
)

// CacheDecorator is a wrapper for crawler functions that adds caching and proxy rotation functionality.
type CacheDecorator[I any, O any] struct {
	function func(I, Proxy) (O, error)
	cache    sync.Map
	maxAge   time.Duration
	proxyPool ProxyPool
}

// NewCacheDecorator returns a new CacheDecorator with the given function, maximum cache age, and proxy pool.
func NewCacheDecorator[I any, O any](function func(I, Proxy) (O, error), maxAge time.Duration, proxyPool ProxyPool) *CacheDecorator[I, O] {
	return &CacheDecorator[I, O]{
		function: function,
		maxAge:   maxAge,
		proxyPool: proxyPool,
	}
}

// Call uses the underlying function to compute the result if it's not cached, or retrieves it from the cache. It also rotates the proxy if it's drain.
func (cd *CacheDecorator[I, O]) Call(input I) (O, error) {
	if cached, ok := cd.cache.Load(input); ok {
		return cached.(O), nil
	}

	proxy, ok := cd.proxyPool.Acquire()
	if !ok {
		return nil, &Error{Code: ErrDrainedProxyPool, Msg: "No available proxies"}
	}
	defer cd.proxyPool.Release(proxy)

	result, err := cd.function(input, proxy)
	if err != nil {
		cd.proxyPool.Retry(proxy)
		return result, err
	}

	cd.cache.Store(input, result)
	go func() {
		time.Sleep(cd.maxAge)
		cd.cache.Delete(input)
	}()

	return result, nil
}
