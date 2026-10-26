package crawler

import (
	"regexp"
	"sync"
)

// TinyParser represents an HTML extractor for high-throughput crawling.
type TinyParser struct {
	hrefPattern *regexp.Regexp
	linkChan    chan string // Channel for link discoveries
	mu          sync.Mutex   // Mutex for safe link deduplication
}

// NewTinyParser returns a new TinyParser instance.
func NewTinyParser() *TinyParser {
	return &TinyParser{
		hrefPattern: regexp.MustCompile(`<a\s+[^>]*href\s*=\s*"([^"]+)"`),
		linkChan:    make(chan string),
	}
}

// ExtractLinks extracts all unique href values found in the given HTML string.
func (p *TinyParser) ExtractLinks(html string) []string {
	matches := p.hrefPattern.FindAllStringSubmatch(html, -1)
	links := make(map[string]bool) // Use a map for faster link deduplication

	for _, match := range matches {
		if len(match) >= 2 {
			link := match[1]
			p.mu.Lock()
			if !links[link] {
				links[link] = true
				p.linkChan <- link // Emit event for each unique link found
			}
			p.mu.Unlock()
		}
	}
	close(p.linkChan) // Close the channel after we're done
	linkList := make([]string, 0, len(links))
	for link := range links {
		linkList = append(linkList, link)
	}
	return linkList
}

// LinkChan returns a channel that receives unique links as they're discovered.
func (p *TinyParser) LinkChan() <-chan string {
	return p.linkChan
}
