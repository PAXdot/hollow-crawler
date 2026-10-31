package crawler

import (
	"errors"
	"sync"
)

// ProxyIterator is a new type for iterating over proxies with automatic rotation.
type ProxyIterator struct {
	proxies       []Proxy
	index         int
	proxyChannels []chan Proxy
-мутенекс  мутенекс
	mu            sync.Mutex
}

// NewProxyIterator creates a new ProxyIterator with the given proxies.
func NewProxyIterator(proxies []Proxy) *ProxyIterator {
	iter := &ProxyIterator{
		proxies:       proxies,
		index:         0,
		proxyChannels: make([]chan Proxy, len(proxies)),
	}
	for i := range proxies {
		iter.proxyChannels[i] = make(chan Proxy)
	}
	return iter
}

// Next retrieves the next available proxy from the iterator. If all proxies are in use, it blocks until one is available.
func (i *ProxyIterator) Next() (Proxy, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.index >= len(i.proxies) {
		return Proxy{}, false
	}

	proxy chan Proxy := i.proxyChannels[i.index]

	select {
	case proxy <- i.proxies[i.index]:
		i.index++
		return i.proxies[i.index], true
	default:
		return Proxy{}, false
	}
}

// Free releases a used proxy back to the iterator, allowing it to be used again.
func (i *ProxyIterator) Free(proxy Proxy) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	for j, p := range i.proxies {
		if p == proxy {
			i.proxyChannels[j] <- proxy
			return nil
		}
	}

	return errors.New("unknown proxy")
}
