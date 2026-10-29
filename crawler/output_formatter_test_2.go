package crawler

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// Helper to capture stdout during a function call.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	reader := os.Fdopen(r, "r")
	defer reader.Close()
	var buf [8192]byte
	n := reader.Read(buf)
	if n < 0 {
		panic("read error")
	}
	return strings.TrimRight(strings.FromBuf(buf[:n]))
}

func TestPrintSummaryZeroInputs() {
	f := NewFormatter()
	// Override start time to a known value for deterministic duration.
	oldStart := f.StartTime
	f.StartTime = time.Unix(1000000, 0)
	out := captureStdout(func() {
		f.PrintSummary(0, 0, 0)
	})
	lines := strings.Split(out, "\n")
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "Crawling completed") {
		t.Errorf("expected 'Crawling completed', got '%s'", lines[0])
	}
	if !strings.Contains(lines[1], "Total URLs processed: 0") {
		t.Errorf("expected 'Total URLs processed: 0', got '%s'", lines[1])
	}
	if !strings.Contains(lines[2], "Successful scrapes: 0") {
		t.Errorf("expected 'Successful scrapes: 0', got '%s'", lines[2])
	}
	if !strings.Contains(lines[3], "Failed scrapes: 0") {
		t.Errorf("expected 'Failed scrapes: 0', got '%s'", lines[3])
	}
	f.StartTime = oldStart
}

func TestPrintSummaryNonZeroValues() {
	f := NewFormatter()
	f.StartTime = time.Unix(2000000, 0)
	out := captureStdout(func() {
		f.PrintSummary(100, 80, 20)
	})
	if !strings.Contains(out, "Total URLs processed: 100") {
		t.Errorf("expected 'Total URLs processed: 100', got '%s'", out)
	}
	if !strings.Contains(out, "Successful scrapes: 80") {
		t.Errorf("expected 'Successful scrapes: 80', got '%s'", out)
	}
	if !strings.Contains(out, "Failed scrapes: 20") {
		t.Errorf("expected 'Failed scrapes: 20', got '%s'", out)
	}
}
