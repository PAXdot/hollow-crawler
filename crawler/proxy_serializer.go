package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

// ProxySerializer provides methods to save and load proxy pool state
// to and from JSON files. Useful for persisting proxy rotation state
// across crawler restarts.
type ProxySerializer struct {
	mu sync.Mutex
}

// ProxyPoolSnapshot is a JSON-serializable representation of a proxy pool.
type ProxyPoolSnapshot struct {
	Proxies []ProxyEntry `json:"proxies"`
}

// ProxyEntry holds a single proxy URL and its current health score.
type ProxyEntry struct {
	URL   string  `json:"url"`
	Score float64 `json:"score"`
}

// NewProxySerializer creates a new ProxySerializer.
func NewProxySerializer() *ProxySerializer {
	return &ProxySerializer{}
}

// Export writes the current proxy pool state to a JSON file.
// It reads the proxy list and their scores from the provided rotator.
func (ps *ProxySerializer) Export(filename string, rotator *ProxyRotator) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if rotator == nil {
		return errors.New("proxy rotator is nil")
	}

	snapshot := ProxyPoolSnapshot{}

	// Collect all proxies and their scores from the rotator.
	// This assumes ProxyRotator has exported methods to list proxies.
	// If not, adapt to your actual rotator API.
	for url, score := range rotator.ProxyScores() {
		snapshot.Proxies = append(snapshot.Proxies, ProxyEntry{
			URL:   url,
			Score: score,
		})
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal proxy snapshot: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("write proxy snapshot: %w", err)
	}

	return nil
}

// Import reads a JSON file and loads the proxy pool state into the rotator.
// Existing proxies in the rotator are not cleared; imported entries are added or updated.
func (ps *ProxySerializer) Import(filename string, rotator *ProxyRotator) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if rotator == nil {
		return errors.New("proxy rotator is nil")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read proxy snapshot: %w", err)
	}

	var snapshot ProxyPoolSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("unmarshal proxy snapshot: %w", err)
	}

	for _, entry := range snapshot.Proxies {
		if err := rotator.AddProxy(entry.URL, entry.Score); err != nil {
			return fmt.Errorf("add proxy %q: %w", entry.URL, err)
		}
	}

	return nil
}
