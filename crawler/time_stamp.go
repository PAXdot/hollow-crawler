package crawler

import (
	"fmt"
	"html"
	"regexp"
	"time"
)

// TimeStamp extracts and formats the earliest valid timestamp from a HTML string.
// It searches for 'datetime' attributes in 'time' tags, preferring the earliest.
// If none found or all are invalid, it returns an empty string.
func TimeStamp(htmlStr string) string {
	reg := regexp.MustCompile(`<time\s[^>]*(?:datetime="?([^"]+)")?`)
	matches := reg.FindAllStringSubmatch(htmlStr, -1)
	if len(matches) == 0 {
		return ""
	}

	var earliestTime time.Time
	for _, match := range matches {
		if match[1] == "" {
			continue
		}
		dt, err := time.Parse(time.RFC3339, match[1])
		if err != nil || dt.After(time.Now()) {
			continue
		}
		if earliestTime.IsZero() || dt.Before(earliestTime) {
			earliestTime = dt
		}
	}

	if !earliestTime.IsZero() {
		return fmt.Sprintf("Published on %s", html.EscapeString(earliestTime.Format("January 2, 2006 15:04:05 MST")))
	}

	return ""
}
