package timer

import (
	"errors"
	"time"
	log "github.com/sirupsen/logrus"
)

const (
	PeriodDuration      = time.Minute * 30 // Default duration for a bout period.
	LineupDuration      = time.Second * 30 // Default time between jams.
	JamDuration         = time.Minute * 2  // Default maximum jam time.
	TeamTimeoutDuration = time.Minute * 1  // Default duration of a team timeout.
)

type PeriodClockState string

// PeriodClockState reflects the running state of the period clock.
const (
	JamPending    PeriodClockState = "pending" // JamPending indicates the lineup period before the start of a jam.
	JamActive     PeriodClockState = "active"  // JamActive indicates and active jam.
	PeriodTimeout PeriodClockState = "timeout" // PeriodTimeout indicates a timeout is in effect.
)

// PeriodClock contains the elapsed time for the current period, lineup, jam, and timeout.
type PeriodClock struct {
	State          PeriodClockState
	PeriodElapsed  time.Duration
	LineupElapsed  time.Duration
	JamElapsed     time.Duration
	TimeoutElapsed time.Duration
}

// PeriodRemaining returns the amount of time remaining in the current period.
func (c PeriodClock) PeriodRemaining() (time.Duration, error) {
	if c.PeriodElapsed == 0 {
		log.Debug("Period has not yet started")
		err := errors.New("Period has not yet started")
		return 0, err
	}
	return PeriodDuration - c.PeriodElapsed, nil
}

// JamRemaining returns the amount of time remaining in the current jam.
func (c PeriodClock) JamRemaining() (time.Duration, error) {
	if c.JamElapsed == 0 {
		log.Debug("Jam has not yet started")
		err := errors.New("Jam has not yet started")
		return 0, err
	}
	return JamDuration - c.JamElapsed, nil
}
