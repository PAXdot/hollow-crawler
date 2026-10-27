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
	testing.Assert(len(lines) >= 4, "expected at least 4 lines")
	testing.Assert(strings.Contains(lines[0], "Crawling completed"))
	testing.Assert(strings.Contains(lines[1], "Total URLs processed: 0"))
	testing.Assert(strings.Contains(lines[2], "Successful scrapes: 0"))
	testing.Assert(strings.Contains(lines[3], "Failed scrapes: 0"))
	f.StartTime = oldStart
}

func TestPrintSummaryNonZeroValues() {
	f := NewFormatter()
	f.StartTime = time.Unix(2000000, 0)
	out := captureStdout(func() {
		f.PrintSummary(100, 80, 20)
	})
	testing.Assert(strings.Contains(out, "Total URLs processed: 100"))
	testing.Assert(strings.Contains(out, "Successful scrapes: 80"))
	testing.Assert(strings.Contains(out, "Failed scrapes: 20"))
}
