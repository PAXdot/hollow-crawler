package crawler

import (
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

// RetryBackoffStats stores statistics about recent retry and backoff operations.
type RetryBackoffStats struct {
	TotalRetries   int
	MaxConsecutive int
	LastRetryTime  time.Time
}

// DedupScore represents a numeric quality or freshness score for a URL
// during deduplication and export pipeline decisions.
type DedupScore uint32

const (
	// DedupScoreMin is the lowest possible dedup score.
	DedupScoreMin DedupScore = 0
	// DedupScoreMax is the highest possible dedup score.
	DedupScoreMax DedupScore = 100
	// MaxPathDepth is the maximum depth of path segments considered for scoring.
	MaxPathDepth = 15
	// MaxQueryParams is the maximum number of query parameters considered for scoring.
	MaxQueryParams = 5
	// RetryScoreBonus is the score bonus added for each recent retry.
	RetryScoreBonus DedupScore = 5
)

// scoreURL computes a DedupScore for a given URL based on its length,
// path segments, query parameters, and recent retry history. Higher score
// indicates higher priority for export and lower likelihood of being deduplicated away.
// The function is self-contained, deterministic, and uses only standard library.
func scoreURL(rawURL string, retryStats *RetryBackoffStats) DedupScore {
	if rawURL == "" {
		return DedupScoreMin
	}

	lengthScore := len(rawURL)
	if lengthScore > 100 {
		lengthScore = 100
	}

	depth := min(countPathSegments(rawURL), MaxPathDepth)
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

// countPathSegments returns the number of path segments in the URL.
func countPathSegments(url string) int {
	return len(strings.Split(url, "/")) - 1
}

// scoreParamsIfPresent returns the count of query parameters in the URL, up to a maximum.
func scoreParamsIfPresent(url string) int {
	u, _ := url.Parse(url)
	return min(u.RawQueryCount(), MaxQueryParams)
}

// min returns the minimum of the given integers.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// IsHighScore returns true if the score is above a typical threshold (>= 75).
func IsHighScore(s DedupScore) bool {
	return s >= 75
}
