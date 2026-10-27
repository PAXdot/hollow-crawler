package crawler

import (
	"sync"
	"time"
)

// ProxyScore struct to store proxy score and last used time
type ProxyScore struct {
	Score     float64
	LastUsed  time.Time
}

// ProxyScorer service for managing and scoring proxies
type ProxyScorer struct {
	proxies    map[string]ProxyScore
	proxyMutex sync.Mutex
}

// NewProxyScorer creates a new ProxyScorer instance
func NewProxyScorer() *ProxyScorer {
	return &ProxyScorer{
		proxies: make(map[string]ProxyScore),
	}
}

// UpdateScore updates the score of a proxy based on success or failure
func (ps *ProxyScorer) UpdateScore(proxyID string, success bool) {
	ps.proxyMutex.Lock()
	defer ps.proxyMutex.Unlock()

	score, ok := ps.proxies[proxyID]
	if !ok {
		score = ProxyScore{Score: 0, LastUsed: time.Now()}
	}

	if success {
		score.Score += 0.1
	} else {
		score.Score -= 0.1
	}
	score.LastUsed = time.Now()

	ps.proxies[proxyID] = score
}

// GetBestProxy returns the proxy with the highest score that has not been used recently
func (ps *ProxyScorer) GetBestProxy(maxIdleTime time.Duration) (string, bool) {
	ps.proxyMutex.Lock()
	defer ps.proxyMutex.Unlock()

 var bestProxy string
 var bestScore float64
 for proxyID, score := range ps.proxies {
	 if time.Since(score.LastUsed) <= maxIdleTime && score.Score > bestScore {
			bestProxy = proxyID
			bestScore = score.Score
		}
	}

	return bestProxy, bestScore > 0
}

// RemoveProxyFromCache removes a proxy from the score cache
func (ps *ProxyScorer) RemoveProxyFromCache(proxyID string) {
	ps.proxyMutex.Lock()
	defer ps.proxyMutex.Unlock()

	delete(ps.proxies, proxyID)
}
