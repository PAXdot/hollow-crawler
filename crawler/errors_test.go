package crawler

import (
	"errors"
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
		{"wrapped error", NewWrappedError(504, "Gateway timeout", errors.New("backend failed")), false, "[504] Gateway timeout"},
		{"rate limit", RateLimitError(30 * time.Second), true, "[429] rate limited, retry after 30s"},
		{"timeout", TimeoutError("https://example.com", 10 * time.Second), true, "[408] timeout fetching https://example.com after 10s"},
		{"robots exclusion", RobotsExcludedError("https://example.com"), false, "[403] robots.txt excludes https://example.com"},
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
