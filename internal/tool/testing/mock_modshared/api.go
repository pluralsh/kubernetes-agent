// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared (interfaces: RpcApi)

// Package mock_modshared is a generated GoMock package.
package mock_modshared

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	retry "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/retry"
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

// HandleSendError mocks base method.
func (m *MockRpcApi) HandleSendError(arg0 *zap.Logger, arg1 string, arg2 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleSendError", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleSendError indicates an expected call of HandleSendError.
func (mr *MockRpcApiMockRecorder) HandleSendError(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleSendError", reflect.TypeOf((*MockRpcApi)(nil).HandleSendError), arg0, arg1, arg2)
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