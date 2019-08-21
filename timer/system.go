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

	source   *time.Ticker
	time     TimeProvider
	ticker   chan time.Duration
	command  chan Command
	response chan error
}

// NewSystemTimer returns a pointer to an initialized SystemTimer, given a core time.Ticker struct and an
// implementation of the clock.RevertibleTimer interface.
func NewSystemTimer(source *time.Ticker, p TimeProvider) *SystemTimer {
	t := &SystemTimer{
		source:   source,
		time:     p,
		ticker:   make(chan time.Duration, 1),
		command:  make(chan Command, 1),
		response: make(chan error, 1),
	}
	t.C = t.ticker
	go t.run()
	return t
}

func (t *SystemTimer) run() {
	defer close(t.ticker)
	defer close(t.response)
	state := StateStart
	base := []time.Time{t.time.Now()}
	elapsed := []time.Duration{time.Second * 0}
	for {

		select {
		case cmd := <-t.command:

			switch cmd {
			case CmdStart:
				switch state {
				case StateStart, StateStopped, StateStoppedRevertible:
					base = append([]time.Time{t.time.Now()}, base...)
					state = StateRunningRevertible
					log.WithFields(log.Fields{
						"base":    base,
						"elapsed": elapsed,
					}).Debug("moving to RUNNING_REVERTIBLE")
					t.response <- nil
				default:
					t.response <- errors.New("cannot start a running timer")
				}

			case CmdStop:
				switch state {
				case StateRunning, StateRunningRevertible:
					elapsed = append([]time.Duration{t.time.Now().Sub(base[0])}, elapsed...)
					state = StateStoppedRevertible
					log.WithFields(log.Fields{
						"base":    base,
						"elapsed": elapsed,
					}).Debug("moving to STOPPED_REVERTIBLE")
					t.response <- nil
				default:
					t.response <- errors.New("cannot stop a stopped timer")
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
					t.response <- nil
				case StateStoppedRevertible:
					elapsed = elapsed[1:]
					state = StateRunning
					log.WithFields(log.Fields{
						"base":    base,
						"elapsed": elapsed,
					}).Debug("moving to RUNNING")
					t.response <- nil
				default:
					t.response <- errors.New("revert not available")
				}

			case CmdReset:
				base = []time.Time{t.time.Now()}
				elapsed = []time.Duration{time.Second * 0}
				state = StateStart
				log.WithFields(log.Fields{
					"base":    base,
					"elapsed": elapsed,
				}).Debug("reset to START")
				t.response <- nil

			case CmdClose:
				log.WithFields(log.Fields{
					"base":    base,
					"elapsed": elapsed,
				}).Debug("shutting down timer")
				t.response <- nil
				return
			}

		case tc := <-t.source.C:
			log.WithFields(log.Fields{
				"base":    base,
				"elapsed": elapsed,
				"tick":    tc,
			}).Debug("tick received")
			switch state {
			case StateRunningRevertible, StateRunning:
				t.ticker <- tc.Sub(base[0]) + elapsed[0]
			default:
				t.ticker <- elapsed[0]
			}
		}
	}
}

// Start causes the timer to begin running and recording elapsed time.
func (t *SystemTimer) Start() error {
	t.command <- CmdStart
	return <-t.response
}

// Stop halts the timer and freezes the elapsed time.
func (t *SystemTimer) Stop() error {
	t.command <- CmdStop
	return <-t.response
}

// Revert moves back to the previous running or stopped state.  Revert may only be called once after advancing to a
// new state.
func (t *SystemTimer) Revert() error {
	t.command <- CmdRevert
	return <-t.response
}

// Reset moves the timer back to the start state and clears the elapsed time.
func (t *SystemTimer) Reset() error {
	t.command <- CmdReset
	return <-t.response
}

// Close terminates the timer goroutine.
func (t *SystemTimer) Close() error {
	t.command <- CmdClose
	return <-t.response
}
