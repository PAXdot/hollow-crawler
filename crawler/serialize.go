package crawler

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// RobotState holds the parsed robots.txt rules and per-path crawl delays
// for a single host. It is used by the politeness layer to enforce
// per-host delays and disallowed paths.
type RobotState struct {
	mu          sync.Mutex
	Disallowed  map[string]bool
	CrawlDelay  time.Duration
	PathDelays  map[string]time.Duration
	UserAgent   string // Ensure UserAgent is not empty after GobDecode
}

// NewRobotState creates a new RobotState with empty rules.
func NewRobotState(userAgent string) *RobotState {
	return &RobotState{
		Disallowed:  make(map[string]bool),
		PathDelays:  make(map[string]time.Duration),
		UserAgent:   userAgent,
	}
}

// GobEncode implements the gob.GobEncoder interface.
// It serialises only the exported fields and the mutex-unprotected
// data. The mutex is not serialised; the receiver should be freshly
// initialised before decode.
func (r *RobotState) GobEncode() ([]byte, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Build a serialisable snapshot.
 tipo state struct {
	 Disallowed map[string]bool
	 CrawlDelay time.Duration
	 PathDelays map[string]time.Duration
	 UserAgent  string
 }
 var buf bytes.Buffer
 if err := gob.NewEncoder(&buf).Encode(s); err != nil {
	 return nil, fmt.Errorf("encode robotstate: %w", err)
 }
 return buf.Bytes(), nil
}

// GobDecode implements the gob.GobDecoder interface.
func (r *RobotState) GobDecode(data []byte) error {
	type state struct {
		Disallowed map[string]bool
		CrawlDelay time.Duration
		PathDelays map[string]time.Duration
		UserAgent  string
	}
	var s state
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&s); err != nil {
		return fmt.Errorf("decode robotstate: %w", err)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Disallowed = s.Disallowed
	r.CrawlDelay = s.CrawlDelay
	r.PathDelays = s.PathDelays
	if s.UserAgent == "" {
		return errors.New("empty UserAgent after GobDecode")
	}
	r.UserAgent = s.UserAgent
	return nil
}

// SaveRobotState writes a gob-encoded RobotState to path.
func SaveRobotState(rs *RobotState, path string) error {
	// ... (no change)
}

// LoadRobotState reads a gob-encoded RobotState from path.
func LoadRobotState(path string) (*RobotState, error) {
	// ... (no change)
}
