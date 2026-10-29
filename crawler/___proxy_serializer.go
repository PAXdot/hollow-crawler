// crawler/proxy_serializer.go

package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

// ProxySerializer implements a simple JSON serialization for Proxy.
type ProxySerializer struct{}

// Serialize converts a Proxy to a JSON string.
func (s ProxySerializer) Serialize(p *Proxy) (string, error) {
	if p == nil {
		return "", errors.New("cannot serialize nil proxy")
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
	return &p, nil
}

// Proxy represents a proxy with its address and properties.
type Proxy struct {
	IP       net.IP
	Port     int
	Scheme   string
	Type     string
	Username string
	Password string
}
