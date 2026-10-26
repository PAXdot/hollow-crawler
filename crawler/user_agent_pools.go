package crawler

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// UserAgentPool represents the configuration for user-agent pools.
type UserAgentPool struct {
	Pools map[string][]string `yaml:"pools"`
}

// LoadConfig loads the user-agent pool configuration from the given file path.
func LoadConfig(filePath string) (*UserAgentPool, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config UserAgentPool
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Validate checks if the configuration is valid.
func (c *UserAgentPool) Validate() error {
	for _, pool := range c.Pools {
		if len(pool) == 0 {
			return fmt.Errorf("empty user-agent pool found")
		}
	}

	return nil
}
