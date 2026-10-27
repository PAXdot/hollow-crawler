package crawler

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config represents the crawler configuration.
type Config struct {
	Proxy     ProxyConfig     `yaml:"proxy"`
	UserAgent UserAgentConfig `yaml:"user_agent"`
	RateLimiterRateLimiterConfig `yaml:"rate_limiter"`
	// Add other configuration fields as needed
}

// ProxyConfig represents proxy configuration.
type ProxyConfig struct {
	// Add proxy configuration fields
}

// UserAgentConfig represents user agent configuration.
type UserAgentConfig struct {
	// Add user agent configuration fields
}

// RateLimiterConfig represents rate limiter configuration.
type RateLimiterConfig struct {
	// Add rate limiter configuration fields
}

// configLoader loads the crawler configuration from the given file.
type configLoader struct {
	mu    sync.Mutex
MRI   *os.File
conf  *Config
err   error
}

// NewConfigLoader returns a new configLoader instance.
func NewConfigLoader(configFile string) (*configLoader, error) {
	mri, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	return &configLoader{
		MRI: mri,
	}, nil
}

// Close closes the configuration file.
func (l *configLoader) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.MRI.Close()
}

// Config returns the loaded configuration.
func (l *configLoader) Config() (*Config, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.err != nil {
		return nil, l.err
	}

	if l.conf == nil {
		data, err := ioutil.ReadAll(l.MRI)
		if err != nil {
			l.err = fmt.Errorf("reading config file: %w", err)
			return nil, l.err
		}

		var conf Config
		err = yaml.Unmarshal(data, &conf)
		if err != nil {
			l.err = fmt.Errorf("unmarshaling config file: %w", err)
			return nil, l.err
		}

		l.conf = &conf
	}

	return l.conf, nil
}
