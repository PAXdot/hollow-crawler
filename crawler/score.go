package crawler

import (
	"sort"
	"time"
)

// Score represents the score of a URL in the deduplication and export ranking.
type Score struct {
	LastCrawled   time.Time
	SuccessCount  int
	FailureCount  int
	Depth         int
	Importance    float64 // arbitrary importance score
}

// URLScore keeps track of scores for each URL.
type URLScore map[string]Score

// Update updates the score of a URL based on the given crawl result.
func (u URLScore) Update(url string, crawledTime time.Time, success bool) {
	score, ok := u[url]
	if !ok {
		score = Score{
			LastCrawled: time.Time{},
			SuccessCount:  0,
			FailureCount:  0,
			Depth:        0,
			Importance:   1.0, // initial importance score
		}
	}
	score.LastCrawled = crawledTime
	if success {
		score.SuccessCount++
	} else {
		score.FailureCount++
	}
	u[url] = score
}

// RankedURLs returns a slice of URLs sorted by their scores in descending order.
func (u URLScore) RankedURLs() []string {
	scores := make([]Score, 0, len(u))
	for _, score := range u {
		scores = append(scores, score)
	}
	sort.Slice(scores, func(i, j int) bool {
		return compareScores(scores[i], scores[j])
	})
	var rankedURLs []string
	for _, score := range scores {
		rankedURLs = append(rankedURLs, score)
	}
	return rankedURLs
}

// compareScores compares two scores and returns true if the first score should be considered
// "better" than the second score according to the ranking criteria.
func compareScores(s1, s2 Score) bool {
	// Adjust weights as needed to fit the desired ranking criteria.
	return (s1.SuccessCount*2 - s1.FailureCount) > (s2.SuccessCount*2 - s2.FailureCount) &&
		s1.Depth > s2.Depth &&
		s1.Importance > s2.Importance &&
		s1.LastCrawled.After(s2.LastCrawled)
}
