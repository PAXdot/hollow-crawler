package crawler

import (
	"strings"
	"testing"
)

// Tests for UserAgentPoolSampler

func TestNewUserAgentPoolSampler(t *testing.T) {
	s := NewUserAgentPoolSampler(42)
	t.Assert(s != nil, "sampler should not be nil")
	t.Assert(s.rng != nil, "rng should be initialized")
}

func TestPickEmptyPools(t *testing.T) {
	s := NewUserAgentPoolSampler(42)
	result := s.Pick(map[string][]string{}, nil)
	t.Assert(result == "", "should return empty string for empty pools")
}

func TestPickUniformWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(12345)
	pools := map[string][]string{
		"a": []string{"ua1", "ua2"},
		"b": []string{"ua3", "ua4"},
		"c": []string{"ua5"},
	}
	// Collect many picks to verify distribution
	seen := map[string]int{}
	for i := 0; i < 1000; i++ {
		ua := s.Pick(pools, nil)
		t.Assert(ua != "", "should not be empty")
		seen[ua]++
	}
	// All user agents should appear at least once
	for _, agents := range pools {
		for _, a := range agents {
			t.Assert(seen[a] > 0, "each user-agent should appear: " + a)
		}
	}
}

func TestPickWithWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(999)
	pools := map[string][]string{
		"light": []string{"l1", "l2"},
		"heavy": []string{"h1", "h2"},
	}
	weights := map[string]float64{
		"light": 1.0,
		"heavy": 100.0,
	}
	seen := map[string]int{}
	for i := 0; i < 500; i++ {
		ua := s.Pick(pools, weights)
		seen[ua]++
	}
	// Heavy pool should dominate
	heavyCount := seen["h1"] + seen["h2"]
	lightCount := seen["l1"] + seen["l2"]
	t.Assert(heavyCount > lightCount, "heavy pool should be selected more often")
}

func TestPickNegativeWeight(t *testing.T) {
	s := NewUserAgentPoolSampler(7)
	pools := map[string][]string{
		"a": []string{"ua1"},
		"b": []string{"ua2"},
	}
	weights := map[string]float64{
		"a": -5.0,
		"b": 10.0,
	}
	result := s.Pick(pools, weights)
	t.Assert(result != "", "should pick something")
	t.Assert(result == "ua2", "should pick from pool with positive weight")
}

func TestPickAllZeroWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(1)
	pools := map[string][]string{
		"a": []string{"ua1"},
	}
	weights := map[string]float64{
		"a": 0.0,
	}
	result := s.Pick(pools, weights)
	t.Assert(result == "", "should return empty when total weight is zero")
}

func TestPickPoolSelectionEdgeCase(t *testing.T) {
	s := NewUserAgentPoolSampler(42)
	pools := map[string][]string{
		"only": []string{"unique_ua"},
	}
	result := s.Pick(pools, nil)
	t.Assert(result == "unique_ua", "should select from the only pool")
}

func TestPickFromPool(t *testing.T) {
	s := NewUserAgentPoolSampler(100)
	agents := []string{"alpha", "beta", "gamma"}
	seen := map[string]int{}
	for i := 0; i < 300; i++ {
		ua := s.PickFromPool(agents)
		t.Assert(ua != "", "should not be empty")
		seen[ua]++
	}
	for _, a := range agents {
		t.Assert(seen[a] > 0, "each agent should appear: " + a)
	}
}

func TestPickFromPoolEmpty(t *testing.T) {
	s := NewUserAgentPoolSampler(0)
	result := s.PickFromPool([]string{})
	t.Assert(result == "", "should return empty for empty list")
}

func TestPickFromPoolSingle(t *testing.T) {
	s := NewUserAgentPoolSampler(77)
	agents := []string{"sole"}
	for i := 0; i < 10; i++ {
		result := s.PickFromPool(agents)
		t.Assert(result == "sole", "should always return sole")
	}
}

func TestPickDeterministicSeed(t *testing.T) {
	s1 := NewUserAgentPoolSampler(42)
	s2 := NewUserAgentPoolSampler(42)
	pools := map[string][]string{
		"x": []string{"a", "b"},
		"y": []string{"c", "d"},
	}
	for i := 0; i < 20; i++ {
		r1 := s1.Pick(pools, nil)
		r2 := s2.Pick(pools, nil)
		t.Assert(r1 == r2, "same seed should produce same sequence")
	}
}

func TestPickPoolNotFoundInWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(123)
	pools := map[string][]string{
		"present": []string{"ua"},
		"missing": []string{"another"},
	}
	weights := map[string]float64{
		"nonexistent": 100.0,
	}
	result := s.Pick(pools, weights)
	t.Assert(result != "", "should still pick from pools")
	// Both pools get weight 1.0 (default), so either is valid
	ok := strings.HasPrefix(result, "ua") || strings.HasPrefix(result, "another")
	t.Assert(ok, "unexpected user-agent: " + result)
}
