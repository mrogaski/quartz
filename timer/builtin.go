package timer

import "time"

type BuiltinTimer struct {
	ticker  <-chan time.Duration // The channel on which the ticks are delivered.
	command chan<- string        // The command channel to control the timer.
}

func NewBuiltinTimer(d time.Duration) *BuiltinTimer {
	ticker := make(chan time.Duration)
	command := make(chan string)
	running := false
	go func() {
		var offset time.Time
		var curr time.Duration
		source := time.NewTicker(d)
		defer source.Stop()
		for {
			select {
			case cmd := <-command:
				switch cmd {
				case "start":
					offset = time.Now()
					running = true
				case "close":
					return
				}
			case t := <-source.C:
				if running {
					curr = t.Sub(offset)
				}
				ticker <- curr
			}
		}
	}()
	return &BuiltinTimer{ticker: ticker, command: command}
}

func (timer *BuiltinTimer) GetTickChannel() <-chan time.Duration {
	return timer.ticker
}

func (timer *BuiltinTimer) Start() {
	timer.command <- "start"
}

func (timer *BuiltinTimer) Close() {
	timer.command <- "close"
}
