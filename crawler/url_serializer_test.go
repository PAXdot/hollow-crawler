package crawler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLSerializer_SerializeAndDeserialize(t *testing.T) {
	serializer := &URLSerializer{}
	testURLs := []string{"https://example.com", "https://golang.org"}

	tempFile, err := ioutil.TempFile("", "test_url_serializer.*")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = serializer.Serialize(testURLs, tempFile.Name())
	assert.NoError(t, err)

	deserializedURLs, err := serializer.Deserialize(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, testURLs, deserializedURLs)

	_, err = serializer.Deserialize("non_existent_file")
	assert.ErrorIs(t, err, ErrInvalidFilePath)

	err = serializer.Serialize(nil, tempFile.Name())
	assert.NoError(t, err)
}

func TestURLSerializer_LoadURLsFromFile(t *testing.T) {
	serializer := &URLSerializer{}

	tempFile, err := ioutil.TempFile("", "test_url_serializer.*")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	err = serializer.SaveURLsToFile([]string{"https://example.com"}, tempFile.Name())
	assert.NoError(t, err)

	urls, err := serializer.LoadURLsFromFile(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, []string{"https://example.com"}, urls)

	_, err = serializer.LoadURLsFromFile("non_existent_file")
	assert.ErrorIs(t, err, ErrInvalidFilePath)
}
