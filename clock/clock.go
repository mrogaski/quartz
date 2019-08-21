package clock

import (
	"errors"
	"time"
)

const (
	// TickDuration is the default ticker interval.
	TickDuration = time.Second

	// PeriodDuration is the efault duration for a bout period.
	PeriodDuration = time.Minute * 30

	// LineupDuration is the default time between jams.
	LineupDuration = time.Second * 30

	// JamDuration is the efault maximum jam time.
	JamDuration = time.Minute * 2

	// TeamTimeoutDuration is the default duration of a team timeout.
	TeamTimeoutDuration = time.Minute
)

// PeriodClockState reflects the running state of the period clock.
const (
	// PeriodStart is the initial state prior to the period.
	PeriodStart = "start"

	// JamPending indicates the lineup period before the start of a jam.
	JamPending = "pending"

	// JamActive indicates and active jam.
	JamActive = "active"

	// PeriodTimeout indicates a timeout is in effect.
	PeriodTimeout = "timeout"

	// PeriodEnd is the end of normal play.
	PeriodEnd = "end"
)

// PeriodClock contains the elapsed time for the current period, lineup, jam, and timeout.
type PeriodClock struct {
	// State is the current state of the bout period.
	State          string
	PeriodElapsed  time.Duration
	LineupElapsed  time.Duration
	JamElapsed     time.Duration
	TimeoutElapsed time.Duration
}

// PeriodRemaining returns the amount of time remaining in the current period.
func (c PeriodClock) PeriodRemaining() (time.Duration, error) {
	if c.PeriodElapsed == 0 {
		err := errors.New("period has not yet started")
		return 0, err
	}
	return PeriodDuration - c.PeriodElapsed, nil
}

// LineupRemaining returns the amount of time remaining in the current period.
func (c PeriodClock) LineupRemaining() (time.Duration, error) {
	if c.LineupElapsed == 0 {
		err := errors.New("lineup has not yet started")
		return 0, err
	}
	return LineupDuration - c.LineupElapsed, nil
}

// JamRemaining returns the amount of time remaining in the current jam.
func (c PeriodClock) JamRemaining() (time.Duration, error) {
	if c.JamElapsed == 0 {
		err := errors.New("jam has not yet started")
		return 0, err
	}
	return JamDuration - c.JamElapsed, nil
}
