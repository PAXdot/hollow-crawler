package crawler

import (
	"regexp"
)

// TinyParser represents an HTML extractor for high-throughput crawling.
type TinyParser struct {
	hrefPattern *regexp.Regexp
	onLinkFound  chan string // Add an event channel for link discoveries
}

// NewTinyParser returns a new TinyParser instance.
func NewTinyParser() *TinyParser {
	return &TinyParser{
		hrefPattern: regexp.MustCompile(`<a\s+[^>]*href\s*=\s*"([^"]+)"`),
		onLinkFound:  make(chan string),
	}
}

// ExtractLinks returns all href values found in the given HTML string.
func (p *TinyParser) ExtractLinks(html string) []string {
	var links []string
	for _, match := range p.hrefPattern.FindAllStringSubmatch(html) {
		if len(match) >= 2 {
			links = append(links, match[1])
			p.onLinkFound <- match[1] // Emit event for each link found
		}
	}
	close(p.onLinkFound) // Close the channel after we're done
	return links
}

// ListenForLinks blocks until the links are found.
func (p *TinyParser) ListenForLinks() <-chan string {
	return p.onLinkFound
}
