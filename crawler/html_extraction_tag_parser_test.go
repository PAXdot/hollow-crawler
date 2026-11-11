package crawler

import (
	"testing"
)

func Test_parseStartTag(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  htmlTagAttrs
	}{
		{
			name:  "empty input",
			input: "",
			want:  htmlTagAttrs{},
		},
		{
			name:  "bare tag name",
			input: "<a>",
			want: htmlTagAttrs{
				Name: "a",
			},
		},
		{
			name:  "tag with single attribute",
			input: "<a href=\"/x\">",
			want: htmlTagAttrs{
				Name: "a",
				Attrs: map[string]string{
					"href": "/x",
				},
			},
		},
		{
			name:  "tag with multiple attributes",
			input: "<a href=\"/x\" class=\"y\">",
			want: htmlTagAttrs{
				Name: "a",
				Attrs: map[string]string{
					"href":  "/x",
					"class": "y",
				},
			},
		},
		{
			name:  "tag with quoted attribute values",
			input: "<a href=\"/x\" class='y'>",
			want: htmlTagAttrs{
				Name: "a",
				Attrs: map[string]string{
					"href":  "/x",
					"class": "y",
				},
			},
		},
		{
			name:  "self-closing tag with single attribute",
			input: "<meta charset=\"utf-8\">",
			want: htmlTagAttrs{
				Name: "meta",
				Attrs: map[string]string{
					"charset": "utf-8",
				},
			},
		},
		{
			name:  "tag with leading/trailing space",
			input: "  <a href=\"/x\">  ",
			want: htmlTagAttrs{
				Name: "a",
				Attrs: map[string]string{
					"href": "/x",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseStartTag(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseStartTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
