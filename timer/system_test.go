package timer

import (
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockTimeProvider struct {
	mock.Mock
	mu  sync.Mutex
	now time.Time
}

func (m *mockTimeProvider) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}

func (m *mockTimeProvider) SetNow(t time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = t
}

type SystemTimerTestSuite struct {
	suite.Suite
	time   *mockTimeProvider
	ticker chan time.Time
}

func (ts *SystemTimerTestSuite) SetupTest() {
	ts.time = new(mockTimeProvider)
	ts.ticker = make(chan time.Time, 1)
}

func TestSystemTimerTestSuite(t *testing.T) {
	suite.Run(t, new(SystemTimerTestSuite))
}

func (ts *SystemTimerTestSuite) TestNewSystemTimer() {
	st := NewSystemTimer(ts.ticker, ts.time)
	assert.NotNil(ts.T(), st)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStart() {
	t := time.Now()
	ts.time.now = t

	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.ticker <- t
	assert.Equal(ts.T(), time.Second*0, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerRunningRevertible() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	ts.ticker <- t2
	assert.Equal(ts.T(), time.Second*1, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerRunRevertibleToStopped() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	ts.ticker <- t2
	assert.Equal(ts.T(), time.Second*0, <-c)

}

func (ts *SystemTimerTestSuite) TestSystemTimerRevertStopped() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Revert()
	assert.EqualError(ts.T(), err, "revert not available")
	ts.ticker <- t3
	assert.Equal(ts.T(), time.Second*0, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStoppedToRunningRevert() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)
	err = st.Revert()
	assert.NoError(ts.T(), err)
	err = st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	ts.ticker <- t2
	assert.Equal(ts.T(), time.Second*1, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerRunRevertibleToStoppedRevertible() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	ts.ticker <- t3
	assert.Equal(ts.T(), time.Second*1, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStoppedRevertibleToRunningRevertible() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)
	t4 := t0.Add(time.Second * 4)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t4)
	ts.ticker <- t4
	assert.Equal(ts.T(), time.Second*2, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStoppedRevertibleToRunning() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)
	t4 := t0.Add(time.Second * 4)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t4)
	ts.ticker <- t4
	assert.Equal(ts.T(), time.Second*3, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerRevertRunning() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)
	t4 := t0.Add(time.Second * 4)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t4)
	err = st.Revert()
	assert.EqualError(ts.T(), err, "revert not available")
	ts.ticker <- t4
	assert.Equal(ts.T(), time.Second*3, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerRunningToStoppedRevertible() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)
	t4 := t0.Add(time.Second * 4)
	t5 := t0.Add(time.Second * 5)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t4)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t5)
	ts.ticker <- t5
	assert.Equal(ts.T(), time.Second*3, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStartRunningRevertible() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	ts.ticker <- t2
	assert.Equal(ts.T(), time.Second*1, <-c)

	err = st.Start()
	assert.EqualError(ts.T(), err, "cannot start a running timer")
	ts.ticker <- t2
	assert.Equal(ts.T(), time.Second*1, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStartRunning() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)
	t4 := t0.Add(time.Second * 4)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t4)
	err = st.Start()
	assert.EqualError(ts.T(), err, "cannot start a running timer")
	ts.ticker <- t4
	assert.Equal(ts.T(), time.Second*3, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStopStoppedRevertible() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Stop()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Stop()
	assert.EqualError(ts.T(), err, "cannot stop a stopped timer")
	ts.ticker <- t3
	assert.Equal(ts.T(), time.Second*1, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerStopStopped() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Revert()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	err = st.Stop()
	assert.EqualError(ts.T(), err, "cannot stop a stopped timer")
	ts.ticker <- t3
	assert.Equal(ts.T(), time.Second*0, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerReset() {
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)
	t3 := t0.Add(time.Second * 3)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Reset()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t3)
	ts.ticker <- t3
	assert.Equal(ts.T(), time.Second*0, <-c)
}

func (ts *SystemTimerTestSuite) TestSystemTimerClose() {
	logrus.SetLevel(logrus.DebugLevel)
	t0 := time.Now()
	t1 := t0.Add(time.Second * 1)
	t2 := t0.Add(time.Second * 2)

	ts.time.SetNow(t0)
	st := NewSystemTimer(ts.ticker, ts.time)
	c := st.TickChannel()

	ts.time.SetNow(t1)
	err := st.Start()
	assert.NoError(ts.T(), err)

	ts.time.SetNow(t2)
	err = st.Close()
	assert.NoError(ts.T(), err)
	assert.Equal(ts.T(), time.Second*0, <-c)
}
