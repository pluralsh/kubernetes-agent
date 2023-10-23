// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pluralsh/kuberentes-agent/pkg/module/modserver (interfaces: RpcApi,AgentRpcApi)
//
// Generated by this command:
//
//	mockgen -typed -destination rpc_api.go -package mock_modserver github.com/pluralsh/kuberentes-agent/pkg/module/modserver RpcApi,AgentRpcApi
//
// Package mock_modserver is a generated GoMock package.
package mock_modserver

import (
	context "context"
	reflect "reflect"

	api "github.com/pluralsh/kuberentes-agent/pkg/api"
	retry "github.com/pluralsh/kuberentes-agent/pkg/tool/retry"
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
func (mr *MockRpcApiMockRecorder) HandleIoError(arg0, arg1, arg2 any) *RpcApiHandleIoErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIoError", reflect.TypeOf((*MockRpcApi)(nil).HandleIoError), arg0, arg1, arg2)
	return &RpcApiHandleIoErrorCall{Call: call}
}

// RpcApiHandleIoErrorCall wrap *gomock.Call
type RpcApiHandleIoErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RpcApiHandleIoErrorCall) Return(arg0 error) *RpcApiHandleIoErrorCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RpcApiHandleIoErrorCall) Do(f func(*zap.Logger, string, error) error) *RpcApiHandleIoErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RpcApiHandleIoErrorCall) DoAndReturn(f func(*zap.Logger, string, error) error) *RpcApiHandleIoErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HandleProcessingError mocks base method.
func (m *MockRpcApi) HandleProcessingError(arg0 *zap.Logger, arg1 int64, arg2 string, arg3 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleProcessingError", arg0, arg1, arg2, arg3)
}

// HandleProcessingError indicates an expected call of HandleProcessingError.
func (mr *MockRpcApiMockRecorder) HandleProcessingError(arg0, arg1, arg2, arg3 any) *RpcApiHandleProcessingErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProcessingError", reflect.TypeOf((*MockRpcApi)(nil).HandleProcessingError), arg0, arg1, arg2, arg3)
	return &RpcApiHandleProcessingErrorCall{Call: call}
}

// RpcApiHandleProcessingErrorCall wrap *gomock.Call
type RpcApiHandleProcessingErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RpcApiHandleProcessingErrorCall) Return() *RpcApiHandleProcessingErrorCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RpcApiHandleProcessingErrorCall) Do(f func(*zap.Logger, int64, string, error)) *RpcApiHandleProcessingErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RpcApiHandleProcessingErrorCall) DoAndReturn(f func(*zap.Logger, int64, string, error)) *RpcApiHandleProcessingErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Log mocks base method.
func (m *MockRpcApi) Log() *zap.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*zap.Logger)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockRpcApiMockRecorder) Log() *RpcApiLogCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockRpcApi)(nil).Log))
	return &RpcApiLogCall{Call: call}
}

// RpcApiLogCall wrap *gomock.Call
type RpcApiLogCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RpcApiLogCall) Return(arg0 *zap.Logger) *RpcApiLogCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RpcApiLogCall) Do(f func() *zap.Logger) *RpcApiLogCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RpcApiLogCall) DoAndReturn(f func() *zap.Logger) *RpcApiLogCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PollWithBackoff mocks base method.
func (m *MockRpcApi) PollWithBackoff(arg0 retry.PollConfig, arg1 retry.PollWithBackoffFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollWithBackoff", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PollWithBackoff indicates an expected call of PollWithBackoff.
func (mr *MockRpcApiMockRecorder) PollWithBackoff(arg0, arg1 any) *RpcApiPollWithBackoffCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollWithBackoff", reflect.TypeOf((*MockRpcApi)(nil).PollWithBackoff), arg0, arg1)
	return &RpcApiPollWithBackoffCall{Call: call}
}

// RpcApiPollWithBackoffCall wrap *gomock.Call
type RpcApiPollWithBackoffCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RpcApiPollWithBackoffCall) Return(arg0 error) *RpcApiPollWithBackoffCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RpcApiPollWithBackoffCall) Do(f func(retry.PollConfig, retry.PollWithBackoffFunc) error) *RpcApiPollWithBackoffCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RpcApiPollWithBackoffCall) DoAndReturn(f func(retry.PollConfig, retry.PollWithBackoffFunc) error) *RpcApiPollWithBackoffCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockAgentRpcApiMockRecorder) AgentInfo(arg0, arg1 any) *AgentRpcApiAgentInfoCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentInfo", reflect.TypeOf((*MockAgentRpcApi)(nil).AgentInfo), arg0, arg1)
	return &AgentRpcApiAgentInfoCall{Call: call}
}

// AgentRpcApiAgentInfoCall wrap *gomock.Call
type AgentRpcApiAgentInfoCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRpcApiAgentInfoCall) Return(arg0 *api.AgentInfo, arg1 error) *AgentRpcApiAgentInfoCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRpcApiAgentInfoCall) Do(f func(context.Context, *zap.Logger) (*api.AgentInfo, error)) *AgentRpcApiAgentInfoCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRpcApiAgentInfoCall) DoAndReturn(f func(context.Context, *zap.Logger) (*api.AgentInfo, error)) *AgentRpcApiAgentInfoCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// AgentToken mocks base method.
func (m *MockAgentRpcApi) AgentToken() api.AgentToken {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AgentToken")
	ret0, _ := ret[0].(api.AgentToken)
	return ret0
}

// AgentToken indicates an expected call of AgentToken.
func (mr *MockAgentRpcApiMockRecorder) AgentToken() *AgentRpcApiAgentTokenCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AgentToken", reflect.TypeOf((*MockAgentRpcApi)(nil).AgentToken))
	return &AgentRpcApiAgentTokenCall{Call: call}
}

// AgentRpcApiAgentTokenCall wrap *gomock.Call
type AgentRpcApiAgentTokenCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRpcApiAgentTokenCall) Return(arg0 api.AgentToken) *AgentRpcApiAgentTokenCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRpcApiAgentTokenCall) Do(f func() api.AgentToken) *AgentRpcApiAgentTokenCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRpcApiAgentTokenCall) DoAndReturn(f func() api.AgentToken) *AgentRpcApiAgentTokenCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HandleIoError mocks base method.
func (m *MockAgentRpcApi) HandleIoError(arg0 *zap.Logger, arg1 string, arg2 error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleIoError", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleIoError indicates an expected call of HandleIoError.
func (mr *MockAgentRpcApiMockRecorder) HandleIoError(arg0, arg1, arg2 any) *AgentRpcApiHandleIoErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleIoError", reflect.TypeOf((*MockAgentRpcApi)(nil).HandleIoError), arg0, arg1, arg2)
	return &AgentRpcApiHandleIoErrorCall{Call: call}
}

// AgentRpcApiHandleIoErrorCall wrap *gomock.Call
type AgentRpcApiHandleIoErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRpcApiHandleIoErrorCall) Return(arg0 error) *AgentRpcApiHandleIoErrorCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRpcApiHandleIoErrorCall) Do(f func(*zap.Logger, string, error) error) *AgentRpcApiHandleIoErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRpcApiHandleIoErrorCall) DoAndReturn(f func(*zap.Logger, string, error) error) *AgentRpcApiHandleIoErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HandleProcessingError mocks base method.
func (m *MockAgentRpcApi) HandleProcessingError(arg0 *zap.Logger, arg1 int64, arg2 string, arg3 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleProcessingError", arg0, arg1, arg2, arg3)
}

// HandleProcessingError indicates an expected call of HandleProcessingError.
func (mr *MockAgentRpcApiMockRecorder) HandleProcessingError(arg0, arg1, arg2, arg3 any) *AgentRpcApiHandleProcessingErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProcessingError", reflect.TypeOf((*MockAgentRpcApi)(nil).HandleProcessingError), arg0, arg1, arg2, arg3)
	return &AgentRpcApiHandleProcessingErrorCall{Call: call}
}

// AgentRpcApiHandleProcessingErrorCall wrap *gomock.Call
type AgentRpcApiHandleProcessingErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRpcApiHandleProcessingErrorCall) Return() *AgentRpcApiHandleProcessingErrorCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRpcApiHandleProcessingErrorCall) Do(f func(*zap.Logger, int64, string, error)) *AgentRpcApiHandleProcessingErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRpcApiHandleProcessingErrorCall) DoAndReturn(f func(*zap.Logger, int64, string, error)) *AgentRpcApiHandleProcessingErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Log mocks base method.
func (m *MockAgentRpcApi) Log() *zap.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log")
	ret0, _ := ret[0].(*zap.Logger)
	return ret0
}

// Log indicates an expected call of Log.
func (mr *MockAgentRpcApiMockRecorder) Log() *AgentRpcApiLogCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockAgentRpcApi)(nil).Log))
	return &AgentRpcApiLogCall{Call: call}
}

// AgentRpcApiLogCall wrap *gomock.Call
type AgentRpcApiLogCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRpcApiLogCall) Return(arg0 *zap.Logger) *AgentRpcApiLogCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRpcApiLogCall) Do(f func() *zap.Logger) *AgentRpcApiLogCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRpcApiLogCall) DoAndReturn(f func() *zap.Logger) *AgentRpcApiLogCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// PollWithBackoff mocks base method.
func (m *MockAgentRpcApi) PollWithBackoff(arg0 retry.PollConfig, arg1 retry.PollWithBackoffFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PollWithBackoff", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PollWithBackoff indicates an expected call of PollWithBackoff.
func (mr *MockAgentRpcApiMockRecorder) PollWithBackoff(arg0, arg1 any) *AgentRpcApiPollWithBackoffCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollWithBackoff", reflect.TypeOf((*MockAgentRpcApi)(nil).PollWithBackoff), arg0, arg1)
	return &AgentRpcApiPollWithBackoffCall{Call: call}
}

// AgentRpcApiPollWithBackoffCall wrap *gomock.Call
type AgentRpcApiPollWithBackoffCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRpcApiPollWithBackoffCall) Return(arg0 error) *AgentRpcApiPollWithBackoffCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRpcApiPollWithBackoffCall) Do(f func(retry.PollConfig, retry.PollWithBackoffFunc) error) *AgentRpcApiPollWithBackoffCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRpcApiPollWithBackoffCall) DoAndReturn(f func(retry.PollConfig, retry.PollWithBackoffFunc) error) *AgentRpcApiPollWithBackoffCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
