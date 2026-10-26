package crawler

import (
	"time"
)

// DefaultCrawlDelay is the default delay used when a robots.txt
// does not specify a Crawl-delay directive.
const DefaultCrawlDelay = 1 * time.Second

// MaxCrawlDelay is the upper bound for a crawl delay to avoid
// excessively long waits between requests.
const MaxCrawlDelay = 5 * time.Minute

// ResolveCrawlDelay returns the effective crawl delay for a given
// user agent, applying the specified delay if it is positive and
// within the maximum allowed bounds. Otherwise, it falls back to
// the provided default.
func ResolveCrawlDelay(robotsDelay time.Duration, ua string) time.Duration {
	// In a real implementation, this could consult a per-user-agent
	// override table. For now, it simply validates and returns the
	// delay.
	if robotsDelay <= 0 {
		return DefaultCrawlDelay
	}
	if robotsDelay > MaxCrawlDelay {
		return MaxCrawlDelay
	}
	_ = ua // Reserved for future user-agent specific adjustments.
	return robotsDelay
}
