package crawler

import "net/url"

// CleanFrontierURL normalizes a raw URL string for frontier insertion.
// It removes fragment, decodes percent-encoding, and lowercases scheme/host.
// Returns nil if parsing fails or the URL is empty after trimming spaces.
func CleanFrontierURL(raw string) *url.URL {
	if raw == "" {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil
	}
	u.Fragment = ""
	u.RawQuery = u.Query().Encode()
	u.RawPath = u.EscapedPath()
	if u.Scheme != "" {
		u.Scheme = lowerASCII(u.Scheme)
	}
	if u.Host != "" {
		u.Host = lowerASCII(u.Host)
	}
	if u.String() == "" {
		return nil
	}
	return u
}

// lowerASCII returns a lowercased copy of s for ASCII letters.
func lowerASCII(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}
