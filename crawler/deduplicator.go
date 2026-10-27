// deduplicator.go

package crawler

import (
	"encoding/json"
	"fmt"
)

type Page struct {
	URL      string          `json:"url"`
	Content  []byte          `json:"content"`
	Hash     [32]byte        `json:"-"`
	Metadata map[string]interface{} `json:"metadata"`
}

func SerializePage(page *Page) ([]byte, error) {
	return json.MarshalIndent(page, "", "  ")
}

func DeserializePage(data []byte) (*Page, error) {
	var page Page
	err := json.Unmarshal(data, &page)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize page: %w", err)
	}
	return &page, nil
}
