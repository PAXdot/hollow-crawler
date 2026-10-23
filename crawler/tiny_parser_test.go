package crawler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractLinks(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		links  []string
	}{
		{
			name: "simple case",
			input: "<a href=/one>One</a><a href=/two>Two</a>",
			links: []string{"/one", "/two"},
		},
		{
			name: "empty string",
			input: "",
			links: []string{},
		},
		{
			name: "no anchor tags",
			input: "No anchor tags here.",
			links: []string{},
		},
		{
			name: "invalid href attributes",
			input: "<a href>Invalid</a>",
			links: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			links := ExtractLinks(tc.input)
			assert.Equal(t, tc.links, links)
		})
	}
}
