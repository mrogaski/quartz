package clock

import "time"

// The Timer interface defines a Ticker-style stopwatch goroutine that can be started or stopped via methods.
type Timer interface {
	// GetTickChannel returns a channel that delivers the time.Duration ticks.
	GetTickChannel() (chan time.Duration)

	// Start begins or resumes the timer, returning the current elapsed time before advancing the clock.
	Start() (time.Duration, error)

	// Stop halts the timer and returns the current elapsed time.
	Stop() (time.Duration, error)

	// Reset zeroes the elapsed time and returns the elapsed time prior to the reset.
	Reset() (time.Duration, error)

	// Close stops the timer and destroys the timer goroutine.
	Close() error
}
