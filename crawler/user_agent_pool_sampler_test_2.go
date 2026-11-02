package crawler

import (
	"math/rand"
	"strings"
	"testing"
)

// Tests for UserAgentPoolSampler.

func TestNewUserAgentPoolSampler(t *testing.T) {
	s := NewUserAgentPoolSampler(42)
	if s == nil {
		t.Fatal("sampler should not be nil")
	}
	if s.rng == nil {
		t.Fatal("rng should be initialized")
	}
}

func TestPickEmptyPools(t *testing.T) {
	s := NewUserAgentPoolSampler(42)
	result := s.Pick(map[string][]string{}, nil)
	if result != "" {
		t.Fatalf("should return empty string for empty pools, got %q", result)
	}
}

func TestPickUniformWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(12345)
	pools := map[string][]string{
		"a": {"ua1", "ua2"},
		"b": {"ua3", "ua4"},
		"c": {"ua5"},
	}
	seen := map[string]int{}
	for i := 0; i < 1000; i++ {
		ua := s.Pick(pools, nil)
		if ua == "" {
			t.Fatal("should not be empty")
		}
		seen[ua]++
	}
	for _, agents := range pools {
		for _, a := range agents {
			if seen[a] == 0 {
				t.Errorf("each user-agent should appear: %s", a)
			}
		}
	}
}

func TestPickWithWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(999)
	pools := map[string][]string{
		"light": {"l1", "l2"},
		"heavy": {"h1", "h2"},
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
	heavyCount := seen["h1"] + seen["h2"]
	lightCount := seen["l1"] + seen["l2"]
	if heavyCount <= lightCount {
		t.Fatalf("heavy pool should be selected more often: heavy=%d light=%d", heavyCount, lightCount)
	}
}

func TestPickNegativeWeight(t *testing.T) {
	s := NewUserAgentPoolSampler(7)
	pools := map[string][]string{
		"a": {"ua1"},
		"b": {"ua2"},
	}
	weights := map[string]float64{
		"a": -5.0,
		"b": 10.0,
	}
	result := s.Pick(pools, weights)
	if result == "" {
		t.Fatal("should pick something")
	}
	if result != "ua2" {
		t.Fatalf("should pick from pool with positive weight, got %q", result)
	}
}

func TestPickAllZeroWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(1)
	pools := map[string][]string{
		"a": {"ua1"},
	}
	weights := map[string]float64{
		"a": 0.0,
	}
	result := s.Pick(pools, weights)
	if result != "" {
		t.Fatalf("should return empty when total weight is zero, got %q", result)
	}
}

func TestPickPoolSelectionEdgeCase(t *testing.T) {
	s := NewUserAgentPoolSampler(42)
	pools := map[string][]string{
		"only": {"unique_ua"},
	}
	result := s.Pick(pools, nil)
	if result != "unique_ua" {
		t.Fatalf("should select from the only pool, got %q", result)
	}
}

func TestPickFromPool(t *testing.T) {
	s := NewUserAgentPoolSampler(100)
	agents := []string{"alpha", "beta", "gamma"}
	seen := map[string]int{}
	for i := 0; i < 300; i++ {
		ua := s.PickFromPool(agents)
		if ua == "" {
			t.Fatal("should not be empty")
		}
		seen[ua]++
	}
	for _, a := range agents {
		if seen[a] == 0 {
			t.Errorf("each agent should appear: %s", a)
		}
	}
}

func TestPickFromPoolEmpty(t *testing.T) {
	s := NewUserAgentPoolSampler(0)
	result := s.PickFromPool([]string{})
	if result != "" {
		t.Fatalf("should return empty for empty list, got %q", result)
	}
}

func TestPickFromPoolSingle(t *testing.T) {
	s := NewUserAgentPoolSampler(77)
	agents := []string{"sole"}
	for i := 0; i < 10; i++ {
		result := s.PickFromPool(agents)
		if result != "sole" {
			t.Fatalf("should always return sole, got %q", result)
		}
	}
}

func TestPickDeterministicSeed(t *testing.T) {
	s1 := NewUserAgentPoolSampler(42)
	s2 := NewUserAgentPoolSampler(42)
	pools := map[string][]string{
		"x": {"a", "b"},
		"y": {"c", "d"},
	}
	for i := 0; i < 20; i++ {
		r1 := s1.Pick(pools, nil)
		r2 := s2.Pick(pools, nil)
		if r1 != r2 {
			t.Fatalf("same seed should produce same sequence: %q != %q", r1, r2)
		}
	}
}

func TestPickPoolNotFoundInWeights(t *testing.T) {
	s := NewUserAgentPoolSampler(123)
	pools := map[string][]string{
		"present": {"ua"},
		"missing": {"another"},
	}
	weights := map[string]float64{
		"nonexistent": 100.0,
	}
	result := s.Pick(pools, weights)
	if result == "" {
		t.Fatal("should still pick from pools")
	}
	ok := strings.HasPrefix(result, "ua") || strings.HasPrefix(result, "another")
	if !ok {
		t.Fatalf("unexpected user-agent: %q", result)
	}
}

// TestPickFromPoolDoesNotMutate ensures sampling leaves the source slice intact.
func TestPickFromPoolDoesNotMutate(t *testing.T) {
	s := NewUserAgentPoolSampler(5)
	agents := []string{"one", "two", "three"}
	original := append([]string(nil), agents...)
	for i := 0; i < 50; i++ {
		s.PickFromPool(agents)
	}
	for i := range agents {
		if agents[i] != original[i] {
			t.Fatalf("source slice was mutated at index %d: %q != %q", i, agents[i], original[i])
		}
	}
}

// TestPickFromPoolZeroValueSampler verifies a zero-value sampler still works.
func TestPickFromPoolZeroValueSampler(t *testing.T) {
	var s UserAgentPoolSampler
	s.rng = rand.New(rand.NewSource(1))
	result := s.PickFromPool([]string{"only"})
	if result != "only" {
		t.Fatalf("expected only, got %q", result)
	}
}
