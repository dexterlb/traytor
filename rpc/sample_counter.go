package rpc

import "sync"

// SampleCounter is a thread-safe  decreasing counter
type SampleCounter struct {
	sync.Mutex
	Counter int
}

// NewSampleCounter initializes the counter with a value
func NewSampleCounter(value int) *SampleCounter {
	return &SampleCounter{Counter: value + 1}
}

// Dec decreases the counter by value, but doesn't go below zero (if
// value > sc.value, the counter's value will become 0). Returns the
// new value.
func (sc *SampleCounter) Dec(value int) int {
	sc.Lock()
	defer sc.Unlock()

	if sc.Counter >= value {
		sc.Counter -= value
	} else {
		value = sc.Counter
		sc.Counter = 0
	}
	return value
}
