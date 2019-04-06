package core

import (
	"errors"
	"time"
)

const (
	PeriodDuration      = time.Minute * 30
	LineupDuration      = time.Second * 30
	JamDuration         = time.Minute * 2
	TeamTimeoutDuration = time.Minute * 1
)

// PeriodClock contains the elapsed time for the current period, lineup, jam, and timeout.
type PeriodClock struct {
	PeriodElapsed  time.Duration
	LineupElapsed  time.Duration
	JamElapsed     time.Duration
	TimeoutElapsed time.Duration
}

// PeriodRemaining returns the amount of time remaining in the current period.
func (c PeriodClock) PeriodRemaining() (time.Duration, error) {
	if c.PeriodElapsed == 0 {
		err := errors.New("Period has not yet started")
		return 0, err
	}
	return PeriodDuration - c.PeriodElapsed, nil
}

// JamRemaining returns the amount of time remaining in the current jam.
func (c PeriodClock) JamRemaining() (time.Duration, error) {
	if c.JamElapsed == 0 {
		err := errors.New("Jam has not yet started")
		return 0, err
	}
	return JamDuration - c.JamElapsed, nil
}
