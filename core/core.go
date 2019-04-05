package core

import "time"

const (
	PeriodDuration        = time.Minute * 30
	MaximumJamDuration    = time.Minute * 2
	JamTransitionDuration = time.Second * 30
	TeamTimeoutDuration   = time.Minute * 1
)

// ElapsedTime contains the current remaining time in the period and in the timeout, if applicable.
type ElapsedTime struct {
	Period  time.Duration
	Timeout time.Duration
}
