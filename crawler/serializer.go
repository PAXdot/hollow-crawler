package crawler

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"time"
)

// Serializer provides stateless helpers for encoding and decoding
// common crawler state types (dedup keys, retry records, timestamps).
// All methods are safe for concurrent use.
type Serializer struct{}

// DedupKey represents a normalized URL and its computed hash for dedup.
type DedupKey struct {
	URL    string
	Hash64 uint64
}

// RetryRecord holds the state for a single retry attempt.
type RetryRecord struct {
	URL       string
	Attempt   int
	LastError string
	Backoff   time.Duration
	NextTry   time.Time
}

// NewDedupKey hashes the URL and returns a DedupKey.
func (s *Serializer) NewDedupKey(url string) DedupKey {
	h := fnv.New64a()
	h.Write([]byte(url))
	return DedupKey{
		URL:    url,
		Hash64: h.Sum64(),
	}
}

// EncodeDedupKey writes the DedupKey as a binary blob: 8-byte hash + length-prefixed URL.
func (s *Serializer) EncodeDedupKey(w io.Writer, dk DedupKey) error {
	if err := binary.Write(w, binary.LittleEndian, dk.Hash64); err != nil {
		return fmt.Errorf("write hash64: %w", err)
	}
	urlBytes := []byte(dk.URL)
	urlLen := uint32(len(urlBytes))
	if err := binary.Write(w, binary.LittleEndian, urlLen); err != nil {
		return fmt.Errorf("write url length: %w", err)
	}
	if _, err := w.Write(urlBytes); err != nil {
		return fmt.Errorf("write url: %w", err)
	}
	return nil
}

// DecodeDedupKey reads a DedupKey from a binary blob written by EncodeDedupKey.
func (s *Serializer) DecodeDedupKey(r io.Reader) (DedupKey, error) {
	var hash64 uint64
	if err := binary.Read(r, binary.LittleEndian, &hash64); err != nil {
		return DedupKey{}, fmt.Errorf("read hash64: %w", err)
	}
	var urlLen uint32
	if err := binary.Read(r, binary.LittleEndian, &urlLen); err != nil {
		return DedupKey{}, fmt.Errorf("read url length: %w", err)
	}
	urlBytes := make([]byte, urlLen)
	if _, err := io.ReadFull(r, urlBytes); err != nil {
		return DedupKey{}, fmt.Errorf("read url: %w", err)
	}
	return DedupKey{
		URL:    string(urlBytes),
		Hash64: hash64,
	}, nil
}

// MarshalRetryRecord serializes a RetryRecord as JSON.
func (s *Serializer) MarshalRetryRecord(rr RetryRecord) ([]byte, error) {
	data, err := json.Marshal(rr)
	if err != nil {
		return nil, fmt.Errorf("marshal retry record: %w", err)
	}
	return data, nil
}

// UnmarshalRetryRecord deserializes JSON bytes into a RetryRecord.
func (s *Serializer) UnmarshalRetryRecord(data []byte) (RetryRecord, error) {
	var rr RetryRecord
	if err := json.Unmarshal(data, &rr); err != nil {
		return RetryRecord{}, fmt.Errorf("unmarshal retry record: %w", err)
	}
	return rr, nil
}

// TimestampToBytes converts a time.Time to an 8-byte UnixNano representation.
func (s *Serializer) TimestampToBytes(t time.Time) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(t.UnixNano()))
	return buf
}

// BytesToTimestamp reads an 8-byte slice and returns the corresponding time.Time.
func (s *Serializer) BytesToTimestamp(buf []byte) (time.Time, error) {
	if len(buf) < 8 {
		return time.Time{}, fmt.Errorf("buffer too short for timestamp: %d bytes", len(buf))
	}
	nano := binary.LittleEndian.Uint64(buf[:8])
	return time.Unix(0, int64(nano)), nil
}
