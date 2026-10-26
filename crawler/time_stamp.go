package crawler

import (
	"fmt"
	"html"
	"time"
)

// TimeStamp extracts and formats the timestamp from a HTML string.
// It looks for the earliest 'datetime' attribute in 'time' tags.
// If none found, it returns an empty string.
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
	dt, _ := time.Parse(time.RFC3339, dtStr)
	return fmt.Sprintf("Published on %s", html.EscapeString(dt.Format(time.RFC1123)))
}
