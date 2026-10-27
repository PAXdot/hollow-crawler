package crawler

import (
	"math/rand"
	"sync"
)

// UserAgentPoolSampler provides weighted random selection from a set of
// user-agent pools. It is designed to work alongside UserAgentPoolManager
// without duplicating its state, by accepting an existing pool map.
type UserAgentPoolSampler struct {
	mu sync.RWMutex
	rng *rand.Rand
}

// NewUserAgentPoolSampler creates a sampler with an optional random seed.
func NewUserAgentPoolSampler(seed int64) *UserAgentPoolSampler {
	s := new(UserAgentPoolSampler)
	s.rng = rand.New(rand.NewPCG(seed))
	return s
}

// Pick returns a randomly selected user-agent string from the given pool map,
// respecting each pool's weight. The pool map is map[string][]string where
// keys are pool names and values are lists of user-agent strings.
// An optional weight map map[string]float64 assigns selection probability to
// each pool; if nil, pools are selected uniformly. Returns empty string if no
// pools are present.
func (s *UserAgentPoolSampler) Pick(
	pools map[string][]string,
	weights map[string]float64,
) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(pools) == 0 {
		return ""
	}

	// Normalize weights or default to uniform.
	poolNames := make([]string, len(pools))
	i := 0
	for k := range pools {
		poolNames[i] = k
		i++
	}
	normWeights := make([]float64, len(poolNames))
	total := 0.0
	for j, name := range poolNames {
		var w float64 = 1.0
		if weights != nil {
			if v, ok := weights[name]; ok {
				w = v
			}
		}
		if w < 0 {
			w = 0
		}
		normWeights[j] = w
		total += w
	}

	if total <= 0 {
		return ""
	}

	// Weighted random selection of pool.
	r := s.rng.Float64() * total
	cum := 0.0
	selected := ""
	for j, name := range poolNames {
		cum += normWeights[j]
		if r <= cum {
			selected = name
			break
		}
	}
	if selected == "" {
		selected = poolNames[len(poolNames)-1]
	}

	// Pick a random user-agent from selected pool.
	agents := pools[selected]
	if len(agents) == 0 {
		return ""
	}
	idx := s.rng.IntN(len(agents))
	return agents[idx]
}

// PickFromPool selects a random user-agent from a single pool list.
func (s *UserAgentPoolSampler) PickFromPool(agents []string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(agents) == 0 {
		return ""
	}
	idx := s.rng.IntN(len(agents))
	return agents[idx]
}
