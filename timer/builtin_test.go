package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuiltinTimer(t *testing.T) {
	timer := NewBuiltinTimer(time.Second)
	prev := time.Millisecond
	n := 8
	ticker := timer.GetTickChannel()
	timer.Start()
	for curr := range ticker {
		assert.Less(t, int64(prev), int64(curr))
		n--
		if n < 0 {
			return
		}
	}
}
