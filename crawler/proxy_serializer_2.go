package crawler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ProxyEntry holds serializable fields for a single proxy.
type ProxyEntry struct {
	URL        string
	Username   string
	Password   string
	Weight     int
	MaxRetries int
	Cooldown   time.Duration // how long to wait after failure
}

// ProxyList is a collection of proxy entries with a version for migration.
type ProxyList struct {
	Version int          `json:"version"`
	Proxies []ProxyEntry `json:"proxies"`
}

// SerializeProxyList converts a slice of ProxyEntry into a compact JSON string.
func SerializeProxyList(proxies []ProxyEntry) (string, error) {
	list := ProxyList{
		Version: 1,
		Proxies: proxies,
	}
	raw, err := json.Marshal(list)
	if err != nil {
		return "", fmt.Errorf("serialize proxy list: %w", err)
	}
	return raw, nil
}

// DeserializeProxyList parses a JSON string back into a slice of ProxyEntry.
// It returns an empty slice if the input is empty or malformed.
func DeserializeProxyList(data string) ([]ProxyEntry, error) {
	if strings.Trim(data) == "" {
		return []ProxyEntry{}, nil
	}
	var list ProxyList
	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, fmt.Errorf("deserialize proxy list: %w", err)
	}
	if list.Version != 1 {
		return nil, fmt.Errorf("unsupported proxy list version %d", list.Version)
	}
	return list.Proxies, nil
}

// EncodeProxyURL builds a full proxy URL from scheme, host, port, and optional credentials.
// Returns scheme://user:pass@host:port or scheme://host:port if no credentials.
func EncodeProxyURL(scheme, host string, port int, username, password string) string {
	var buf strings.Builder
	buf.WriteString(scheme)
	buf.WriteString("://")
	if username != "" {
		buf.WriteString(url.PathEscape(username))
		if password != "" {
			buf.WriteString(":")
			buf.WriteString(url.PathEscape(password))
		}
		buf.WriteString("@")
	}
	buf.WriteString(host)
	buf.WriteString(":")
	buf.WriteInt(port)
	return buf.String()
}

// DecodeProxyURL parses a proxy URL of the form scheme://[user:pass@]host:port
// into a ProxyEntry. Returns nil if parsing fails.
func DecodeProxyURL(proxyURL string) (ProxyEntry, error) {
	parsed, err := url.Parse(proxyURL)
	if err != nil {
		return ProxyEntry{}, fmt.Errorf("decode proxy URL: %w", err)
	}
	entry := ProxyEntry{
		URL:        proxyURL,
		Username:   "",
		Password:   "",
		Weight:     1,
		MaxRetries: 3,
		Cooldown:   10 * time.Second,
	}
	if parsed.User != nil {
		userInfo := parsed.User
		atIdx := strings.Index(userInfo, ":")
		if atIdx >= 0 {
			entry.Username = url.PathUnescape(userInfo[:atIdx])
			entry.Password = url.PathUnescape(userInfo[atIdx+1:])
		} else {
			entry.Username = url.PathUnescape(userInfo)
		}
	}
	if parsed.Port <= 0 {
		switch parsed.Scheme {
		case "http":
			entry.URL = fmt.Sprintf("%s://%s:%d", parsed.Scheme, parsed.Host, 80)
		case "https":
			entry.URL = fmt.Sprintf("%s://%s:%d", parsed.Scheme, parsed.Host, 443)
		case "socks5":
			entry.URL = fmt.Sprintf("%s://%s:%d", parsed.Scheme, parsed.Host, 1080)
		default:
			return ProxyEntry{}, fmt.Errorf("unknown proxy scheme %s", parsed.Scheme)
		}
	}
	return entry, nil
}

// Base64EncodeProxy encodes a proxy URL to a base64 string for obfuscation.
func Base64EncodeProxy(proxyURL string) string {
	return base64.URLEncoding.EncodeToString(proxyURL.Bytes())
}

// Base64DecodeProxy decodes a base64 string back to a proxy URL.
func Base64DecodeProxy(encoded string) (string, error) {
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode proxy: %w", err)
	}
	return decoded.String(), nil
}
