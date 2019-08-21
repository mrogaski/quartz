package timer

import (
	"errors"
	"time"
)

type Command uint8

const (
	CmdStart Command = iota
	CmdStop
	CmdRevert
	CmdReset
	CmdClose
)

type State uint8

const (
	StateStart State = iota
	StateRunningRevertible
	StateRunning
	StateStoppedRevertible
	StateStopped
)

type SystemTimer struct {
	ticker   <-chan time.Duration // The channel on which the ticks are delivered.
	command  chan<- Command       // The command channel to control the timer.
	response <-chan error         // A response channel, reports an error if the command failed or nil otherwise.
}

func NewSystemTimer(c <-chan time.Time, p TimeProvider) *SystemTimer {
	tc := make(chan time.Duration, 1)
	cc := make(chan Command, 1)
	rc := make(chan error, 1)
	go func(source <-chan time.Time, p TimeProvider) {
		state := StateStart
		base := []time.Time{p.Now()}
		elapsed := []time.Duration{time.Second * 0}
		for {

			select {
			case cmd := <-cc:

				switch cmd {
				case CmdStart:
					switch state {
					case StateStart, StateStopped, StateStoppedRevertible:
						base = append([]time.Time{p.Now()}, base...)
						state = StateRunningRevertible
						rc <- nil
					default:
						rc <- errors.New("cannot start a running timer")
					}

				case CmdStop:
					switch state {
					case StateRunning, StateRunningRevertible:
						elapsed = append([]time.Duration{p.Now().Sub(base[0])}, elapsed...)
						state = StateStoppedRevertible
						rc <- nil
					default:
						rc <- errors.New("cannot stop a stopped timer")
					}

				case CmdRevert:
					switch state {
					case StateRunningRevertible:
						elapsed = elapsed[:1]
						state = StateStopped
						rc <- nil
					case StateStoppedRevertible:
						base = base[:1]
						state = StateRunning
						rc <- nil
					default:
						rc <- errors.New("revert not available")
					}

				case CmdReset:
					base = []time.Time{p.Now()}
					elapsed = []time.Duration{time.Second * 0}

				case CmdClose:
					return
				}

			case t := <-source:
				switch state {
				case StateRunningRevertible, StateRunning:
					tc <- t.Sub(base[0])
				default:
					tc <- elapsed[0]
				}
			}
		}
	}(c, p)
	return &SystemTimer{ticker: tc, command: cc, response: rc}
}

func (timer *SystemTimer) TickChannel() <-chan time.Duration {
	return timer.ticker
}

func (timer *SystemTimer) Start() error {
	timer.command <- CmdStart
	return <-timer.response
}

func (timer *SystemTimer) Stop() error {
	timer.command <- CmdStop
	return <-timer.response
}

func (timer *SystemTimer) Revert() error {
	timer.command <- CmdRevert
	return <-timer.response
}

func (timer *SystemTimer) Reset() error {
	timer.command <- CmdReset
	return <-timer.response
}

func (timer *SystemTimer) Close() error {
	timer.command <- CmdClose
	return <-timer.response
}
