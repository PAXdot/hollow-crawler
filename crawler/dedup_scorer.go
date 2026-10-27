package crawler

// DedupScore represents a numeric quality or freshness score for a URL
// during deduplication and export pipeline decisions.
type DedupScore uint32

const (
	// DedupScoreMin is the lowest possible dedup score.
	DedupScoreMin DedupScore = 0
	// DedupScoreMax is the highest possible dedup score.
	DedupScoreMax DedupScore = 100
)

// ScoreURL computes a DedupScore for a given URL string based on its length,
// path depth, and query parameters. A higher score indicates higher priority
// for export and lower likelihood of being deduplicated away.
// The function is self-contained, deterministic, and uses only standard library.
func ScoreURL(rawURL string) DedupScore {
	if rawURL == "" {
		return DedupScoreMin
	}
	// Prefer longer URLs as they tend to be more specific.
	lengthScore := len(rawURL)
	if lengthScore > 100 {
		lengthScore = 100
	}
	// Penalize excessive path depth; deeper paths may indicate pagination.
	depth := 0
	for i := 0; i < len(rawURL); i++ {
		if rawURL[i] == '/' {
			depth++
		}
	}
	depthPenalty := depth
	if depthPenalty > 10 {
		depthPenalty = 10
	}
	// Reward presence of query parameters (dynamic content).
	queryBonus := 0
	for i := 0; i < len(rawURL); i++ {
		if rawURL[i] == '?' {
			queryBonus = 15
			break
		}
	}
	raw := uint32(lengthScore) - uint32(depthPenalty) + uint32(queryBonus)
	if raw > uint32(DedupScoreMax) {
		return DedupScoreMax
	}
	return DedupScore(raw)
}

// IsHighScore returns true if the score is above a typical threshold (>= 75).
func IsHighScore(s DedupScore) bool {
	return s >= 75
}
