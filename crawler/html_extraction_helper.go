package crawler

import (
	"strings"
)

// htmlExractionHelper provides utility functions for HTML extraction.
type htmlExractionHelper struct{}

// NewHTMLExractionHelper returns a new instance of htmlExractionHelper.
func NewHTMLExractionHelper() *htmlExractionHelper {
	return &htmlExractionHelper{}
}

// ExtractHTMLTitle extracts the title from the given HTML content.
func (h *htmlExractionHelper) ExtractHTMLTitle(htmlContent string) string {
 titres := []string{"title", "TITLE", " Title", " Title "}

	for _, t := range titres {
		if strings.HasPrefix(htmlContent, t) {
			return strings.TrimSpace(htmlContent[len(t):])
		}
	}
	// Default title if none found
	return "No Title"
}

// ExtractHTMLHeaders extracts headers (H1 to H6) from the given HTML content.
func (h *htmlExractionHelper) ExtractHTMLHeaders(htmlContent string) []string {
	var headers []string
	lines := strings.Split(htmlContent, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		for i := 1; i <= 6; i++ {
			if strings.HasPrefix(line, fmt.Sprintf("H%d", i)) {
				headers = append(headers, strings.TrimSpace(line[len(fmt.Sprintf("H%d", i)):]))
				break
			}
		}
	}

	return headers
}
