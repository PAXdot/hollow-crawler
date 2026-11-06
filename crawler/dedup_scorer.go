package crawler

import (
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hollow-crawler/hollow-crawler/internal/logger"
)

// scoreURL calculates a deduplication score for a given URL. The score is based on URL length, path depth,
// presence of query parameters, and retry history. Normalizes and handles invalid URLs.
func scoreURL(rawURL string, retryStats *RetryBackoffStats) DedupScore {
	if rawURL == "" {
		return DedupScoreMin
	}

	// Normalize URL by removing query parameters and handling invalid URLs
	u, err := url.Parse(rawURL)
	if err != nil {
		logger.Warnf("Invalid URL '%s': %v", rawURL, err)
		return DedupScoreMin
	}
	normalizedURL := u.String()

	lengthScore := len(normalizedURL)
	if lengthScore > 100 {
		lengthScore = 100
	}

	// ... (rest of the file remains unchanged)
}
