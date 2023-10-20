// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pluralsh/kuberentes-agent/cmd/agentk/agentkapp (interfaces: Runner,LeaderElector)
//
// Generated by this command:
//
//	mockgen -typed -destination mock_for_test.go -package agentkapp github.com/pluralsh/kuberentes-agent/cmd/agentk/agentkapp Runner,LeaderElector
//
// Package agentkapp is a generated GoMock package.
package agentkapp

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRunner is a mock of Runner interface.
type MockRunner struct {
	ctrl     *gomock.Controller
	recorder *MockRunnerMockRecorder
}

// MockRunnerMockRecorder is the mock recorder for MockRunner.
type MockRunnerMockRecorder struct {
	mock *MockRunner
}

// NewMockRunner creates a new mock instance.
func NewMockRunner(ctrl *gomock.Controller) *MockRunner {
	mock := &MockRunner{ctrl: ctrl}
	mock.recorder = &MockRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRunner) EXPECT() *MockRunnerMockRecorder {
	return m.recorder
}

// RunWhenLeader mocks base method.
func (m *MockRunner) RunWhenLeader(arg0 context.Context, arg1 ModuleStartFunc, arg2 ModuleStopFunc) (CancelRunWhenLeaderFunc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunWhenLeader", arg0, arg1, arg2)
	ret0, _ := ret[0].(CancelRunWhenLeaderFunc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunWhenLeader indicates an expected call of RunWhenLeader.
func (mr *MockRunnerMockRecorder) RunWhenLeader(arg0, arg1, arg2 any) *RunnerRunWhenLeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunWhenLeader", reflect.TypeOf((*MockRunner)(nil).RunWhenLeader), arg0, arg1, arg2)
	return &RunnerRunWhenLeaderCall{Call: call}
}

// RunnerRunWhenLeaderCall wrap *gomock.Call
type RunnerRunWhenLeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RunnerRunWhenLeaderCall) Return(arg0 CancelRunWhenLeaderFunc, arg1 error) *RunnerRunWhenLeaderCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RunnerRunWhenLeaderCall) Do(f func(context.Context, ModuleStartFunc, ModuleStopFunc) (CancelRunWhenLeaderFunc, error)) *RunnerRunWhenLeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RunnerRunWhenLeaderCall) DoAndReturn(f func(context.Context, ModuleStartFunc, ModuleStopFunc) (CancelRunWhenLeaderFunc, error)) *RunnerRunWhenLeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockLeaderElector is a mock of LeaderElector interface.
type MockLeaderElector struct {
	ctrl     *gomock.Controller
	recorder *MockLeaderElectorMockRecorder
}

// MockLeaderElectorMockRecorder is the mock recorder for MockLeaderElector.
type MockLeaderElectorMockRecorder struct {
	mock *MockLeaderElector
}

// NewMockLeaderElector creates a new mock instance.
func NewMockLeaderElector(ctrl *gomock.Controller) *MockLeaderElector {
	mock := &MockLeaderElector{ctrl: ctrl}
	mock.recorder = &MockLeaderElectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLeaderElector) EXPECT() *MockLeaderElectorMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockLeaderElector) Run(arg0 context.Context, arg1, arg2 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0, arg1, arg2)
}

// Run indicates an expected call of Run.
func (mr *MockLeaderElectorMockRecorder) Run(arg0, arg1, arg2 any) *LeaderElectorRunCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockLeaderElector)(nil).Run), arg0, arg1, arg2)
	return &LeaderElectorRunCall{Call: call}
}

// LeaderElectorRunCall wrap *gomock.Call
type LeaderElectorRunCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *LeaderElectorRunCall) Return() *LeaderElectorRunCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *LeaderElectorRunCall) Do(f func(context.Context, func(), func())) *LeaderElectorRunCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *LeaderElectorRunCall) DoAndReturn(f func(context.Context, func(), func())) *LeaderElectorRunCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
