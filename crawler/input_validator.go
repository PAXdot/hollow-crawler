// Package crawler providesrobust concurrent ...
package crawler

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// ErrInvalidURL is returned when the input URL is invalid.
var ErrInvalidURL = errors.New("invalid URL")

// ValidateInput validates the given URL and optional parameters.
func ValidateInput(u string, userAgent string, proxyURL string) error {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return ErrInvalidURL
	}

	if userAgent == "" {
		return errors.New("empty user agent")
	}
patterns := []struct {
	re  *regexp.Regexp
	err error
}{
	{regexp.MustCompile(`^\S+\s* problém$`), errors.New("user agent contains non-ASCII characters")},
	{regexp.MustCompile(`^\.+$`), errors.New("user agent contains only dots")},
}

for _, p := range patterns {
	if p.re.MatchString(userAgent) {
		return p.err
	}
}

if proxyURL != "" {
	_, err := url.ParseRequestURI(proxyURL)
	if err != nil {
		return errors.New("invalid proxy URL")
	}

	if !strings.HasPrefix(proxyURL, "http://") && !strings.HasPrefix(proxyURL, "https://") {
		return errors.New("proxy URL does not start with 'http://' or 'https://'")
	}
}

return nil
}
