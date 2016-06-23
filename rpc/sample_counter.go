package rpc

import "sync"

type SampleCounter struct {
	sync.Mutex
	Counter int
}

func NewSampleCounter(value int) *SampleCounter {
	return &SampleCounter{Counter: value}
}

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
