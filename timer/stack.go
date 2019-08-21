package timer

import (
	"errors"
	"time"
)

// TimeStack is a fixed-size stack for time.Time values.
type TimeStack struct {
	n    int
	list []time.Time
}

// NewTimeStack returns a pointer to a TimeStack with the maximum size set.
func NewTimeStack(size int) (*TimeStack, error) {
	if size < 1 {
		return nil, errors.New("the stack size must be a positive integer")
	}
	return &TimeStack{n: size, list: make([]time.Time, 0)}, nil
}

// Size returns the current number of elements pushed to the stack.
func (s *TimeStack) Size() int {
	return len(s.list)
}

// Push inserts an element on the head of the stack or returns an error if the stack is full.
func (s *TimeStack) Push(t time.Time) error {
	if s.Size() >= s.n {
		return errors.New("stack is full")
	}
	s.list = append([]time.Time{t}, s.list...)
	return nil
}

// Pop removes an element from the head of the stack and returns it, or returns an error if the stack is empty.
func (s *TimeStack) Pop() (time.Time, error) {
	if s.Size() < 1 {
		return time.Time{}, errors.New("stack is empty")
	}
	t := s.list[0]
	s.list = s.list[1:]
	return t, nil
}

// Top returns the element currently at the head of the stack or an error if the stack is empty.
func (s *TimeStack) Top() (time.Time, error) {
	if s.Size() < 1 {
		return time.Time{}, errors.New("stack is empty")
	}
	return s.list[0], nil
}

// Clear empties the stack.
func (s *TimeStack) Clear() {
	s.list = make([]time.Time, 0)
}

// DurationStack is a fixed-size stack for time.Duration values.
type DurationStack struct {
	n    int
	list []time.Duration
}

// NewDurationStack returns a pointer to a DurationStack with the maximum size set.
func NewDurationStack(size int) (*DurationStack, error) {
	if size < 1 {
		return nil, errors.New("the stack size must be a positive integer")
	}
	return &DurationStack{n: size, list: make([]time.Duration, 0)}, nil
}

// Size returns the current number of elements pushed to the stack.
func (s *DurationStack) Size() int {
	return len(s.list)
}

// Push inserts an element on the head of the stack or returns an error if the stack is full.
func (s *DurationStack) Push(t time.Duration) error {
	if s.Size() >= s.n {
		return errors.New("stack is full")
	}
	s.list = append([]time.Duration{t}, s.list...)
	return nil
}

// Pop removes an element from the head of the stack and returns it, or returns an error if the stack is empty.
func (s *DurationStack) Pop() (time.Duration, error) {
	if s.Size() < 1 {
		return 0, errors.New("stack is empty")
	}
	t := s.list[0]
	s.list = s.list[1:]
	return t, nil
}

// Top returns the element currently at the head of the stack or an error if the stack is empty.
func (s *DurationStack) Top() (time.Duration, error) {
	if s.Size() < 1 {
		return 0, errors.New("stack is empty")
	}
	return s.list[0], nil
}

// Clear empties the stack.
func (s *DurationStack) Clear() {
	s.list = make([]time.Duration, 0)
}
