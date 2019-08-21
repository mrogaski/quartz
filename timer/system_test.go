package timer

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
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

func (ts *SystemTimerTestSuite) TestSystemTimerRun() {
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

func (ts *SystemTimerTestSuite) TestSystemTimerRunRevert() {
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
