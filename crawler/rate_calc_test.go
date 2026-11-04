package crawler

import (
	"net/url"
	"testing"
	"time"
)

func urlToEntity(parsedURL *url.URL) string {
	return &parsedURL.URL.String()
}

func TestRateCalc_Record(t *testing.T) {
	rc := NewRateCalcWithSize(6)
	for i := 0; i < 5; i++ {
		rc.Record(10)
		time.Sleep(1 * time.Second)
	}
	if rc.count != 50 {
		t.Fatalf("expected 50, got %d", rc.count)
	}
	if rc.Throughput() != 50 {
		t.Fatalf("expected 50, got %.2f", rc.Throughput())
	}
}

func TestRateCalc_Throughput_Empty(t *testing.T) {
	rc := NewRateCalcWithSize(6)
	if rc.Throughput() != 0 {
		t.Fatalf("expected 0, got %.2f", rc.Throughput())
	}
}

func TestRateCalc_NormalizeEntities(t *testing.T) {
	rc := NewRateCalc()
	normalized, err := rc.NormalizeEntities("http://www.example.com/url%26encoded%26and%26spaces")
	if err != nil {
		t.Fatal(err)
	}
	expected := "http://www.example.com/url&encoded&and&spaces"
	if normalized != expected {
		t.Fatalf("expected %s, got %s", expected, normalized)
	}
}
