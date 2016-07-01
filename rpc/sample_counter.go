package rpc

import "sync/atomic"

// SampleCounter is a thread-safe  decreasing counter
type SampleCounter struct {
	Counter int64
}

// NewSampleCounter initializes the counter with a value
func NewSampleCounter(value int) *SampleCounter {
	return &SampleCounter{Counter: int64(value)}
}

// Dec decreases the counter by value, but doesn't go below zero (if
// value > sc.value, the counter's value will become 0). Returns the
// actual amount decremented (so the return value is 0 if the counter
// has already been 0)
func (sc *SampleCounter) Dec(value int) int {
	difference := int64(value)
	newValue := atomic.AddInt64(&sc.Counter, -difference)

	if newValue >= 0 {
		return int(difference)
	}

	atomic.StoreInt64(&sc.Counter, 0)

	if -newValue < difference {
		return int(difference + newValue)
	} else {
		return 0
	}
}
