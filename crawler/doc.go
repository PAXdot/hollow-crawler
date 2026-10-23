// Package crawler implements a concurrent web crawler and scraper with improved
// deduplication state machine.
package crawler

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// DeduplicationState represents the state of a URL in the deduplication process.
type DeduplicationState int

const (
	// UnknownState indicates the initial state of a URL.
	UnknownState DeduplicationState = iota
	// PendingState indicates that the URL is waiting to be processed.
	PendingState
	// ScrapingState indicates that the URL is being scraped.
	ScrapingState
	// ScrapedState indicates that the URL has been successfully scraped.
	ScrapedState
	// FailedState indicates that the URL failed to be scraped and should be retried.
	FailedState
)

// DeduplicationManager manages the deduplication process of URLs.
type DeduplicationManager struct {
	mu       sync.Mutex
	states   map[string]DeduplicationState
	exportCh chan string
}

// NewDeduplicationManager creates a new DeduplicationManager instance.
func NewDeduplicationManager(exportCh chan string) *DeduplicationManager {
	return &DeduplicationManager{
		states:   make(map[string]DeduplicationState),
		exportCh: exportCh,
	}
}

// Transition updates the state of a URL and returns the next action to be taken.
func (dm *DeduplicationManager) Transition(ctx context.Context, url string, action DeduplicationAction) (DeduplicationAction, error) {
	dm.mu.Lock()
	prevState := dm.states[url]
	dm.states[url] = action.State
	dm.mu.Unlock()

	switch action {
	case AddURLAction:
		if prevState == UnknownState {
			return ScrapeURLAction, nil
		}
		return SkipURLAction, fmt.Errorf("url already in process")
	case ScrapeURLAction:
		switch prevState {
		case PendingState, UnknownState:
			return ScrapeURLAction, nil
		case ScrapedState, FailedState:
			return RetryURLAction, fmt.Errorf("url already processed")
		default:
			return SkipURLAction, fmt.Errorf("invalid state transition")
		}
	case RetryURLAction:
		if prevState == FailedState {
			return ScrapeURLAction, nil
		}
		return SkipURLAction, fmt.Errorf("url not in failed state")
	case ExportURLAction:
		if prevState == ScrapedState {
			select {
			case dm.exportCh <- url:
				return SkipURLAction, nil
			case <-ctx.Done():
				return SkipURLAction, ctx.Err()
			case <-time.After(5 * time.Second):
				return RetryURLAction, fmt.Errorf("export channel blocked")
			}
		}
		return SkipURLAction, fmt.Errorf("url not in scraped state")
	default:
		return SkipURLAction, fmt.Errorf("invalid action")
	}
}
