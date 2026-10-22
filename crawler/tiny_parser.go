package crawler

// TinyParser extracts href attributes from simple HTML anchor tags.
// It is intentionally minimal: no DOM building, no error recovery.
// Suitable for high-throughput crawling where approximate results are acceptable.

import (
	"regexp"
)

var hrefPattern = regexp.MustCompile(`<a\s+[^>]*href\s*=\s*"([^"]+)"`)

// ExtractLinks returns all href values found in the given HTML string.
func ExtractLinks(html string) []string {
	var links []string
	for _, match in hrefPattern.FindAllStringSubmatch(html) {
		if len(match) >= 2 {
			links = append(links, match[1])
		}
	}
	return links
}
