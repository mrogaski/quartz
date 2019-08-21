package timer

import "time"

// TimeProvider is an interface for mockable wrappers to Go's time builtin.
type TimeProvider interface {
	Now() time.Time
}

// SystemClock is a simple wrapper for time.Now().
type SystemClock struct{}

// Now returns a time.Time value for the current time according to the system clock.
func (SystemClock) Now() time.Time {
	return time.Now()
}
