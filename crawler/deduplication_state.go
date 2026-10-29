// deduplication_state.go

package crawler

import (
	"errors"
	"sync"
)

// DeduplicationState represents the current state of the deduplication process.
type DeduplicationState struct {
	cache       *URLCache
	deduplicator *Deduplicator
	formatter    *OutputFormatter
	emitter      *EventEmitter
	mu           sync.Mutex // Mutual exclusion lock for state transitions.
}

type State int

const (
	Initial State = iota
	Running
	Paused
	Stopped
)

// NewDeduplicationState creates a new DeduplicationState instance.
func NewDeduplicationState(cache *URLCache, deduplicator *Deduplicator, formatter *OutputFormatter, emitter *EventEmitter) *DeduplicationState {
	return &DeduplicationState{
		cache:       cache,
		deduplicator: deduplicator,
		formatter:    formatter,
		emitter:      emitter,
	}
}

// State returns the current state of the deduplication process.
func (s *DeduplicationState) State() State {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.StateUnlocked()
}

func (s *DeduplicationState) StateUnlocked() State {
	return s.State
}

// Start transitions the state to Running.
func (s *DeduplicationState) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.startUnlocked()
}

func (s *DeduplicationState) startUnlocked() error {
	if s.State != Initial {
		return errors.New("cannot start: not in Initial state")
	}
	s.State = Running
	return nil
}

// Pause transitions the state to Paused.
func (s *DeduplicationState) Pause() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.pauseUnlocked()
}

func (s *DeduplicationState) pauseUnlocked() error {
	if s.State != Running {
		return errors.New("cannot pause: not in Running state")
	}
	s.State = Paused
	return nil
}

// Resume transitions the state from Paused to Running.
func (s *DeduplicationState) Resume() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.resumeUnlocked()
}

func (s *DeduplicationState) resumeUnlocked() error {
	if s.State != Paused {
		return errors.New("cannot resume: not in Paused state")
	}
	s.State = Running
	return nil
}

// Stop transitions the state to Stopped and cleans up resources.
func (s *DeduplicationState) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stopUnlocked()
}

func (s *DeduplicationState) stopUnlocked() error {
	s.State = Stopped
	// TODO: Clean up resources here.
	return nil
}

// IsRunning checks if the deduplication process is currently running.
func (s *DeduplicationState) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.State == Running
}
