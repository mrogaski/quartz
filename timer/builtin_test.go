package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const MinimumTickCount = 4

func TestBuiltinTimer(t *testing.T) {
	timer := NewBuiltinTimer(time.Second)
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
	timer := NewBuiltinTimer(time.Second)
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
	timer := NewBuiltinTimer(time.Second)
	timer.Close()
}
