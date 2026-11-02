package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
)

// ProxySerializer implements a simple JSON serialization for Proxy.
type ProxySerializer struct{}

// Serialize converts a Proxy to a JSON string. It checks if the proxy is valid before serialization.
func (s ProxySerializer) Serialize(p *Proxy) (string, error) {
	if err := validateProxy(p); err != nil {
		return "", err
	}
	data, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("json marshal failed: %v", err)
	}
	return string(data), nil
}

// Deserialize converts a JSON string to a Proxy.
func (s ProxySerializer) Deserialize(data string) (*Proxy, error) {
	var p Proxy
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %v", err)
	}
	if err := validateProxy(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Proxy represents a proxy with its address and properties.
type Proxy struct {
	IP       string `json:"ip"` // Use string for IP to make it human-readable
	Port     int    `json:"port"`
	Scheme   string `json:"scheme"`
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Password string `json:"-"` // Exclude password from JSON output
}

func (p *Proxy) MarshalJSON() ([]byte, error) {
	type ProxyAlias Proxy
	_ = (*ProxyAlias)(p)
	return json.Marshal(ProxyAlias(p))
}

func (p *Proxy) UnmarshalJSON(data []byte) error {
	type ProxyAlias Proxy
	alias := &ProxyAlias{}
	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}
	*p = Proxy(alias)
	return nil
}

// validateProxy checks if the proxy is valid.
func validateProxy(p *Proxy) error {
	if p == nil {
		return errors.New("cannot serialize nil proxy")
	}
	if p.IP == "" || p.Port == 0 || p.Scheme == "" || p.Type == "" {
		return errors.New("invalid proxy: missing field(s)")
	}
	if _, err := net.ResolveTCPAddr(p.Scheme, fmt.Sprintf("%s:%d", p.IP, p.Port)); err != nil {
		return fmt.Errorf("invalid proxy address: %v", err)
	}
	if _, err := regexp.MatchString(`^(http|https)$`, p.Scheme); err != nil {
		return fmt.Errorf("invalid scheme: %s", p.Scheme)
	}
	return nil
}
