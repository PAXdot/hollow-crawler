package crawler

import (
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
	Disallowed  map[string]bool    // path -> disallowed
	CrawlDelay  time.Duration      // default crawl delay for this host
	PathDelays  map[string]time.Duration // path-specific crawl delay overrides
	UserAgent   string             // user-agent these rules apply to
}

// NewRobotState creates a new RobotState with empty rules.
func NewRobotState(userAgent string) *RobotState {
	return &RobotState{
		Disallowed: make(map[string]bool),
		PathDelays: make(map[string]time.Duration),
		UserAgent:  userAgent,
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
	type state struct {
		Disallowed map[string]bool
		CrawlDelay time.Duration
		PathDelays map[string]time.Duration
		UserAgent  string
	}
	s := state{
		Disallowed: r.Disallowed,
		CrawlDelay: r.CrawlDelay,
		PathDelays: r.PathDelays,
		UserAgent:  r.UserAgent,
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
	r.UserAgent = s.UserAgent
	return nil
}

// SaveRobotState writes a gob-encoded RobotState to path.
func SaveRobotState(rs *RobotState, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create robotstate file: %w", err)
	}
	defer f.Close()
	if err := gob.NewEncoder(f).Encode(rs); err != nil {
		return fmt.Errorf("encode robotstate: %w", err)
	}
	return nil
}

// LoadRobotState reads a gob-encoded RobotState from path.
func LoadRobotState(path string) (*RobotState, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open robotstate file: %w", err)
	}
	defer f.Close()
	var rs RobotState
	if err := gob.NewDecoder(f).Decode(&rs); err != nil {
		return nil, fmt.Errorf("decode robotstate: %w", err)
	}
	return &rs, nil
}
