package crawler

import (
	"time"
)

// RateCalc calculates HTML extraction rates.
type RateCalc struct {
extractStarts int
extractEnds   int
startTime     time.Time
}

// NewRateCalc returns a new RateCalc instance.
func NewRateCalc() *RateCalc {
	return &RateCalc{startTime: time.Now()}
}

// OnExtractStart increments extractStarts and updates the start time.
func (r *RateCalc) OnExtractStart() {
	r.extractStarts++
	r.startTime = time.Now()
}

// OnExtractEnd increments extractEnds.
func (r *RateCalc) OnExtractEnd() {
	r.extractEnds++
}

// Rate returns the current extraction rate (extracted per second).
func (r *RateCalc) Rate() float64 {
	elapsed := time.Since(r.startTime).Seconds()
	return float64(r.extractEnds) / elapsed
}

// TotalExtracted returns the total number of extracted HTMLs.
func (r *RateCalc) TotalExtracted() int {
	return r.extractEnds
}
