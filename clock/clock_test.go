package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPeriodClockPeriodRemaining(t *testing.T) {
	c := PeriodClock{PeriodElapsed: time.Minute * 25}
	tm, err := c.PeriodRemaining()
	assert.Equal(t, PeriodDuration-time.Minute*25, tm)
	assert.Nil(t, err)
}

func TestPeriodClockPeriodRemainingError(t *testing.T) {
	c := PeriodClock{}
	tm, err := c.PeriodRemaining()
	assert.Zero(t, tm)
	assert.EqualError(t, err, "period has not yet started")
}

func TestPeriodClockLineupRemaining(t *testing.T) {
	c := PeriodClock{LineupElapsed: time.Minute * 25}
	tm, err := c.LineupRemaining()
	assert.Equal(t, LineupDuration-time.Minute*25, tm)
	assert.Nil(t, err)
}

func TestPeriodClockLineupRemainingRemainingError(t *testing.T) {
	c := PeriodClock{}
	tm, err := c.LineupRemaining()
	assert.Zero(t, tm)
	assert.EqualError(t, err, "lineup has not yet started")
}

func TestPeriodClockJamRemaining(t *testing.T) {
	c := PeriodClock{JamElapsed: time.Second * 25}
	tm, err := c.JamRemaining()
	assert.Equal(t, JamDuration-time.Second*25, tm)
	assert.Nil(t, err)
}

func TestPeriodClockJamRemainingError(t *testing.T) {
	c := PeriodClock{}
	tm, err := c.JamRemaining()
	assert.Zero(t, tm)
	assert.EqualError(t, err, "jam has not yet started")
}
