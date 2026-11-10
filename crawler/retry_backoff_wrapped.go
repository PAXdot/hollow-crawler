package crawler

import (
	"time"
)

// RetryBackoffWrapper wraps the retry/backoff logic to handle deduplication and export pipeline errors.
// It uses the provided dedup interface for deduplication and the pipeline event emitter for exporting errors.
type RetryBackoffWrapper struct {
	dedup Deduplicator
	ee    *EventEmitter
}

// NewRetryBackoffWrapper returns a new instance of RetryBackoffWrapper.
func NewRetryBackoffWrapper(d Deduplicator, ee *EventEmitter) *RetryBackoffWrapper {
	return &RetryBackoffWrapper{dedup: d, ee: ee}
}

// Wait handles the retry/backoff logic and emits pipeline events.
func (r *RetryBackoffWrapper) Wait(err error, maxRetry int, initialBackoff time.Duration, backoffFactor float64) error {
	// Deduplicate the error
	if err := r.dedup.Deduplicate(err); err != nil {
		r.ee.emitPipelineEvent(PipelineEvent{Type: "deduplicate_error", Error: err})
	}

	// Implement retry/backoff logic based on standard Go library and project conventions
	// ...

	return nil
}
