package crawler

import (
	"sync"
	"time"
)

// PipelineStage identifies the part of the crawl/export pipeline that emits an
// event. It is used by PipelineEvent to let subscribers discriminate between
// deduplication and export activity.
type PipelineStage int

const (
	// StageDedup marks events originating from the deduplicator (cache hits,
	// misses, scorer decisions).
	StageDedup PipelineStage = iota
	// StageExport marks events originating from the export/formatter pipeline.
	StageExport
)

// String returns a short human readable name for the stage.
func (s PipelineStage) String() string {
	switch s {
	case StageDedup:
		return "dedup"
	case StageExport:
		return "export"
	default:
		return "unknown"
	}
}

// PipelineEventKind classifies the concrete event being emitted.
type PipelineEventKind string

const (
	// EventDedupHit is emitted when a URL is found in the dedup cache.
	EventDedupHit PipelineEventKind = "dedup.hit"
	// EventDedupMiss is emitted when a URL is not present in the dedup cache.
	EventDedupMiss PipelineEventKind = "dedup.miss"
	// EventDedupScored is emitted after a dedup score is computed.
	EventDedupScored PipelineEventKind = "dedup.scored"
	// EventExportStarted is emitted when an export batch begins.
	EventExportStarted PipelineEventKind = "export.started"
	// EventExportFinished is emitted when an export batch completes.
	EventExportFinished PipelineEventKind = "export.finished"
	// EventExportFailed is emitted when an export batch fails.
	EventExportFailed PipelineEventKind = "export.failed"
)

// String returns the string form of the kind.
func (k PipelineEventKind) String() string { return string(k) }

// PipelineEvent is the payload delivered to subscribers. Fields are populated
// on a best-effort basis; consumers must not assume any field is non-zero.
type PipelineEvent struct {
	Stage      PipelineStage
	Kind       PipelineEventKind
	URL        string
	Count      int
	Score      float64
	Err        error
	OccurredAt time.Time
}

// PipelineEventHandler is a callback invoked synchronously for each event.
// Handlers run while the emitter holds a read lock, so they must not call back
// into the emitter to avoid deadlocks.
type PipelineEventHandler func(PipelineEvent)

// PipelineEventEmitter fan-outs pipeline events to registered handlers. It is
// safe for concurrent use. Handlers are invoked in registration order.
type PipelineEventEmitter struct {
	mu       sync.RWMutex
	nextID   uint64
	handlers map[uint64]PipelineEventHandler
	closed   bool
}

// NewPipelineEventEmitter returns an empty, ready to use emitter.
func NewPipelineEventEmitter() *PipelineEventEmitter {
	return &PipelineEventEmitter{handlers: make(map[uint64]PipelineEventHandler)}
}

// Subscribe registers h and returns an identifier that can be passed to
// Unsubscribe. A nil handler is ignored and yields a zero identifier.
func (e *PipelineEventEmitter) Subscribe(h PipelineEventHandler) uint64 {
	if h == nil {
		return 0
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return 0
	}
	e.nextID++
	id := e.nextID
	e.handlers[id] = h
	return id
}

// Unsubscribe removes the handler registered under id. Unknown identifiers are
// silently ignored.
func (e *PipelineEventEmitter) Unsubscribe(id uint64) {
	e.mu.Lock()
	delete(e.handlers, id)
	e.mu.Unlock()
}

// Emit delivers ev to all current subscribers. It is a no-op once Close has
// been called. OccurredAt is filled in when left zero.
func (e *PipelineEventEmitter) Emit(ev PipelineEvent) {
	if ev.OccurredAt.IsZero() {
		ev.OccurredAt = time.Now()
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.closed {
		return
	}
	for _, h := range e.handlers {
		h(ev)
	}
}

// Len reports the number of registered handlers.
func (e *PipelineEventEmitter) Len() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.handlers)
}

// Close drops all handlers and prevents further emissions. It is idempotent.
func (e *PipelineEventEmitter) Close() {
	e.mu.Lock()
	e.closed = true
	e.handlers = make(map[uint64]PipelineEventHandler)
	e.mu.Unlock()
}
