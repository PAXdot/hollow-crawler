package crawler

import "sync"

// HTMLExtractionEvent represents an event during HTML extraction.
type HTMLExtractionEvent int

const (
	// HTMLExtractionStart indicates the beginning of extraction.
	HTMLExtractionStart HTMLExtractionEvent = iota
	// HTMLExtractionComplete indicates successful extraction.
	HTMLExtractionComplete
	// HTMLExtractionError indicates an extraction error.
	HTMLExtractionError
	// HTMLExtractionSkip indicates extraction was skipped.
	HTMLExtractionSkip
)

// HTMLExtractionListener is a callback for HTML extraction events.
type HTMLExtractionListener func(event HTMLExtractionEvent, data interface{})

// HTMLExtractionEmitter manages event listeners.
type HTMLExtractionEmitter struct {
	mu        sync.RWMutex
	listeners []HTMLExtractionListener
}

// NewHTMLExtractionEmitter creates a new emitter.
func NewHTMLExtractionEmitter() *HTMLExtractionEmitter {
	return &HTMLExtractionEmitter{}
}

// On registers a listener for all events.
func (e *HTMLExtractionEmitter) On(listener HTMLExtractionListener) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.listeners = append(e.listeners, listener)
}

// Off removes a specific listener by value comparison.
func (e *HTMLExtractionEmitter) Off(listener HTMLExtractionListener) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, l := range e.listeners {
		if &l == &listener {
			e.listeners = append(e.listeners[:i], e.listeners[i+1:]...)
			return
		}
	}
}

// Emit sends an event to all registered listeners.
func (e *HTMLExtractionEmitter) Emit(event HTMLExtractionEvent, data interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, listener := range e.listeners {
		listener(event, data)
	}
}
