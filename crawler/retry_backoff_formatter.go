package crawler

import (
	"fmt"
	"strings"
	"time"
)

// RetryBackoffFormatter produces human-readable strings from retry/backoff data.
type RetryBackoffFormatter struct {
	buf strings.Builder
}

// NewRetryBackoffFormatter creates a new RetryBackoffFormatter.
func NewRetryBackoffFormatter() *RetryBackoffFormatter {
	return &RetryBackoffFormatter{}
}

// FormatAttempt formats a single retry attempt with attempt number, delay, and optional error.
func (f *RetryBackoffFormatter) FormatAttempt(attempt int, delay time.Duration, err error) string {
	f.buf.Reset()
	f.buf.WriteString(fmt.Sprintf("attempt %d: delay %s", attempt, delay))
	if err != nil {
		f.buf.WriteString(fmt.Sprintf(", error: %v", err))
	}
	return f.buf.String()
}

// FormatStats formats retry statistics from a RetryStats instance.
func (f *RetryBackoffFormatter) FormatStats(s *RetryStats) string {
	f.buf.Reset()
	f.buf.WriteString("retry stats: ")
	f.buf.WriteString(fmt.Sprintf("total attempts: %d", s.TotalAttempts))
	f.buf.WriteString(fmt.Sprintf(", retries: %d", s.Retries))
	f.buf.WriteString(fmt.Sprintf(", failures: %d", s.Failures))
	f.buf.WriteString(fmt.Sprintf(", success: %d", s.Success))
	f.buf.WriteString(fmt.Sprintf(", last attempt: %s", s.LastAttempt.Format(time.RFC3339)))
	return f.buf.String()
}

// FormatSummary creates a compact one-line summary.
func (f *RetryBackoffFormatter) FormatSummary(s *RetryStats) string {
	f.buf.Reset()
	f.buf.WriteString(fmt.Sprintf("attempts=%d retries=%d failures=%d success=%d",
		s.TotalAttempts, s.Retries, s.Failures, s.Success))
	return f.buf.String()
}

// FormatBackoffStep formats a single backoff step given its index and duration.
func (f *RetryBackoffFormatter) FormatBackoffStep(step int, d time.Duration) string {
	return fmt.Sprintf("step %d: %s", step, d)
}

// FormatBackoffSequence formats an ordered list of backoff durations.
func (f *RetryBackoffFormatter) FormatBackoffSequence(backoffs []time.Duration) string {
	if len(backoffs) == 0 {
		return "backoff sequence: (empty)"
	}
	f.buf.Reset()
	f.buf.WriteString("backoff sequence: [")
	for i, d := range backoffs {
		if i > 0 {
			f.buf.WriteString(", ")
		}
		f.buf.WriteString(d.String())
	}
	f.buf.WriteString("]")
	return f.buf.String()
}
