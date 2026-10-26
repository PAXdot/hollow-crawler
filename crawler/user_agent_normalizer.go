package crawler

import (
	"regexp"
	"strings"
)

// userAgentNormalizer normalizes user-agent strings.
type userAgentNormalizer struct{}

// Normalize takes a user-agent string and returns a normalized version.
func (u userAgentNormalizer) Normalize(ua string) string {
	ua = strings.ToLower(ua)
	ua = regexp.MustCompile(`\s+`).ReplaceAllString(ua, " ")
	ua = regexp.MustCompile(`/[^\d\s.]+|;|,`).ReplaceAllString(ua, "")
	return ua
}
