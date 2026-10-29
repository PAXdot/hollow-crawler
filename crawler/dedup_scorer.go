package crawler

import (
	"strings"
	"unicode/utf8"
)

// DedupScore represents a numeric quality or freshness score for a URL
// during deduplication and export pipeline decisions.
type DedupScore uint32

const (
	// DedupScoreMin is the lowest possible dedup score.
	DedupScoreMin DedupScore = 0
	// DedupScoreMax is the highest possible dedup score.
	DedupScoreMax DedupScore = 100
)

// scoreURL computes a DedupScore for a given URL based on its length,
// path segments, and query parameters. A higher score indicates higher priority
// for export and lower likelihood of being deduplicated away.
// The function is self-contained, deterministic, and uses only standard library.
func scoreURL(rawURL string) DedupScore {
	if rawURL == "" {
		return DedupScoreMin
	}

	lengthScore := len(rawURL)
	if lengthScore > 100 {
		lengthScore = 100
	}

	depth := countPathSegments(rawURL)
	depthPenalty := depth
	if depthPenalty > 10 {
		depthPenalty = 10
	}

	queryBonus := scoreParamsIfPresent(rawURL)

	raw := uint32(lengthScore) - uint32(depthPenalty) + uint32(queryBonus)
	if raw > uint32(DedupScoreMax) {
		return DedupScoreMax
	}

	return DedupScore(raw)
}

// countPathSegments returns the number of path segments in the URL.
func countPathSegments(url string) int {
	return len(strings.Split(url, "/")) - 1
}

// scoreParamsIfPresent returns 15 if the URL has query parameters, 0 otherwise.
func scoreParamsIfPresent(url string) int {
	return strings.Count(url, "?")
}

// IsHighScore returns true if the score is above a typical threshold (>= 75).
func IsHighScore(s DedupScore) bool {
	return s >= 75
}
