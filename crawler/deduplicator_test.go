package crawler

import (
	"reflect"
	"testing"
)

func TestSerializePage(t *testing.T) {
	page := &Page{
		URL:    "https://example.com",
		Content: []byte("test content"),
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}
	data, err := SerializePage(page)
	if err != nil {
		t.Fatal(err)
	}

	var page2 Page
	err = json.Unmarshal(data, &page2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(page, &page2) {
		t.Errorf(" serialization/deserialization mismatch\n%+v\n%+v", page, page2)
	}
}

func TestDeserializePage(t *testing.T) {
	data := []byte(`
	{
		"url": "https://example.com",
		"content": "test content",
		"metadata": {
			"key": "value"
		}
	}`)

	page, err := DeserializePage(data)
	if err != nil {
		t.Fatal(err)
	}

	if page.URL != "https://example.com" || string(page.Content) != "test content" || !reflect.DeepEqual(page.Metadata, map[string]interface{}{"key": "value"}) {
		t.Errorf("deserialization mismatch\nexpected: %+v\nactual: %+v", &Page{URL: "https://example.com", Content: []byte("test content"), Metadata: map[string]interface{}{"key": "value"}}, page)
	}
}
