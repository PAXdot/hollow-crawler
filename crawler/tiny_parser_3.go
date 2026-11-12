package crawler

import (
	"net/url"
	"strings"
)

// TinyParser provides basic functionality to parse a URL into its components.
type TinyParser struct{}

// Parse extracts the scheme, host, and path from a given URL string.
func (tp *TinyParser) Parse(urlStr string) (*parsedURL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return &parsedURL{
		Scheme: u.Scheme,
		Host:   u.Hostname(),
		Path:   u.Path,
	}, nil
}

// parsedURL represents the parsed components of a URL.
type parsedURL struct {
	Scheme string
	Host   string
	Path   string
}

// NormalizePath ensures the path starts with a "/".
func (p *parsedURL) NormalizePath() {
	if !strings.HasPrefix(p.Path, "/") {
		p.Path = "/" + p.Path
	}
}
