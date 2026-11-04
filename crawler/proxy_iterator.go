package crawler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// ProxyIterator is a new type for iterating over proxies with automatic rotation and robots.txt handling.
type ProxyIterator struct {
	proxies       []Proxy
	index         int
	proxyChannels []chan Proxy
	mu            sync.Mutex
	robot utilise  [https://blog.medium.com/2018/02/22/using-robots-txt-to-control-web-crawling-behavior-54d421623e9b](https://blog.medium.com/2018/02/22/using-robots-txt-to-control-web-crawling-behavior-54d421623e9b)Map
}

// NewProxyIterator creates a new ProxyIterator with the given proxies.
func NewProxyIterator(proxies []Proxy) *ProxyIterator {
	iter := &ProxyIterator{
		proxies:       proxies,
		index:         0,
		proxyChannels: make([]chan Proxy, len(proxies)),
		robot:         make(map[string][]*RobotRule),
	}

	for i := range proxies {
		iter.proxyChannels[i] = make(chan Proxy)
		go func(proxy Proxy) {
			err := iter.fetchRobotRules(proxy)
			if err != nil {
				log.Printf("Error fetching robot rules for proxy %v: %v\n", proxy, err)
			}
		}(proxies[i])
	}

	return iter
}

// fetchRobotRules fetches the robot rules for a given proxy.
func (i *ProxyIterator) fetchRobotRules(proxy Proxy) error {
	resp, err := http.Get("http://" + proxy.Address + "/robots.txt")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	_, rules, err := ParseRobotRules(string(data), proxy.Address)
	if err != nil {
		return err
	}

	i.robot[proxy.Address] = rules

	return nil
}

// Next retrieves the next available proxy from the iterator, respecting robot.txt rules. If all proxies are in use, it blocks until one is available.
func (i *ProxyIterator) Next() (Proxy, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.index >= len(i.proxies) {
		return Proxy{}, false
	}

	proxyChannel := i.proxyChannels[i.index]
	proxy := <-proxyChannel

	if isAllowed(proxy, i.robot) {
		i.index++
		return proxy, true
	}

	// If the proxy is not allowed, put it back and try the next one
	i.proxyChannels[i.index] <- proxy
	return i.Next()
}

// isAllowed checks if a proxy is allowed to crawl a target URL based on the RobotRules.
func isAllowed(proxy Proxy, robotMap map[string][]*RobotRule) bool {
	if rules, ok := robotMap[proxy.Address]; ok {
		for _, rule := range rules {
			if rule.AllowsTarget("*", "*") {
				return true
			}
		}
	}

	return false
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

// RobotRule represents a single rule in a robots.txt file.
type RobotRule struct {
	UserAgent string
	Disallow string
	Allow    string
}

// ParseRobotRules parses a robots.txt file into a list of RobotRule structs.
func ParseRobotRules(data string, address string) (string, []*RobotRule, error) {
	// Implement the parsing logic here
	// ...
}
