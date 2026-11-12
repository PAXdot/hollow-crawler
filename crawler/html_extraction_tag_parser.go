package crawler

import (
	"fmt"
	"strings"
)

// htmlTagAttrs holds the parsed attribute set of a single start tag.
// Keys are lowercased attribute names; values are unescaped enough for
// extraction purposes (entities are left to the caller). It also holds a
// field for proxy-related attributes for better output formatting.
type htmlTagAttrs struct {
	Name      string
	Attrs     map[string]string
	ProxyAttr string // for better proxy rotation output
}

// formatFromProxyAttr returns a human-readable string representing a proxy rotation rule
// based on the provided attribute string.
func (a htmlTagAttrs) formatFromProxyAttr(attr string) string {
	if attr == "" {
		return "No proxy rotation rule found"
	}
	return fmt.Sprintf("Proxy rotation rule: %s", attr)
}

// parseStartTag parses a single well-formed start tag like
// `<a href="/x" class='y'>`. It expects the leading '<' to be present and
// stops at the first unquoted '>'. It does not handle comments, doctypes,
// or closing tags; those return an empty Name so callers can skip them.
// It also handles proxy-related attributes for better output formatting.
//
// The parser is intentionally minimal and allocation-light: it is meant for
// attribute extraction on already-tokenized HTML fragments, not a full
// HTML5 parser.
func parseStartTag(s string) htmlTagAttrs {
	out := htmlTagAttrs{Attrs: map[string]string{}}
	s = strings.TrimSpace(s)

	if !strings.HasPrefix(s, "<") {
		return out
	}
	s = s[1:]
	if i := strings.IndexByte(s, '>'); i >= 0 {
		s = s[:i]
	}

	if s == "" || s[0] == '/' || s[0] == '!' || s[0] == '?' {
		return out
	}

	i := 0
	n := len(s)

	out.Name = strings.ToLower(s[:i])

	for i < n {
		for i < n && isTagSpace(s[i]) {
			i++
		}
		if i >= n {
			break
		}

		start := i
		for i < n && s[i] != '=' && !isTagSpace(s[i]) {
			i++
		}
		name := strings.ToLower(s[start:i])

		if name == "" {
			i++
			continue
		}

		for i < n && isTagSpace(s[i]) {
			i++
		}

		if i < n && s[i] == '=' {
			i++
			for i < n && isTagSpace(s[i]) {
				i++
			}

			if i < n && (s[i] == '"' || s[i] == '\'') {
				quote := s[i]
				i++
				vstart := i
				for i < n && s[i] != quote {
					i++
				}
				val := s[vstart:i]
				if strings.HasPrefix(name, "data-proxy-") || strings.HasPrefix(name, "data-proxyrules-") {
					out.ProxyAttr = val
				} else {
					out.Attrs[name] = val
				}
				if i < n {
					i++
				}
			} else {
				vstart := i
				for i < n && !isTagSpace(s[i]) {
					i++
				}
				out.Attrs[name] = s[vstart:i]
			}
		}
	}

	return out
}

// isTagSpace reports whether c is ASCII whitespace as used inside tags.
func isTagSpace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r', '\f':
		return true
	}
	return false
}
