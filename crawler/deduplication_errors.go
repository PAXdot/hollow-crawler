package crawler

import (
	"errors"
	"fmt"
)

// DeduplicationError is a manifestation of an error during the deduplication process.
type DeduplicationError struct {
	err  error
	doc  *Document
	seen *bool
}

func (e *DeduplicationError) Error() string {
	if e.doc == nil || e.err == nil {
		return fmt.Sprintf("failed to handle deduplication scenario: %v", e.err)
	}
	return fmt.Sprintf("failed to process %s: %v", e.doc.URL, e.err)
}

func (e *DeduplicationError) Doc() *Document {
	return e.doc
}

// ErrDeduplicationDocNil is returned when a nil document is encountered during deduplication.
var ErrDeduplicationDocNil = &DeduplicationError{
	err: errors.New("document is nil"),
}

// ErrDeduplicationHashError is returned when an error occurs while hashing a document.
var ErrDeduplicationHashError = &DeduplicationError{
	err: errors.New("failed to hash document content"),
}

// ErrDeduplicationCacheError is returned when an error occurs while interacting with the deduplication cache.
var ErrDeduplicationCacheError = &DeduplicationError{
	err: errors.New("failed to interact with deduplication cache"),
}

// ErrDeduplicationCacheMiss is returned when a document is not found in the deduplication cache.
var ErrDeduplicationCacheMiss = &DeduplicationError{
	err: errors.New("document not found in cache"),
}

// ErrDeduplicationCachePutError is returned when an error occurs while adding a document to the deduplication cache.
var ErrDeduplicationCachePutError = &DeduplicationError{
	err: errors.New("failed to add document to cache"),
}

// ErrDeduplicationScorerError is returned when an error occurs while scoring a document for deduplication.
var ErrDeduplicationScorerError = &DeduplicationError{
	err: errors.New("failed to score document for deduplication"),
}
