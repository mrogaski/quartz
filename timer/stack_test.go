package timer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TimeStackTestSuite struct {
	suite.Suite
}

func TestTimeStackTestSuite(t *testing.T) {
	suite.Run(t, new(TimeStackTestSuite))
}

func (ts *TimeStackTestSuite) TestNewTimeStack() {
	s, err := NewTimeStack(2)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 0, s.Size())
}

func (ts *TimeStackTestSuite) TestNewTimeStackError() {
	s, err := NewTimeStack(0)
	assert.EqualError(ts.T(), err, "the stack size must be a positive integer")
	assert.Nil(ts.T(), s)
}

func (ts *TimeStackTestSuite) TestTimeStackPush() {
	s, err := NewTimeStack(2)
	assert.NoError(ts.T(), err)

	t := time.Now()
	err = s.Push(t)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 1, s.Size())

	err = s.Push(t.Add(time.Second * 1))
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 2, s.Size())

	err = s.Push(t.Add(time.Second * 1))
	assert.EqualError(ts.T(), err, "stack is full")
	assert.Equal(ts.T(), 2, s.Size())
}

func (ts *TimeStackTestSuite) TestTimeStackPop() {
	s, err := NewTimeStack(2)
	assert.NoError(ts.T(), err)

	t := time.Now()
	err = s.Push(t)
	assert.NoError(ts.T(), err)
	err = s.Push(t.Add(time.Second * 1))
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 2, s.Size())

	e, err := s.Pop()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), e, t.Add(time.Second*1))
	assert.Equal(ts.T(), 1, s.Size())

	e, err = s.Pop()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), e, t)
	assert.Equal(ts.T(), 0, s.Size())

	e, err = s.Pop()
	assert.EqualError(ts.T(), err, "stack is empty")
	assert.Zero(ts.T(), e)
	assert.Equal(ts.T(), 0, s.Size())
}

func (ts *TimeStackTestSuite) TestTimeStackTop() {
	s, err := NewTimeStack(2)
	assert.NoError(ts.T(), err)

	e, err := s.Top()
	assert.EqualError(ts.T(), err, "stack is empty")
	assert.Zero(ts.T(), e)

	t := time.Now()
	err = s.Push(t)
	assert.NoError(ts.T(), err)
	e, err = s.Top()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), t, e)

	err = s.Push(t.Add(time.Second * 1))
	assert.NoError(ts.T(), err)
	e, err = s.Top()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), t.Add(time.Second*1), e)

	assert.Equal(ts.T(), 2, s.Size())
}

func (ts *TimeStackTestSuite) TestTimeStackClear() {
	s, err := NewTimeStack(2)
	assert.NoError(ts.T(), err)

	t := time.Now()
	err = s.Push(t)
	assert.NoError(ts.T(), err)
	err = s.Push(t.Add(time.Second * 1))
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 2, s.Size())

	s.Clear()
	assert.Equal(ts.T(), 0, s.Size())
}

type DurationStackTestSuite struct {
	suite.Suite
}

func TestDurationStackTestSuite(t *testing.T) {
	suite.Run(t, new(DurationStackTestSuite))
}

func (ts *DurationStackTestSuite) TestNewDurationStack() {
	s, err := NewDurationStack(2)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 0, s.Size())
}

func (ts *DurationStackTestSuite) TestNewDurationStackError() {
	s, err := NewDurationStack(0)
	assert.EqualError(ts.T(), err, "the stack size must be a positive integer")
	assert.Zero(ts.T(), s)
}

func (ts *DurationStackTestSuite) TestDurationStackPush() {
	s, err := NewDurationStack(2)
	assert.NoError(ts.T(), err)

	err = s.Push(time.Second * 1)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 1, s.Size())

	err = s.Push(time.Second * 2)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 2, s.Size())

	err = s.Push(time.Second * 3)
	assert.EqualError(ts.T(), err, "stack is full")
	assert.Equal(ts.T(), 2, s.Size())
}

func (ts *DurationStackTestSuite) TestDurationStackPop() {
	s, err := NewDurationStack(2)
	assert.NoError(ts.T(), err)

	err = s.Push(time.Second * 1)
	assert.NoError(ts.T(), err)
	err = s.Push(time.Second * 2)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 2, s.Size())

	e, err := s.Pop()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), e, time.Second*2)
	assert.Equal(ts.T(), 1, s.Size())

	e, err = s.Pop()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), e, time.Second*1)
	assert.Equal(ts.T(), 0, s.Size())

	e, err = s.Pop()
	assert.EqualError(ts.T(), err, "stack is empty")
	assert.Zero(ts.T(), e)
	assert.Equal(ts.T(), 0, s.Size())
}

func (ts *DurationStackTestSuite) TestDurationStackTop() {
	s, err := NewDurationStack(2)
	assert.NoError(ts.T(), err)

	e, err := s.Top()
	assert.EqualError(ts.T(), err, "stack is empty")
	assert.Zero(ts.T(), e)

	err = s.Push(time.Second * 1)
	assert.NoError(ts.T(), err)
	e, err = s.Top()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), time.Second*1, e)

	err = s.Push(time.Second * 2)
	assert.NoError(ts.T(), err)
	e, err = s.Top()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), time.Second*2, e)

	assert.Equal(ts.T(), 2, s.Size())
}

func (ts *DurationStackTestSuite) TestDurationStackClear() {
	s, err := NewDurationStack(2)
	assert.NoError(ts.T(), err)

	err = s.Push(time.Second * 1)
	assert.NoError(ts.T(), err)
	err = s.Push(time.Second * 2)
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), 2, s.Size())

	s.Clear()
	assert.Equal(ts.T(), 0, s.Size())
}
