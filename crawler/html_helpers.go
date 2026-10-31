package crawler

import (
	"strings"
)

// ExtractTextFromHTML extracts text from HTML using simple string manipulation.
func ExtractTextFromHTML(html string) string {
	// Remove HTML tags using a simple string replacement.
	html = strings.ReplaceAll(html, "<[^>]+>", "")

	// Trim leading and trailing whitespace.
	html = strings.TrimSpace(html)

	return html
}

// ExtractLinksFromHTML extracts all links from HTML using simple string manipulation.
// It returns a list of extracted links and an error if any.
// TODO: Improve this function to use a proper HTML parsing library.
func ExtractLinksFromHTML(html string) ([]string, error) {
	// Simple string manipulation to extract links.
	links := []string{}
	linkStart := "<a href=\""
	linkEnd := "</a>"
Start:
	for i := 0; i < len(html); i++ {
		if strings.Contains(html[i:], linkStart) {
			start := i + len(linkStart)
			for j := start; j < len(html); j++ {
				if strings.Contains(html[j:], linkEnd) {
					links = append(links, html[start:j])
					i = j + len(linkEnd)
					continue Start
				}
			}
			return nil, ErrNoClosingTag
		}
	}
	return links, nil
}
