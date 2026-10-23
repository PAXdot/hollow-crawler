package crawler

import (
	"strings"

	"golang.org/x/net/html"
)

// RobotsMetaTag holds the discovered Robots meta tag.
type RobotsMetaTag struct {
	NoIndex bool
	NoFollow bool
}

// ParseRobots handles the parsing of Robots meta tag.
func ParseRobots(n *html.Node) (*RobotsMetaTag, error) {
	if n.Type != html.ElementNode || n.Data != "meta" {
		return nil, nil
	}

	var content, name string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			content = c.Data
		} else if c.Type == html.ElementNode && c.Data == "name" {
			name = getAttribute(c, "content")
		}
	}

	if name != "robots" || content == "" {
		return nil, nil
	}

	tag := &RobotsMetaTag{}
	for _, part := range strings.Split(content, ",") {
		part = strings.TrimSpace(part)
		switch part {
		case "noindex":
			tag.NoIndex = true
		case "nofollow":
			tag.NoFollow = true
		}
	}

	return tag, nil
}

// getAttribute retrieves the specified attribute value from a node.
func getAttribute(n *html.Node, attr string) string {
	attrVal, ok := n.Attr[attr]
	if !ok {
		return ""
	}
	return attrVal
}
