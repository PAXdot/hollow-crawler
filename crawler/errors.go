package crawler

import (
	"fmt"
	"time"
)

// CrawlError defines the hierarchy of errors produced by the crawler pipeline.
type CrawlError struct {
	Code      int
	Message   string
	Retryable bool
	Inner     error
}

func (e *CrawlError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewCrawlError creates a basic CrawlError.
func NewCrawlError(code int, msg string) *CrawlError {
	return &CrawlError{Code: code, Message: msg, Retryable: false, Inner: nil}
}

// NewRetryableError creates a CrawlError that should trigger a retry.
func NewRetryableError(code int, msg string) *CrawlError {
	return &CrawlError{Code: code, Message: msg, Retryable: true, Inner: nil}
}

// NewWrappedError wraps an existing error into a CrawlError.
func NewWrappedError(code int, msg string, inner error) *CrawlError {
	return &CrawlError{Code: code, Message: msg, Retryable: false, Inner: inner}
}

// IsRetryable returns true if the error indicates a retry is possible.
func IsRetryable(err *CrawlError) bool {
	return err.Retryable
}

// RateLimitError is a convenience constructor for 429-like errors.
func RateLimitError(wait time.Duration) *CrawlError {
	return &CrawlError{
		Code: 429,
		Message: fmt.Sprintf("rate limited, retry after %v", wait),
		Retryable: true,
		Inner: nil,
	}
}

// TimeoutError is a convenience constructor for timeout conditions.
func TimeoutError(url string, dur time.Duration) *CrawlError {
	return &CrawlError{
		Code: 408,
		Message: fmt.Sprintf("timeout fetching %s after %v", url, dur),
		Retryable: true,
		Inner: nil,
	}
}

// RobotsExcludedError is a convenience constructor for robots.txt denial.
func RobotsExcludedError(url string) *CrawlError {
	return &CrawlError{
		Code: 403,
		Message: fmt.Sprintf("robots.txt excludes %s", url),
		Retryable: false,
		Inner: nil,
	}
}
