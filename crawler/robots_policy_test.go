package crawler

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests for RobotsPolicy
func TestRobotsPolicy(t *testing.T) {
	now := time.Now()
	u, _ := url.Parse("https://example.com/")

	testCases := []struct {
		name     string
		policy   RobotsPolicy
		url      *url.URL
		requestTime time.Time
		expected bool
	}{
		{"Allow All", NewRobotsPolicy(AllowAll), u, now, true},
		{"Disallow All", NewRobotsPolicy(DisallowAll), u, now, false},
		{"Allow Path", NewRobotsPolicy(AllowPath("example.com", "/allow-path")), u, now, true},
		{"Disallow Path", NewRobotsPolicy(DisallowPath("example.com", "/disallow-path")), u, now, false},
		{"Allow User Agent", NewRobotsPolicy(AllowUserAgent("example.com", "test-user-agent")), u, now, true},
		{"Disallow User Agent", NewRobotsPolicy(DisallowUserAgent("example.com", "test-user-agent")), u, now, false},
		{"Allow Time Range", NewRobotsPolicy(AllowTimeRange("example.com", "MON, TUE WED, THU, FRI SAT, SUN", "06:00", "18:00")), u, now, true},
		{"Disallow Time Range", NewRobotsPolicy(DisallowTimeRange("example.com", "MON, TUE WED, THU, FRI SAT, SUN", "06:00", "18:00")), u, now, false},
		{"Allow All but Disallow Time Range", NewRobotsPolicy(DisallowTimeRange("example.com", "MON, TUE WED, THU, FRI SAT, SUN", "06:00", "18:00"), AllowAll), u, now, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.policy.Allow(tc.url, tc.requestTime)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
