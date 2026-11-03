package crawler

import (
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

// ...

func scoreURL(rawURL string, retryStats *RetryBackoffStats) DedupScore {
	if rawURL == "" {
		return DedupScoreMin
	}

	// Normalize URL by removing query parameters
	u, err := url.Parse(rawURL)
	if err != nil {
		// Handle invalid URL, e.g., return a low score or log and return the raw URL
		log.Printf("Warning: Invalid URL '%s': %v\n", rawURL, err)
		return DedupScoreMin
	}
	normalizedURL := u.String()

	lengthScore := len(normalizedURL)
	if lengthScore > 100 {
		lengthScore = 100
	}

	depth := min(countPathSegments(normalizedURL), MaxPathDepth)
	depthPenalty := DedupScore(depth)
	if depthPenalty > DedupScoreMax/2 {
		depthPenalty = DedupScoreMax / 2
	}

	queryBonus := scoreParamsIfPresent(rawURL)
	if queryBonus > MaxQueryParams {
		queryBonus = MaxQueryParams
	}

	retryBonus := DedupScore(retryStats.TotalRetries) * RetryScoreBonus
	if time.Since(retryStats.LastRetryTime) > 5*time.Minute {
		retryBonus = 0
	}

	raw := uint32(lengthScore) - uint32(depthPenalty) + uint32(queryBonus) + uint32(retryBonus)
	if raw > uint32(DedupScoreMax) {
		return DedupScoreMax
	}

	return DedupScore(raw)
}
