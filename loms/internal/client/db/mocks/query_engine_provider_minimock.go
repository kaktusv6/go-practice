package mocks

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i route256/libs/db.QueryEngineProvider -o ./mocks\query_engine_provider_minimock.go -n QueryEngineProviderMock

import (
	"context"
	mm_db "route256/libs/db"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// QueryEngineProviderMock implements db.QueryEngineProvider
type QueryEngineProviderMock struct {
	t minimock.Tester

	funcGetQueryEngine          func(ctx context.Context) (q1 mm_db.QueryEngine)
	inspectFuncGetQueryEngine   func(ctx context.Context)
	afterGetQueryEngineCounter  uint64
	beforeGetQueryEngineCounter uint64
	GetQueryEngineMock          mQueryEngineProviderMockGetQueryEngine
}

// NewQueryEngineProviderMock returns a mock for db.QueryEngineProvider
func NewQueryEngineProviderMock(t minimock.Tester) *QueryEngineProviderMock {
	m := &QueryEngineProviderMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetQueryEngineMock = mQueryEngineProviderMockGetQueryEngine{mock: m}
	m.GetQueryEngineMock.callArgs = []*QueryEngineProviderMockGetQueryEngineParams{}

	return m
}

type mQueryEngineProviderMockGetQueryEngine struct {
	mock               *QueryEngineProviderMock
	defaultExpectation *QueryEngineProviderMockGetQueryEngineExpectation
	expectations       []*QueryEngineProviderMockGetQueryEngineExpectation

	callArgs []*QueryEngineProviderMockGetQueryEngineParams
	mutex    sync.RWMutex
}

// QueryEngineProviderMockGetQueryEngineExpectation specifies expectation struct of the QueryEngineProvider.GetQueryEngine
type QueryEngineProviderMockGetQueryEngineExpectation struct {
	mock    *QueryEngineProviderMock
	params  *QueryEngineProviderMockGetQueryEngineParams
	results *QueryEngineProviderMockGetQueryEngineResults
	Counter uint64
}

// QueryEngineProviderMockGetQueryEngineParams contains parameters of the QueryEngineProvider.GetQueryEngine
type QueryEngineProviderMockGetQueryEngineParams struct {
	ctx context.Context
}

// QueryEngineProviderMockGetQueryEngineResults contains results of the QueryEngineProvider.GetQueryEngine
type QueryEngineProviderMockGetQueryEngineResults struct {
	q1 mm_db.QueryEngine
}

// Expect sets up expected params for QueryEngineProvider.GetQueryEngine
func (mmGetQueryEngine *mQueryEngineProviderMockGetQueryEngine) Expect(ctx context.Context) *mQueryEngineProviderMockGetQueryEngine {
	if mmGetQueryEngine.mock.funcGetQueryEngine != nil {
		mmGetQueryEngine.mock.t.Fatalf("QueryEngineProviderMock.GetQueryEngine mock is already set by Set")
	}

	if mmGetQueryEngine.defaultExpectation == nil {
		mmGetQueryEngine.defaultExpectation = &QueryEngineProviderMockGetQueryEngineExpectation{}
	}

	mmGetQueryEngine.defaultExpectation.params = &QueryEngineProviderMockGetQueryEngineParams{ctx}
	for _, e := range mmGetQueryEngine.expectations {
		if minimock.Equal(e.params, mmGetQueryEngine.defaultExpectation.params) {
			mmGetQueryEngine.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetQueryEngine.defaultExpectation.params)
		}
	}

	return mmGetQueryEngine
}

// Inspect accepts an inspector function that has same arguments as the QueryEngineProvider.GetQueryEngine
func (mmGetQueryEngine *mQueryEngineProviderMockGetQueryEngine) Inspect(f func(ctx context.Context)) *mQueryEngineProviderMockGetQueryEngine {
	if mmGetQueryEngine.mock.inspectFuncGetQueryEngine != nil {
		mmGetQueryEngine.mock.t.Fatalf("Inspect function is already set for QueryEngineProviderMock.GetQueryEngine")
	}

	mmGetQueryEngine.mock.inspectFuncGetQueryEngine = f

	return mmGetQueryEngine
}

// Return sets up results that will be returned by QueryEngineProvider.GetQueryEngine
func (mmGetQueryEngine *mQueryEngineProviderMockGetQueryEngine) Return(q1 mm_db.QueryEngine) *QueryEngineProviderMock {
	if mmGetQueryEngine.mock.funcGetQueryEngine != nil {
		mmGetQueryEngine.mock.t.Fatalf("QueryEngineProviderMock.GetQueryEngine mock is already set by Set")
	}

	if mmGetQueryEngine.defaultExpectation == nil {
		mmGetQueryEngine.defaultExpectation = &QueryEngineProviderMockGetQueryEngineExpectation{mock: mmGetQueryEngine.mock}
	}
	mmGetQueryEngine.defaultExpectation.results = &QueryEngineProviderMockGetQueryEngineResults{q1}
	return mmGetQueryEngine.mock
}

// Set uses given function f to mock the QueryEngineProvider.GetQueryEngine method
func (mmGetQueryEngine *mQueryEngineProviderMockGetQueryEngine) Set(f func(ctx context.Context) (q1 mm_db.QueryEngine)) *QueryEngineProviderMock {
	if mmGetQueryEngine.defaultExpectation != nil {
		mmGetQueryEngine.mock.t.Fatalf("Default expectation is already set for the QueryEngineProvider.GetQueryEngine method")
	}

	if len(mmGetQueryEngine.expectations) > 0 {
		mmGetQueryEngine.mock.t.Fatalf("Some expectations are already set for the QueryEngineProvider.GetQueryEngine method")
	}

	mmGetQueryEngine.mock.funcGetQueryEngine = f
	return mmGetQueryEngine.mock
}

// When sets expectation for the QueryEngineProvider.GetQueryEngine which will trigger the result defined by the following
// Then helper
func (mmGetQueryEngine *mQueryEngineProviderMockGetQueryEngine) When(ctx context.Context) *QueryEngineProviderMockGetQueryEngineExpectation {
	if mmGetQueryEngine.mock.funcGetQueryEngine != nil {
		mmGetQueryEngine.mock.t.Fatalf("QueryEngineProviderMock.GetQueryEngine mock is already set by Set")
	}

	expectation := &QueryEngineProviderMockGetQueryEngineExpectation{
		mock:   mmGetQueryEngine.mock,
		params: &QueryEngineProviderMockGetQueryEngineParams{ctx},
	}
	mmGetQueryEngine.expectations = append(mmGetQueryEngine.expectations, expectation)
	return expectation
}

// Then sets up QueryEngineProvider.GetQueryEngine return parameters for the expectation previously defined by the When method
func (e *QueryEngineProviderMockGetQueryEngineExpectation) Then(q1 mm_db.QueryEngine) *QueryEngineProviderMock {
	e.results = &QueryEngineProviderMockGetQueryEngineResults{q1}
	return e.mock
}

// GetQueryEngine implements db.QueryEngineProvider
func (mmGetQueryEngine *QueryEngineProviderMock) GetQueryEngine(ctx context.Context) (q1 mm_db.QueryEngine) {
	mm_atomic.AddUint64(&mmGetQueryEngine.beforeGetQueryEngineCounter, 1)
	defer mm_atomic.AddUint64(&mmGetQueryEngine.afterGetQueryEngineCounter, 1)

	if mmGetQueryEngine.inspectFuncGetQueryEngine != nil {
		mmGetQueryEngine.inspectFuncGetQueryEngine(ctx)
	}

	mm_params := &QueryEngineProviderMockGetQueryEngineParams{ctx}

	// Record call args
	mmGetQueryEngine.GetQueryEngineMock.mutex.Lock()
	mmGetQueryEngine.GetQueryEngineMock.callArgs = append(mmGetQueryEngine.GetQueryEngineMock.callArgs, mm_params)
	mmGetQueryEngine.GetQueryEngineMock.mutex.Unlock()

	for _, e := range mmGetQueryEngine.GetQueryEngineMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.q1
		}
	}

	if mmGetQueryEngine.GetQueryEngineMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetQueryEngine.GetQueryEngineMock.defaultExpectation.Counter, 1)
		mm_want := mmGetQueryEngine.GetQueryEngineMock.defaultExpectation.params
		mm_got := QueryEngineProviderMockGetQueryEngineParams{ctx}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetQueryEngine.t.Errorf("QueryEngineProviderMock.GetQueryEngine got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetQueryEngine.GetQueryEngineMock.defaultExpectation.results
		if mm_results == nil {
			mmGetQueryEngine.t.Fatal("No results are set for the QueryEngineProviderMock.GetQueryEngine")
		}
		return (*mm_results).q1
	}
	if mmGetQueryEngine.funcGetQueryEngine != nil {
		return mmGetQueryEngine.funcGetQueryEngine(ctx)
	}
	mmGetQueryEngine.t.Fatalf("Unexpected call to QueryEngineProviderMock.GetQueryEngine. %v", ctx)
	return
}

// GetQueryEngineAfterCounter returns a count of finished QueryEngineProviderMock.GetQueryEngine invocations
func (mmGetQueryEngine *QueryEngineProviderMock) GetQueryEngineAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetQueryEngine.afterGetQueryEngineCounter)
}

// GetQueryEngineBeforeCounter returns a count of QueryEngineProviderMock.GetQueryEngine invocations
func (mmGetQueryEngine *QueryEngineProviderMock) GetQueryEngineBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetQueryEngine.beforeGetQueryEngineCounter)
}

// Calls returns a list of arguments used in each call to QueryEngineProviderMock.GetQueryEngine.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetQueryEngine *mQueryEngineProviderMockGetQueryEngine) Calls() []*QueryEngineProviderMockGetQueryEngineParams {
	mmGetQueryEngine.mutex.RLock()

	argCopy := make([]*QueryEngineProviderMockGetQueryEngineParams, len(mmGetQueryEngine.callArgs))
	copy(argCopy, mmGetQueryEngine.callArgs)

	mmGetQueryEngine.mutex.RUnlock()

	return argCopy
}

// MinimockGetQueryEngineDone returns true if the count of the GetQueryEngine invocations corresponds
// the number of defined expectations
func (m *QueryEngineProviderMock) MinimockGetQueryEngineDone() bool {
	for _, e := range m.GetQueryEngineMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetQueryEngineMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetQueryEngineCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetQueryEngine != nil && mm_atomic.LoadUint64(&m.afterGetQueryEngineCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetQueryEngineInspect logs each unmet expectation
func (m *QueryEngineProviderMock) MinimockGetQueryEngineInspect() {
	for _, e := range m.GetQueryEngineMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to QueryEngineProviderMock.GetQueryEngine with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetQueryEngineMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetQueryEngineCounter) < 1 {
		if m.GetQueryEngineMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to QueryEngineProviderMock.GetQueryEngine")
		} else {
			m.t.Errorf("Expected call to QueryEngineProviderMock.GetQueryEngine with params: %#v", *m.GetQueryEngineMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetQueryEngine != nil && mm_atomic.LoadUint64(&m.afterGetQueryEngineCounter) < 1 {
		m.t.Error("Expected call to QueryEngineProviderMock.GetQueryEngine")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *QueryEngineProviderMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetQueryEngineInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *QueryEngineProviderMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *QueryEngineProviderMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetQueryEngineDone()
}
