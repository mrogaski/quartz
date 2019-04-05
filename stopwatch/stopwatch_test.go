package stopwatch

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStopWatch(t *testing.T) {
	sw := NewStopWatch(time.Second)
	prev := time.Millisecond
	n := 8
	for curr := range sw.C {
		assert.Less(t, int64(prev), int64(curr))
		n--
		if n < 0 {
			return
		}
	}
}
