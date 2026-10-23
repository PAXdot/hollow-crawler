// Package crawler provides utility functions for high-throughput concurrent crawling and scraping.
package crawler

import (
	"fmt"
	"os"
	"time"
)

// Formatter formats and prints human-readable output.
type Formatter struct {
	StartTime time.Time
}

// NewFormatter creates a new Formatter instance.
func NewFormatter() *Formatter {
	return &Formatter{StartTime: time.Now()}
}

// PrintSummary prints a summary of the crawling process, including duration and stats.
func (f *Formatter) PrintSummary(totalURLs, successfulScrapes, failedScrapes int) {
	duration := time.Since(f.StartTime)
	fmt.Fprintf(os.Stdout, "\nCrawling completed in %v\n", duration)
	fmt.Fprintf(os.Stdout, "Total URLs processed: %d\n", totalURLs)
	fmt.Fprintf(os.Stdout, "Successful scrapes: %d\n", successfulScrapes)
	fmt.Fprintf(os.Stdout, "Failed scrapes: %d\n", failedScrapes)
}
