package timer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSystemClockNow(t *testing.T) {
	p := SystemClock{}
	a := time.Now()
	b := p.Now()
	c := time.Now()
	assert.False(t, a.After(b))
	assert.False(t, c.Before(b))
}
