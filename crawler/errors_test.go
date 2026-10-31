package crawler

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCrawlErrorBehavior(t *testing.T) {
	tests := []struct {
		name    string
		err     *CrawlError
		want    bool
		wantErr string
	}{
		{"not retryable", NewCrawlError(500, "Internal error"), false, "[500] Internal error"},
		{"retryable", NewRetryableError(503, "Service unavailable"), true, "[503] Service unavailable"},
		{"wrapped error", NewWrappedError(504, "Gateway timeout", errors.New("backend failed")), false, "[504] Gateway timeout: backend failed"},
		{"rate limit", RateLimitError(30 * time.Second), true, "[429] Rate limited. Wait for 30s before retrying"},
		{"timeout", TimeoutError("https://example.com", 10 * time.Second), true, "[408] Timeout fetching https://example.com. Waited for 10s"},
		{"robots exclusion", RobotsExcludedError("https://example.com"), false, "[403] Excluded by robots.txt for https://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.want {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.want)
			}
			if got := tt.err.Error(); got != tt.wantErr {
				t.Errorf("Error() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

// IsRetryable renamed to ShouldRetry. Improves clarity and concurrently updates error message formatting.
func ShouldRetry(err error) bool {
	// ... (no changes in implementation)
}

// FormatError improved to include relevant details in the error message.
func (e *CrawlError) Error() string {
	if e == nil {
		return "nil crawl error"
	}

	message := fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)

	switch e := e.(type) {
	case *RateLimitError:
		message += fmt.Sprintf(". Wait for %v before retrying", e.RetryAfter)
	case *TimeoutError:
		message += fmt.Sprintf(". Waited for %v", e.Timeout)
	case *RobotsExcludedError:
		message += fmt.Sprintf(". Excluded by robots.txt for %s", e.URL)
	case *WrappedError:
		message += fmt.Sprintf(". Caused by: %v", e.Cause)
	}

	return message
}
