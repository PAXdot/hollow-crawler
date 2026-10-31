package crawler

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type URLFrontierError struct {
	URL  string
	Err  error
}

func (e *URLFrontierError) Error() string {
	return fmt.Sprintf("error processing URL %q: %v", e.URL, e.Err)
}

func NewURLFrontierError(url string, err error) *URLFrontierError {
	return &URLFrontierError{url, err}
}

// Helper to capture stdout during a function call.
func captureStdout(fn func()) string {
	// ... (no changes to this function)
}

// TestPrintSummaryZeroInputs has been updated with better error handling.
func TestPrintSummaryZeroInputs() {
	f := NewFormatter()
	// Override start time to a known value for deterministic duration.
	oldStart := f.StartTime
	f.StartTime = time.Unix(1000000, 0)
	out := captureStdout(func() {
		err := f.PrintSummary(0, 0, 0)
		if err != nil {
			panic(err)
		}
	})
	lines := strings.Split(out, "\n")
	testing.Assert(len(lines) >= 4, "expected at least 4 lines")
	testing.Assert(strings.Contains(lines[0], "Crawling completed"))
	testing.Assert(strings.Contains(lines[1], "Total URLs processed: 0"))
	testing.Assert(strings.Contains(lines[2], "Successful scrapes: 0"))
	testing.Assert(strings.Contains(lines[3], "Failed scrapes: 0"))
	f.StartTime = oldStart
}

// TestPrintSummaryNonZeroValues is updated with URLFrontierError.
func TestPrintSummaryNonZeroValues() {
	f := NewFormatter()
	f.StartTime = time.Unix(2000000, 0)
	urls := []string{"https://example.com", "https://example.org"}
	frontier := NewURLFrontier(urls)
	errs := []error{
		errors.New("failed to process https://example.com"),
		NewURLFrontierError("https://example.org", errors.New("failed to process https://example.org")),
	}
	frontier.ProcessErrors(errs)
	out := captureStdout(func() {
		f.PrintSummary(frontier.ProcessedURLs(), frontier.SuccessfulScrapes(), len(errs))
	})
	testing.Assert(strings.Contains(out, "Total URLs processed: 2"))
	testing.Assert(strings.Contains(out, "Successful scrapes: 0"))
	testing.Assert(strings.Contains(out, "Failed scrapes: 2"))
	testing.Assert(strings.Contains(out, "https://example.com: failed to process https://example.com"))
	testing.Assert(strings.Contains(out, "https://example.org: failed to process https://example.org"))
}
