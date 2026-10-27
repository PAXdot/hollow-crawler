package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		wantErr  bool
		want     *UserAgentPool
	}{
		{
			name:     "valid config",
			filePath: "testdata/user_agent_pool_valid.yaml",
			want: &UserAgentPool{
				Pools: map[string][]string{
					"desktop": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/5.37.36 Edge/16.16299"},
					"mobile":  {"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_3 like Mac OS X) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.0 Mobile/14G60 Safari/602.1"},
				},
			},
		},
		{
			name:     "invalid config",
			filePath: "testdata/user_agent_pool_invalid.yaml",
			wantErr:  true,
		},
		{
			name:     "file not found",
			filePath: "/nonexistent/file",
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := LoadConfig(tc.filePath)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, config)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		config   *UserAgentPool
		wantErr  bool
	}{
		{
			name: "valid config",
			config: &UserAgentPool{
				Pools: map[string][]string{
					"desktop": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/5.37.36 Edge/16.16299"},
				},
			},
		},
		{
			name:     "empty pool",
			config:   &UserAgentPool{Pools: map[string][]string{"desktop": {}}},
			wantErr:  true,
		},
		{
			name:     "no pools",
			config:   &UserAgentPool{Pools: map[string][]string{}},
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
