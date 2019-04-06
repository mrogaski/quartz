package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPeriodRemaining(t *testing.T) {
	c := PeriodClock{PeriodElapsed: time.Minute * 25}
	tm, err := c.PeriodRemaining()
	assert.Equal(t, PeriodDuration-time.Minute*25, tm)
	assert.Nil(t, err)
}

func TestPeriodRemainingError(t *testing.T) {
	c := PeriodClock{}
	tm, err := c.PeriodRemaining()
	assert.Zero(t, tm)
	assert.Equal(t, "Period has not yet started", err.Error())
}

func TestJamRemaining(t *testing.T) {
	c := PeriodClock{JamElapsed: time.Second * 25}
	tm, err := c.JamRemaining()
	assert.Equal(t, JamDuration-time.Second*25, tm)
	assert.Nil(t, err)
}

func TestJamRemainingError(t *testing.T) {
	c := PeriodClock{}
	tm, err := c.JamRemaining()
	assert.Zero(t, tm)
	assert.Equal(t, "Jam has not yet started", err.Error())
}
