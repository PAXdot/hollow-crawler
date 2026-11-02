package crawler

import (
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ExtractTextFromHTML extracts text from HTML using simple string manipulation.
func ExtractTextFromHTML(html string) string {
	html = strings.ReplaceAll(html, "<[^>]+>", "")
	html = strings.TrimSpace(html)
	return html
}

// ExtractLinksFromHTML extracts all links from HTML using simple string manipulation.
func ExtractLinksFromHTML(html string) ([]string, error) {
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

// RetryGet fetches a URL with retry and backoff.
func RetryGet(url string, client *http.Client) (*http.Response, error) {
	backoff := backoff.NewExponentialBackOff()
	backoff.MaxElapsedTime = 120 * time.Second

	err := backoff.Retry(func() error {
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return ErrNon200Status{resp.StatusCode}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &http.Response{}, nil
}

// ErrNon200Status is returned when a non-200 status code is received.
type ErrNon200Status int

func (e ErrNon200Status) Error() string {
	return "received non-200 status code: " + string(e)
}
