package crawler

import (
	"errors"
	"sync"
)

type DeduplicationState struct {
	// ... existing fields
	extractor *HTMLExtractor
}

// NewDeduplicationState creates a new DeduplicationState instance with additional HTML extractor.
func NewDeduplicationState(cache *URLCache, deduplicator *Deduplicator, formatter *OutputFormatter, emitter *EventEmitter, extractor *HTMLExtractor) *DeduplicationState {
	return &DeduplicationState{
		cache:       cache,
		deduplicator: deduplicator,
		formatter:    formatter,
		emitter:      emitter,
		extractor:    extractor,
	}
}

// ExtractHTML extracts HTML from the given URL and emits the extracted text through the event emitter.
func (s *DeduplicationState) ExtractHTML(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.IsRunning() {
		return errors.New("cannot extract HTML: not running")
	}

	html, err := s.extractor.Extract(url)
	if err != nil {
		return err
	}

	s.emitter.Emit("html_extracted", html)
	return nil
}
