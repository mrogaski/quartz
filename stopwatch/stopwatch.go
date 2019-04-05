package stopwatch

import (
	"time"
)

// A StopWatch holds a channel that delivers a time.Duration tick at a specified interval.  Each tick
// indicates the time elapsed since the StopWatch was created.
type StopWatch struct {
	C    <-chan time.Duration // The channel on which the ticks are delivered.
	done chan<- bool
}

func NewStopWatch(d time.Duration) *StopWatch {
	c := make(chan time.Duration)
	done := make(chan bool)
	go func() {
		start := time.Now()
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				c <- t.Sub(start)
			}
		}
	}()
	return &StopWatch{C: c, done: done}
}

func (s *StopWatch) Stop() {
	s.done <- true
}
