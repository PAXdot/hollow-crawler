package crawler

import (
	"testing"
	"time"
)

func TestRateCalc_Record(t *testing.T) {
	rc := NewRateCalcWithSize(5 * time.Second)
	now := time.Now()

	rc.Record(10)
	rc.Record(20)
	rc.Record(15)

	time.Sleep(1 * time.Second)
	rc.Record(5)
	time.Sleep(1 * time.Second)
	rc.Record(7)
	time.Sleep(1 * time.Second)
	rc.Record(12)

	expected := float64((10+20+15+5+7+12)) / 5
	throughput := rc.Throughput()

	if throughput != expected {
		t.Errorf("Throughput() = %v; want %v", throughput, expected)
	}
}

func TestRateCalc_Reset(t *testing.T) {
	rc := NewRateCalcWithSize(5 * time.Second)
	rc.Record(10)
	rc.Record(20)
	rc.Record(15)

	rc.Reset()

	throughput := rc.Throughput()
	if throughput != 0 {
		t.Errorf("Throughput() after Reset() = %v; want 0", throughput)
	}
}

func TestRateCalc_NormalizeEntities(t *testing.T) {
	rc := NewRateCalc()

	normalized, err := rc.NormalizeEntities("http://example.com/url%2526param=value")
	if err != nil {
		t.Fatalf("NormalizeEntities() returned error: %v", err)
	}

	expected := "http://example.com/url%26param=value"
	if normalized != expected {
		t.Errorf("NormalizeEntities() = %v; want %v", normalized, expected)
	}
}
