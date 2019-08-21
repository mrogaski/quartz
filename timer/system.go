package timer

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

// Command is a type for command constants.
type Command uint8

// Command constants used to control the running state of the timer.
const (
	CmdStart Command = iota
	CmdStop
	CmdRevert
	CmdReset
	CmdClose
)

// State is a type for state constants.
type State uint8

// State constants used to represent the states of the timer's state machine.
const (
	StateStart State = iota
	StateRunningRevertible
	StateRunning
	StateStoppedRevertible
	StateStopped
)

// SystemTimer implements the clock.RevertibleTimer interface.
type SystemTimer struct {
	// C is the channel on which the timer elapsed time ticks are delivered.
	C <-chan time.Duration

	// command is the command channel to control the timer.
	command chan<- Command

	// response is a response channel, reports an error if the command failed or nil otherwise.
	response <-chan error
}

// NewSystemTimer returns a pointer to an initialized SystemTimer, given a core time.Ticker struct and an
// implementation of the clock.RevertibleTimer interface.
func NewSystemTimer(c <-chan time.Time, p TimeProvider) *SystemTimer {
	tc := make(chan time.Duration, 1)
	cc := make(chan Command, 1)
	rc := make(chan error, 1)
	go func(source <-chan time.Time, p TimeProvider) {
		defer close(tc)
		defer close(rc)
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
						log.WithFields(log.Fields{
							"base":    base,
							"elapsed": elapsed,
						}).Debug("moving to RUNNING_REVERTIBLE")
						rc <- nil
					default:
						rc <- errors.New("cannot start a running timer")
					}

				case CmdStop:
					switch state {
					case StateRunning, StateRunningRevertible:
						elapsed = append([]time.Duration{p.Now().Sub(base[0])}, elapsed...)
						state = StateStoppedRevertible
						log.WithFields(log.Fields{
							"base":    base,
							"elapsed": elapsed,
						}).Debug("moving to STOPPED_REVERTIBLE")
						rc <- nil
					default:
						rc <- errors.New("cannot stop a stopped timer")
					}

				case CmdRevert:
					switch state {
					case StateRunningRevertible:
						base = base[1:]
						state = StateStopped
						log.WithFields(log.Fields{
							"base":    base,
							"elapsed": elapsed,
						}).Debug("moving to STOPPED")
						rc <- nil
					case StateStoppedRevertible:
						elapsed = elapsed[1:]
						state = StateRunning
						log.WithFields(log.Fields{
							"base":    base,
							"elapsed": elapsed,
						}).Debug("moving to RUNNING")
						rc <- nil
					default:
						rc <- errors.New("revert not available")
					}

				case CmdReset:
					base = []time.Time{p.Now()}
					elapsed = []time.Duration{time.Second * 0}
					state = StateStart
					log.WithFields(log.Fields{
						"base":    base,
						"elapsed": elapsed,
					}).Debug("reset to START")
					rc <- nil

				case CmdClose:
					log.WithFields(log.Fields{
						"base":    base,
						"elapsed": elapsed,
					}).Debug("shutting down timer")
					rc <- nil
					return
				}

			case t := <-source:
				log.WithFields(log.Fields{
					"base":    base,
					"elapsed": elapsed,
					"tick":    t,
				}).Debug("tick received")
				switch state {
				case StateRunningRevertible, StateRunning:
					tc <- t.Sub(base[0]) + elapsed[0]
				default:
					tc <- elapsed[0]
				}
			}
		}
	}(c, p)
	return &SystemTimer{C: tc, command: cc, response: rc}
}

// Start causes the timer to begin running and recording elapsed time.
func (timer *SystemTimer) Start() error {
	timer.command <- CmdStart
	return <-timer.response
}

// Stop halts the timer and freezes the elapsed time.
func (timer *SystemTimer) Stop() error {
	timer.command <- CmdStop
	return <-timer.response
}

// Revert moves back to the previous running or stopped state.  Revert may only be called once after advancing to a
// new state.
func (timer *SystemTimer) Revert() error {
	timer.command <- CmdRevert
	return <-timer.response
}

// Reset moves the timer back to the start state and clears the elapsed time.
func (timer *SystemTimer) Reset() error {
	timer.command <- CmdReset
	return <-timer.response
}

// Close terminates the timer goroutine.
func (timer *SystemTimer) Close() error {
	timer.command <- CmdClose
	return <-timer.response
}
