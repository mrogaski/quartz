package core

import "time"

type Timestamp struct {
	PeriodStart      time.Time
	PeriodRemaining  time.Duration
	TimeoutStart     time.Time
	TimeoutRemaining time.Duration
}

