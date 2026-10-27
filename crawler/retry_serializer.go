package crawler

import (
	"encoding/json"
	"fmt"
	"time"
	"strings"
)

// RetryState holds the retry/backoff state for a single crawl target.
type RetryState struct {
	URL       string        `json:"url"`
	Attempts  int           `json:"attempts"`
	MaxRetries int          `json:"max_retries"`
	Backoff   time.Duration `json:"backoff"`
	Deadline  time.Time     `json:"deadline"`
}

// String returns a human-readable string representing the RetryState.
func (r RetryState) String() string {
	return fmt.Sprintf("URL: %s, Attempts: %d, Max Retries: %d, Backoff: %v, Deadline: %v",
		r.URL, r.Attempts, r.MaxRetries, r.Backoff, r.Deadline)
}

// RetrySerializer provides methods to serialize and deserialize RetryState.
type RetrySerializer struct{}

// Serialize converts a RetryState into a JSON byte slice with pretty printing.
func (rs RetrySerializer) Serialize(r RetryState) ([]byte, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("retry serializer: marshal: %w", err)
	}
	return data, nil
}

// Deserialize parses a JSON byte slice into a RetryState.
func (rs RetrySerializer) Deserialize(data []byte) (RetryState, error) {
	var r RetryState
	if err := json.Unmarshal(data, &r); err != nil {
		return RetryState{}, fmt.Errorf("retry serializer: unmarshal: %w", err)
	}
	return r, nil
}
