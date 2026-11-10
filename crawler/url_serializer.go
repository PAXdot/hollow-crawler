//go:generate go run generate.go

package crawler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// URLSerializer serializes and deserializes URL frontier.
type URLSerializer struct {
	mu sync.Mutex
}

func (s *URLSerializer) Serialize(urls []string, destPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(urls)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *URLSerializer) Deserialize(srcPath string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}

	var urls []string
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *URLSerializer) LoadURLsFromFile(filePath string) ([]string, error) {
	_, filename := filepath.Split(filePath)
	if filename == "" {
		return nil, ErrInvalidFilePath
	}

	return s.Deserialize(filePath)
}

func (s *URLSerializer) SaveURLsToFile(urls []string, filePath string) error {
	_, filename := filepath.Split(filePath)
	if filename == "" {
		return ErrInvalidFilePath
	}

	return s.Serialize(urls, filePath)
}

// ThroughputCalculator tracks URL processing rate over a sliding window.
type ThroughputCalculator struct {
	mu      sync.Mutex
	window  time.Duration
	entries []throughputEntry
}

type throughputEntry struct {
	timestamp time.Time
	count     int
}

// NewThroughputCalculator creates a calculator with the given window duration.
func NewThroughputCalculator(window time.Duration) *ThroughputCalculator {
	if window <= 0 {
		window = time.Minute
	}
	return &ThroughputCalculator{
		window: window,
	}
}

// Record records that n URLs were processed at the current time.
func (tc *ThroughputCalculator) Record(n int) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	now := time.Now()
	tc.entries = append(tc.entries, throughputEntry{timestamp: now, count: n})
	tc.prune(now)
}

// Rate returns the URLs-per-second rate over the configured window.
func (tc *ThroughputCalculator) Rate() float64 {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	now := time.Now()
	tc.prune(now)

	if len(tc.entries) == 0 {
		return 0.0
	}

	// Use the earliest entry's time as the window start for accuracy.
	start := tc.entries[0].timestamp
	elapsed := now.Sub(start)
	if elapsed <= 0 {
		return 0.0
	}

	total := 0
	for _, e := range tc.entries {
		total += e.count
	}

	return float64(total) / elapsed.Seconds()
}

// Reset clears all recorded entries.
func (tc *ThroughputCalculator) Reset() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.entries = nil
}

// prune removes entries older than the window.
func (tc *ThroughputCalculator) prune(now time.Time) {
	cutoff := now.Add(-tc.window)
	i := 0
	for i < len(tc.entries) && tc.entries[i].timestamp.Before(cutoff) {
		i++
	}
	tc.entries = tc.entries[i:]
}
