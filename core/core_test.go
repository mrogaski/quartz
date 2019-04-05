package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	now := time.Now()
	ts := Timestamp{now, time.Minute * 30, now, time.Second * 30}
	assert.Equal(t, now, ts.PeriodStart)
	assert.Equal(t, time.Minute*30, ts.PeriodRemaining)
	assert.Equal(t, now, ts.TimeoutStart)
	assert.Equal(t, time.Second*30, ts.TimeoutRemaining)
}
