package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const MinimumTickCount = 4

func TestBuiltinTimer(t *testing.T) {
	source := time.NewTicker(time.Second)
	defer source.Stop()
	timer := NewBuiltinTimer(source.C)
	n := MinimumTickCount
	ticker := timer.GetTickChannel()
	for curr := range ticker {
		assert.Zero(t, curr)
		n--
		if n < 0 {
			timer.Close()
			return
		}
	}
}

func TestBuiltinTimerStart(t *testing.T) {
	source := time.NewTicker(time.Second)
	defer source.Stop()
	timer := NewBuiltinTimer(source.C)
	prev := time.Millisecond
	n := MinimumTickCount
	ticker := timer.GetTickChannel()
	timer.Start()
	for curr := range ticker {
		assert.Less(t, int64(prev), int64(curr))
		n--
		if n < 0 {
			timer.Close()
			return
		}
	}
}

func TestBuiltinTimerClose(t *testing.T) {
	source := time.NewTicker(time.Second)
	defer source.Stop()
	timer := NewBuiltinTimer(source.C)
	timer.Close()
}
