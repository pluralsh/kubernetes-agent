// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker (interfaces: Tracker)
//
// Generated by this command:
//
//	mockgen -typed -destination tracker.go github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker Tracker
//
// Package mock_agent_tracker is a generated GoMock package.
package mock_agent_tracker

import (
	context "context"
	reflect "reflect"

	agent_tracker "github.com/pluralsh/kuberentes-agent/internal/module/agent_tracker"
	gomock "go.uber.org/mock/gomock"
)

// MockTracker is a mock of Tracker interface.
type MockTracker struct {
	ctrl     *gomock.Controller
	recorder *MockTrackerMockRecorder
}

// MockTrackerMockRecorder is the mock recorder for MockTracker.
type MockTrackerMockRecorder struct {
	mock *MockTracker
}

// NewMockTracker creates a new mock instance.
func NewMockTracker(ctrl *gomock.Controller) *MockTracker {
	mock := &MockTracker{ctrl: ctrl}
	mock.recorder = &MockTrackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTracker) EXPECT() *MockTrackerMockRecorder {
	return m.recorder
}

// GetConnectedAgentsCount mocks base method.
func (m *MockTracker) GetConnectedAgentsCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectedAgentsCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConnectedAgentsCount indicates an expected call of GetConnectedAgentsCount.
func (mr *MockTrackerMockRecorder) GetConnectedAgentsCount(arg0 any) *TrackerGetConnectedAgentsCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectedAgentsCount", reflect.TypeOf((*MockTracker)(nil).GetConnectedAgentsCount), arg0)
	return &TrackerGetConnectedAgentsCountCall{Call: call}
}

// TrackerGetConnectedAgentsCountCall wrap *gomock.Call
type TrackerGetConnectedAgentsCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *TrackerGetConnectedAgentsCountCall) Return(arg0 int64, arg1 error) *TrackerGetConnectedAgentsCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *TrackerGetConnectedAgentsCountCall) Do(f func(context.Context) (int64, error)) *TrackerGetConnectedAgentsCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *TrackerGetConnectedAgentsCountCall) DoAndReturn(f func(context.Context) (int64, error)) *TrackerGetConnectedAgentsCountCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetConnectionsByAgentId mocks base method.
func (m *MockTracker) GetConnectionsByAgentId(arg0 context.Context, arg1 int64, arg2 agent_tracker.ConnectedAgentInfoCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionsByAgentId", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetConnectionsByAgentId indicates an expected call of GetConnectionsByAgentId.
func (mr *MockTrackerMockRecorder) GetConnectionsByAgentId(arg0, arg1, arg2 any) *TrackerGetConnectionsByAgentIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionsByAgentId", reflect.TypeOf((*MockTracker)(nil).GetConnectionsByAgentId), arg0, arg1, arg2)
	return &TrackerGetConnectionsByAgentIdCall{Call: call}
}

// TrackerGetConnectionsByAgentIdCall wrap *gomock.Call
type TrackerGetConnectionsByAgentIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *TrackerGetConnectionsByAgentIdCall) Return(arg0 error) *TrackerGetConnectionsByAgentIdCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *TrackerGetConnectionsByAgentIdCall) Do(f func(context.Context, int64, agent_tracker.ConnectedAgentInfoCallback) error) *TrackerGetConnectionsByAgentIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *TrackerGetConnectionsByAgentIdCall) DoAndReturn(f func(context.Context, int64, agent_tracker.ConnectedAgentInfoCallback) error) *TrackerGetConnectionsByAgentIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GetConnectionsByProjectId mocks base method.
func (m *MockTracker) GetConnectionsByProjectId(arg0 context.Context, arg1 int64, arg2 agent_tracker.ConnectedAgentInfoCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnectionsByProjectId", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetConnectionsByProjectId indicates an expected call of GetConnectionsByProjectId.
func (mr *MockTrackerMockRecorder) GetConnectionsByProjectId(arg0, arg1, arg2 any) *TrackerGetConnectionsByProjectIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnectionsByProjectId", reflect.TypeOf((*MockTracker)(nil).GetConnectionsByProjectId), arg0, arg1, arg2)
	return &TrackerGetConnectionsByProjectIdCall{Call: call}
}

// TrackerGetConnectionsByProjectIdCall wrap *gomock.Call
type TrackerGetConnectionsByProjectIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *TrackerGetConnectionsByProjectIdCall) Return(arg0 error) *TrackerGetConnectionsByProjectIdCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *TrackerGetConnectionsByProjectIdCall) Do(f func(context.Context, int64, agent_tracker.ConnectedAgentInfoCallback) error) *TrackerGetConnectionsByProjectIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *TrackerGetConnectionsByProjectIdCall) DoAndReturn(f func(context.Context, int64, agent_tracker.ConnectedAgentInfoCallback) error) *TrackerGetConnectionsByProjectIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RegisterConnection mocks base method.
func (m *MockTracker) RegisterConnection(arg0 context.Context, arg1 *agent_tracker.ConnectedAgentInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterConnection", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterConnection indicates an expected call of RegisterConnection.
func (mr *MockTrackerMockRecorder) RegisterConnection(arg0, arg1 any) *TrackerRegisterConnectionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterConnection", reflect.TypeOf((*MockTracker)(nil).RegisterConnection), arg0, arg1)
	return &TrackerRegisterConnectionCall{Call: call}
}

// TrackerRegisterConnectionCall wrap *gomock.Call
type TrackerRegisterConnectionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *TrackerRegisterConnectionCall) Return(arg0 error) *TrackerRegisterConnectionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *TrackerRegisterConnectionCall) Do(f func(context.Context, *agent_tracker.ConnectedAgentInfo) error) *TrackerRegisterConnectionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *TrackerRegisterConnectionCall) DoAndReturn(f func(context.Context, *agent_tracker.ConnectedAgentInfo) error) *TrackerRegisterConnectionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Run mocks base method.
func (m *MockTracker) Run(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockTrackerMockRecorder) Run(arg0 any) *TrackerRunCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockTracker)(nil).Run), arg0)
	return &TrackerRunCall{Call: call}
}

// TrackerRunCall wrap *gomock.Call
type TrackerRunCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *TrackerRunCall) Return(arg0 error) *TrackerRunCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *TrackerRunCall) Do(f func(context.Context) error) *TrackerRunCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *TrackerRunCall) DoAndReturn(f func(context.Context) error) *TrackerRunCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// UnregisterConnection mocks base method.
func (m *MockTracker) UnregisterConnection(arg0 context.Context, arg1 *agent_tracker.ConnectedAgentInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterConnection", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterConnection indicates an expected call of UnregisterConnection.
func (mr *MockTrackerMockRecorder) UnregisterConnection(arg0, arg1 any) *TrackerUnregisterConnectionCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterConnection", reflect.TypeOf((*MockTracker)(nil).UnregisterConnection), arg0, arg1)
	return &TrackerUnregisterConnectionCall{Call: call}
}

// TrackerUnregisterConnectionCall wrap *gomock.Call
type TrackerUnregisterConnectionCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *TrackerUnregisterConnectionCall) Return(arg0 error) *TrackerUnregisterConnectionCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *TrackerUnregisterConnectionCall) Do(f func(context.Context, *agent_tracker.ConnectedAgentInfo) error) *TrackerUnregisterConnectionCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *TrackerUnregisterConnectionCall) DoAndReturn(f func(context.Context, *agent_tracker.ConnectedAgentInfo) error) *TrackerUnregisterConnectionCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
