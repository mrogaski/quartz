package clock

import "time"

// The Timer interface defines a Ticker-style stopwatch goroutine that can be started or stopped via methods.
type Timer interface {
	// GetTickChannel returns a channel that delivers the time.Duration ticks.
	GetTickChannel() <-chan time.Duration

	// Start begins or resumes the timer.
	Start()

	// Stop halts the timer.
	Stop()

	// Reset zeroes the elapsed time.
	Reset()

	// Close stops the timer and destroys the timer goroutine.
	Close()
}
