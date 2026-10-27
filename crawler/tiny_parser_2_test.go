package crawler

import (
	"testing"

	"golang.org/x/net/html"
)

func TestParseRobots(t *testing.T) {
	t.Run("Returns nil when not meta tag", func(t *testing.T) {
		_, err := ParseRobots(&html.Node{Type: html.TextNode, Data: "text"})
		if err != nil {
			t.Fatal("expected no error")
		}
	})

	t.Run("Returns nil when wrong name attribute", func(t *testing.T) {
		metaNode := &html.Node{Type: html.ElementNode, Data: "meta"}
		attrNode := &html.Node{Type: html.ElementNode, Data: "name", Attr: []html.Attribute{{Key: "content", Val: "wrong"}}}
		metaNode.AppendChild(attrNode)
		_, err := ParseRobots(metaNode)
		if err != nil {
			t.Fatal("expected no error")
		}
	})

	t.Run("Returns nil when no content attribute", func(t *testing.T) {
		metaNode := &html.Node{Type: html.ElementNode, Data: "meta"}
		attrNode := &html.Node{Type: html.ElementNode, Data: "name", Attr: []html.Attribute{{Key: "content", Val: "robots"}}}
		metaNode.AppendChild(attrNode)
		_, err := ParseRobots(metaNode)
		if err != nil {
			t.Fatal("expected no error")
		}
	})

	t.Run("Parses Robots meta tag", func(t *testing.T) {
		metaNode := &html.Node{Type: html.ElementNode, Data: "meta"}
		attrNode := &html.Node{Type: html.ElementNode, Data: "name", Attr: []html.Attribute{{Key: "content", Val: "robots"}}}
		contentNode := &html.Node{Type: html.TextNode, Data: "noindex, nofollow"}
		metaNode.AppendChild(attrNode)
		metaNode.AppendChild(contentNode)

		tag, err := ParseRobots(metaNode)
		if err != nil {
			t.Fatal("expected no error")
		}
		if tag.NoIndex != true || tag.NoFollow != true {
			t.Errorf("expected NoIndex and NoFollow to be true")
		}
	})
}
