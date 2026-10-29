package crawler

import "testing"

func TestScoreURL(t *testing.T) {
	cases := []struct {
		input  string
		wanted DedupScore
	}{
		{"", DedupScoreMin},
		{"example.com", 60},
		{"example.com/path", 55},
		{"example.com/path?param=value", 75},
		{strings.Repeat("a", 101), 100},
		{strings.Repeat("/", 11), 45},
	}
	for _, c := range cases {
		score := ScoreURL(c.input)
		if score != c.wanted {
			t.Errorf("ScoreURL(%q) got %d, want %d", c.input, score, c.wanted)
		}
	}
}

func TestIsHighScore(t *testing.T) {
	cases := []struct {
		input  DedupScore
		wanted bool
	}{
		{DedupScoreMin, false},
		{DedupScoreMax, true},
		{50, false},
		{75, true},
		{DedupScore(76), true},
	}
	for _, c := range cases {
		if IsHighScore(c.input) != c.wanted {
			t.Errorf("IsHighScore(%d) got %t, want %t", c.input, !c.wanted, c.wanted)
		}
	}
}
