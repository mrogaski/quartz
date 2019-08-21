package clock

import "time"

// The Timer interface defines a Ticker-style stopwatch goroutine that can be started or stopped via methods.
type Timer interface {
	// GetTickChannel returns a channel that delivers the time.Duration ticks.
	TickChannel() <-chan time.Duration

	// Start begins or resumes the timer.
	Start() error

	// Stop halts the timer.
	Stop() error

	// Revert returns to the previous running or stopped state.
	Revert() error

	// Reset zeroes the elapsed time.
	Reset() error

	// Close stops the timer and destroys the timer goroutine.
	Close() error
}
