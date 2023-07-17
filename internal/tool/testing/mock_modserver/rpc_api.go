// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modserver (interfaces: RpcApi,AgentRpcApi)

// Package mock_modserver is a generated GoMock package.
package mock_modserver

import (
	context "context"
	reflect "reflect"

	api "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	retry "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/retry"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockRpcApi is a mock of RpcApi interface.
type MockRpcApi struct {
	ctrl     *gomock.Controller
	recorder *MockRpcApiMockRecorder
}

// MockRpcApiMockRecorder is the mock recorder for MockRpcApi.
type MockRpcApiMockRecorder struct {
	mock *MockRpcApi
}

// NewMockRpcApi creates a new mock instance.
func NewMockRpcApi(ctrl *gomock.Controller) *MockRpcApi {
	mock := &MockRpcApi{ctrl: ctrl}
	mock.recorder = &MockRpcApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRpcApi) EXPECT() *MockRpcApiMockRecorder {
	return m.recorder
}

// HandleIoError mocks base method.
func (m *MockRpcApi) HandleIoError(arg0 *zap.Logger, arg1 string, arg2 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleIoError", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleIoError indicates an expected call of HandleIoError.
func (mr *MockRpcApiMockRecorder) HandleIoError(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIoError", reflect.TypeOf((*MockRpcApi)(nil).HandleIoError), arg0, arg1, arg2)
}

// HandleProcessingError mocks base method.
func (m *MockRpcApi) HandleProcessingError(arg0 *zap.Logger, arg1 int64, arg2 string, arg3 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleProcessingError", arg0, arg1, arg2, arg3)
}

// HandleProcessingError indicates an expected call of HandleProcessingError.
func (mr *MockRpcApiMockRecorder) HandleProcessingError(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProcessingError", reflect.TypeOf((*MockRpcApi)(nil).HandleProcessingError), arg0, arg1, arg2, arg3)
}

// Log mocks base method.
func (m *MockRpcApi) Log() *zap.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*zap.Logger)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockRpcApiMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockRpcApi)(nil).Log))
}

// PollWithBackoff mocks base method.
func (m *MockRpcApi) PollWithBackoff(arg0 retry.PollConfig, arg1 retry.PollWithBackoffFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollWithBackoff", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PollWithBackoff indicates an expected call of PollWithBackoff.
func (mr *MockRpcApiMockRecorder) PollWithBackoff(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollWithBackoff", reflect.TypeOf((*MockRpcApi)(nil).PollWithBackoff), arg0, arg1)
}

// MockAgentRpcApi is a mock of AgentRpcApi interface.
type MockAgentRpcApi struct {
	ctrl     *gomock.Controller
	recorder *MockAgentRpcApiMockRecorder
}

// MockAgentRpcApiMockRecorder is the mock recorder for MockAgentRpcApi.
type MockAgentRpcApiMockRecorder struct {
	mock *MockAgentRpcApi
}

// NewMockAgentRpcApi creates a new mock instance.
func NewMockAgentRpcApi(ctrl *gomock.Controller) *MockAgentRpcApi {
	mock := &MockAgentRpcApi{ctrl: ctrl}
	mock.recorder = &MockAgentRpcApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentRpcApi) EXPECT() *MockAgentRpcApiMockRecorder {
	return m.recorder
}

// AgentInfo mocks base method.
func (m *MockAgentRpcApi) AgentInfo(arg0 context.Context, arg1 *zap.Logger) (*api.AgentInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentInfo", arg0, arg1)
	ret0, _ := ret[0].(*api.AgentInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AgentInfo indicates an expected call of AgentInfo.
func (mr *MockAgentRpcApiMockRecorder) AgentInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentInfo", reflect.TypeOf((*MockAgentRpcApi)(nil).AgentInfo), arg0, arg1)
}

// AgentToken mocks base method.
func (m *MockAgentRpcApi) AgentToken() api.AgentToken {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentToken")
	ret0, _ := ret[0].(api.AgentToken)
	return ret0
}

// AgentToken indicates an expected call of AgentToken.
func (mr *MockAgentRpcApiMockRecorder) AgentToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentToken", reflect.TypeOf((*MockAgentRpcApi)(nil).AgentToken))
}

// HandleIoError mocks base method.
func (m *MockAgentRpcApi) HandleIoError(arg0 *zap.Logger, arg1 string, arg2 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleIoError", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleIoError indicates an expected call of HandleIoError.
func (mr *MockAgentRpcApiMockRecorder) HandleIoError(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIoError", reflect.TypeOf((*MockAgentRpcApi)(nil).HandleIoError), arg0, arg1, arg2)
}

// HandleProcessingError mocks base method.
func (m *MockAgentRpcApi) HandleProcessingError(arg0 *zap.Logger, arg1 int64, arg2 string, arg3 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleProcessingError", arg0, arg1, arg2, arg3)
}

// HandleProcessingError indicates an expected call of HandleProcessingError.
func (mr *MockAgentRpcApiMockRecorder) HandleProcessingError(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProcessingError", reflect.TypeOf((*MockAgentRpcApi)(nil).HandleProcessingError), arg0, arg1, arg2, arg3)
}

// Log mocks base method.
func (m *MockAgentRpcApi) Log() *zap.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*zap.Logger)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockAgentRpcApiMockRecorder) Log() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockAgentRpcApi)(nil).Log))
}

// PollWithBackoff mocks base method.
func (m *MockAgentRpcApi) PollWithBackoff(arg0 retry.PollConfig, arg1 retry.PollWithBackoffFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollWithBackoff", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PollWithBackoff indicates an expected call of PollWithBackoff.
func (mr *MockAgentRpcApiMockRecorder) PollWithBackoff(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollWithBackoff", reflect.TypeOf((*MockAgentRpcApi)(nil).PollWithBackoff), arg0, arg1)
}
