package crawler

import (
	"fmt"
	"html"
	"regexp"
	"time"
)

// TimeStamp extracts and formats the timestamp from a HTML string.
// It looks for the earliest 'datetime' attribute in 'time' tags.
// If none found, it returns an empty string. If the parsed time is in the future,
// it returns a message indicating that the timestamp is invalid.
func TimeStamp(htmlStr string) string {
	reg := regexp.MustCompile(`<time\s[^>]*(?:\s+datetime="?([^"]+)")?`)
	matches := reg.FindStringSubmatch(htmlStr)
	if len(matches) < 2 {
		return ""
	}
	dtStr := matches[1]
	if dtStr == "" {
		return ""
	}
	dt, err := time.Parse(time.RFC3339, dtStr)
	if err != nil || dt.After(time.Now()) {
		return "Invalid timestamp"
	}
	return fmt.Sprintf("Published on %s", html.EscapeString(dt.Format("January 2, 2006 15:04:05 MST")))
}
