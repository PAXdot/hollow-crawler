// url_serializer.go

package crawler

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// URLSerializer serializes and deserializes URL frontier.
type URLSerializer struct {
	mu sync.Mutex
}

func (s *URLSerializer) Serialize(urls []string, destPath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(urls)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *URLSerializer) Deserialize(srcPath string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return nil, err
	}

	var urls []string
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *URLSerializer) Load URLsFromFile(filePath string) ([]string, error) {
	_, filename := filepath.Split(filePath)
	if filename == "" {
		return nil, ErrInvalidFilePath
	}

	return s.Deserialize(filePath)
}

func (s *URLSerializer) Save URLsToFile(urls []string, filePath string) error {
	_, filename := filepath.Split(filePath)
	if filename == "" {
		return ErrInvalidFilePath
	}

	return s.Serialize(urls, filePath)
}
