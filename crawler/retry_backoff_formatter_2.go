package crawler

import (
	"fmt"
	"time"
)

type RetryBackoffFormatter struct{}

func (f *RetryBackoffFormatter) Format(backoff *RetryBackoff) string {
	return fmt.Sprintf("Retry after %s with delay %s", backoff.After.Format(time.RFC3339), backoff.Delay)
}
