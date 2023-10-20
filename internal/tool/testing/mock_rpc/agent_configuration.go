// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/rpc (interfaces: AgentConfigurationClient,AgentConfiguration_GetConfigurationClient,AgentConfiguration_GetConfigurationServer,ConfigurationWatcherInterface)
//
// Generated by this command:
//
//	mockgen -typed -destination agent_configuration.go -package mock_rpc gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_configuration/rpc AgentConfigurationClient,AgentConfiguration_GetConfigurationClient,AgentConfiguration_GetConfigurationServer,ConfigurationWatcherInterface
//
// Package mock_rpc is a generated GoMock package.
package mock_rpc

import (
	context "context"
	reflect "reflect"

	rpc "github.com/pluralsh/kuberentes-agent/internal/module/agent_configuration/rpc"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockAgentConfigurationClient is a mock of AgentConfigurationClient interface.
type MockAgentConfigurationClient struct {
	ctrl     *gomock.Controller
	recorder *MockAgentConfigurationClientMockRecorder
}

// MockAgentConfigurationClientMockRecorder is the mock recorder for MockAgentConfigurationClient.
type MockAgentConfigurationClientMockRecorder struct {
	mock *MockAgentConfigurationClient
}

// NewMockAgentConfigurationClient creates a new mock instance.
func NewMockAgentConfigurationClient(ctrl *gomock.Controller) *MockAgentConfigurationClient {
	mock := &MockAgentConfigurationClient{ctrl: ctrl}
	mock.recorder = &MockAgentConfigurationClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentConfigurationClient) EXPECT() *MockAgentConfigurationClientMockRecorder {
	return m.recorder
}

// GetConfiguration mocks base method.
func (m *MockAgentConfigurationClient) GetConfiguration(arg0 context.Context, arg1 *rpc.ConfigurationRequest, arg2 ...grpc.CallOption) (rpc.AgentConfiguration_GetConfigurationClient, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfiguration", varargs...)
	ret0, _ := ret[0].(rpc.AgentConfiguration_GetConfigurationClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfiguration indicates an expected call of GetConfiguration.
func (mr *MockAgentConfigurationClientMockRecorder) GetConfiguration(arg0, arg1 any, arg2 ...any) *AgentConfigurationClientGetConfigurationCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfiguration", reflect.TypeOf((*MockAgentConfigurationClient)(nil).GetConfiguration), varargs...)
	return &AgentConfigurationClientGetConfigurationCall{Call: call}
}

// AgentConfigurationClientGetConfigurationCall wrap *gomock.Call
type AgentConfigurationClientGetConfigurationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfigurationClientGetConfigurationCall) Return(arg0 rpc.AgentConfiguration_GetConfigurationClient, arg1 error) *AgentConfigurationClientGetConfigurationCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfigurationClientGetConfigurationCall) Do(f func(context.Context, *rpc.ConfigurationRequest, ...grpc.CallOption) (rpc.AgentConfiguration_GetConfigurationClient, error)) *AgentConfigurationClientGetConfigurationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfigurationClientGetConfigurationCall) DoAndReturn(f func(context.Context, *rpc.ConfigurationRequest, ...grpc.CallOption) (rpc.AgentConfiguration_GetConfigurationClient, error)) *AgentConfigurationClientGetConfigurationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockAgentConfiguration_GetConfigurationClient is a mock of AgentConfiguration_GetConfigurationClient interface.
type MockAgentConfiguration_GetConfigurationClient struct {
	ctrl     *gomock.Controller
	recorder *MockAgentConfiguration_GetConfigurationClientMockRecorder
}

// MockAgentConfiguration_GetConfigurationClientMockRecorder is the mock recorder for MockAgentConfiguration_GetConfigurationClient.
type MockAgentConfiguration_GetConfigurationClientMockRecorder struct {
	mock *MockAgentConfiguration_GetConfigurationClient
}

// NewMockAgentConfiguration_GetConfigurationClient creates a new mock instance.
func NewMockAgentConfiguration_GetConfigurationClient(ctrl *gomock.Controller) *MockAgentConfiguration_GetConfigurationClient {
	mock := &MockAgentConfiguration_GetConfigurationClient{ctrl: ctrl}
	mock.recorder = &MockAgentConfiguration_GetConfigurationClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentConfiguration_GetConfigurationClient) EXPECT() *MockAgentConfiguration_GetConfigurationClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) CloseSend() *AgentConfiguration_GetConfigurationClientCloseSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).CloseSend))
	return &AgentConfiguration_GetConfigurationClientCloseSendCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientCloseSendCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientCloseSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientCloseSendCall) Return(arg0 error) *AgentConfiguration_GetConfigurationClientCloseSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientCloseSendCall) Do(f func() error) *AgentConfiguration_GetConfigurationClientCloseSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientCloseSendCall) DoAndReturn(f func() error) *AgentConfiguration_GetConfigurationClientCloseSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Context mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Context() *AgentConfiguration_GetConfigurationClientContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Context))
	return &AgentConfiguration_GetConfigurationClientContextCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientContextCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientContextCall) Return(arg0 context.Context) *AgentConfiguration_GetConfigurationClientContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientContextCall) Do(f func() context.Context) *AgentConfiguration_GetConfigurationClientContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientContextCall) DoAndReturn(f func() context.Context) *AgentConfiguration_GetConfigurationClientContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Header mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Header() *AgentConfiguration_GetConfigurationClientHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Header))
	return &AgentConfiguration_GetConfigurationClientHeaderCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientHeaderCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientHeaderCall) Return(arg0 metadata.MD, arg1 error) *AgentConfiguration_GetConfigurationClientHeaderCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientHeaderCall) Do(f func() (metadata.MD, error)) *AgentConfiguration_GetConfigurationClientHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientHeaderCall) DoAndReturn(f func() (metadata.MD, error)) *AgentConfiguration_GetConfigurationClientHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Recv mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) Recv() (*rpc.ConfigurationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*rpc.ConfigurationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Recv() *AgentConfiguration_GetConfigurationClientRecvCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Recv))
	return &AgentConfiguration_GetConfigurationClientRecvCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientRecvCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientRecvCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientRecvCall) Return(arg0 *rpc.ConfigurationResponse, arg1 error) *AgentConfiguration_GetConfigurationClientRecvCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientRecvCall) Do(f func() (*rpc.ConfigurationResponse, error)) *AgentConfiguration_GetConfigurationClientRecvCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientRecvCall) DoAndReturn(f func() (*rpc.ConfigurationResponse, error)) *AgentConfiguration_GetConfigurationClientRecvCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RecvMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) RecvMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) RecvMsg(arg0 any) *AgentConfiguration_GetConfigurationClientRecvMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).RecvMsg), arg0)
	return &AgentConfiguration_GetConfigurationClientRecvMsgCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientRecvMsgCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientRecvMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientRecvMsgCall) Return(arg0 error) *AgentConfiguration_GetConfigurationClientRecvMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientRecvMsgCall) Do(f func(any) error) *AgentConfiguration_GetConfigurationClientRecvMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientRecvMsgCall) DoAndReturn(f func(any) error) *AgentConfiguration_GetConfigurationClientRecvMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) SendMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) SendMsg(arg0 any) *AgentConfiguration_GetConfigurationClientSendMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).SendMsg), arg0)
	return &AgentConfiguration_GetConfigurationClientSendMsgCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientSendMsgCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientSendMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientSendMsgCall) Return(arg0 error) *AgentConfiguration_GetConfigurationClientSendMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientSendMsgCall) Do(f func(any) error) *AgentConfiguration_GetConfigurationClientSendMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientSendMsgCall) DoAndReturn(f func(any) error) *AgentConfiguration_GetConfigurationClientSendMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Trailer mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Trailer() *AgentConfiguration_GetConfigurationClientTrailerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Trailer))
	return &AgentConfiguration_GetConfigurationClientTrailerCall{Call: call}
}

// AgentConfiguration_GetConfigurationClientTrailerCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationClientTrailerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationClientTrailerCall) Return(arg0 metadata.MD) *AgentConfiguration_GetConfigurationClientTrailerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationClientTrailerCall) Do(f func() metadata.MD) *AgentConfiguration_GetConfigurationClientTrailerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationClientTrailerCall) DoAndReturn(f func() metadata.MD) *AgentConfiguration_GetConfigurationClientTrailerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockAgentConfiguration_GetConfigurationServer is a mock of AgentConfiguration_GetConfigurationServer interface.
type MockAgentConfiguration_GetConfigurationServer struct {
	ctrl     *gomock.Controller
	recorder *MockAgentConfiguration_GetConfigurationServerMockRecorder
}

// MockAgentConfiguration_GetConfigurationServerMockRecorder is the mock recorder for MockAgentConfiguration_GetConfigurationServer.
type MockAgentConfiguration_GetConfigurationServerMockRecorder struct {
	mock *MockAgentConfiguration_GetConfigurationServer
}

// NewMockAgentConfiguration_GetConfigurationServer creates a new mock instance.
func NewMockAgentConfiguration_GetConfigurationServer(ctrl *gomock.Controller) *MockAgentConfiguration_GetConfigurationServer {
	mock := &MockAgentConfiguration_GetConfigurationServer{ctrl: ctrl}
	mock.recorder = &MockAgentConfiguration_GetConfigurationServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentConfiguration_GetConfigurationServer) EXPECT() *MockAgentConfiguration_GetConfigurationServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) Context() *AgentConfiguration_GetConfigurationServerContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).Context))
	return &AgentConfiguration_GetConfigurationServerContextCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerContextCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerContextCall) Return(arg0 context.Context) *AgentConfiguration_GetConfigurationServerContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerContextCall) Do(f func() context.Context) *AgentConfiguration_GetConfigurationServerContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerContextCall) DoAndReturn(f func() context.Context) *AgentConfiguration_GetConfigurationServerContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RecvMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) RecvMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) RecvMsg(arg0 any) *AgentConfiguration_GetConfigurationServerRecvMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).RecvMsg), arg0)
	return &AgentConfiguration_GetConfigurationServerRecvMsgCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerRecvMsgCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerRecvMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerRecvMsgCall) Return(arg0 error) *AgentConfiguration_GetConfigurationServerRecvMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerRecvMsgCall) Do(f func(any) error) *AgentConfiguration_GetConfigurationServerRecvMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerRecvMsgCall) DoAndReturn(f func(any) error) *AgentConfiguration_GetConfigurationServerRecvMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Send mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) Send(arg0 *rpc.ConfigurationResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) Send(arg0 any) *AgentConfiguration_GetConfigurationServerSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).Send), arg0)
	return &AgentConfiguration_GetConfigurationServerSendCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerSendCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerSendCall) Return(arg0 error) *AgentConfiguration_GetConfigurationServerSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerSendCall) Do(f func(*rpc.ConfigurationResponse) error) *AgentConfiguration_GetConfigurationServerSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerSendCall) DoAndReturn(f func(*rpc.ConfigurationResponse) error) *AgentConfiguration_GetConfigurationServerSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendHeader mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SendHeader(arg0 any) *AgentConfiguration_GetConfigurationServerSendHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SendHeader), arg0)
	return &AgentConfiguration_GetConfigurationServerSendHeaderCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerSendHeaderCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerSendHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerSendHeaderCall) Return(arg0 error) *AgentConfiguration_GetConfigurationServerSendHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerSendHeaderCall) Do(f func(metadata.MD) error) *AgentConfiguration_GetConfigurationServerSendHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerSendHeaderCall) DoAndReturn(f func(metadata.MD) error) *AgentConfiguration_GetConfigurationServerSendHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SendMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SendMsg(arg0 any) *AgentConfiguration_GetConfigurationServerSendMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SendMsg), arg0)
	return &AgentConfiguration_GetConfigurationServerSendMsgCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerSendMsgCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerSendMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerSendMsgCall) Return(arg0 error) *AgentConfiguration_GetConfigurationServerSendMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerSendMsgCall) Do(f func(any) error) *AgentConfiguration_GetConfigurationServerSendMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerSendMsgCall) DoAndReturn(f func(any) error) *AgentConfiguration_GetConfigurationServerSendMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetHeader mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SetHeader(arg0 any) *AgentConfiguration_GetConfigurationServerSetHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SetHeader), arg0)
	return &AgentConfiguration_GetConfigurationServerSetHeaderCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerSetHeaderCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerSetHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerSetHeaderCall) Return(arg0 error) *AgentConfiguration_GetConfigurationServerSetHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerSetHeaderCall) Do(f func(metadata.MD) error) *AgentConfiguration_GetConfigurationServerSetHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerSetHeaderCall) DoAndReturn(f func(metadata.MD) error) *AgentConfiguration_GetConfigurationServerSetHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetTrailer mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SetTrailer(arg0 any) *AgentConfiguration_GetConfigurationServerSetTrailerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SetTrailer), arg0)
	return &AgentConfiguration_GetConfigurationServerSetTrailerCall{Call: call}
}

// AgentConfiguration_GetConfigurationServerSetTrailerCall wrap *gomock.Call
type AgentConfiguration_GetConfigurationServerSetTrailerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentConfiguration_GetConfigurationServerSetTrailerCall) Return() *AgentConfiguration_GetConfigurationServerSetTrailerCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentConfiguration_GetConfigurationServerSetTrailerCall) Do(f func(metadata.MD)) *AgentConfiguration_GetConfigurationServerSetTrailerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentConfiguration_GetConfigurationServerSetTrailerCall) DoAndReturn(f func(metadata.MD)) *AgentConfiguration_GetConfigurationServerSetTrailerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockConfigurationWatcherInterface is a mock of ConfigurationWatcherInterface interface.
type MockConfigurationWatcherInterface struct {
	ctrl     *gomock.Controller
	recorder *MockConfigurationWatcherInterfaceMockRecorder
}

// MockConfigurationWatcherInterfaceMockRecorder is the mock recorder for MockConfigurationWatcherInterface.
type MockConfigurationWatcherInterfaceMockRecorder struct {
	mock *MockConfigurationWatcherInterface
}

// NewMockConfigurationWatcherInterface creates a new mock instance.
func NewMockConfigurationWatcherInterface(ctrl *gomock.Controller) *MockConfigurationWatcherInterface {
	mock := &MockConfigurationWatcherInterface{ctrl: ctrl}
	mock.recorder = &MockConfigurationWatcherInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigurationWatcherInterface) EXPECT() *MockConfigurationWatcherInterfaceMockRecorder {
	return m.recorder
}

// Watch mocks base method.
func (m *MockConfigurationWatcherInterface) Watch(arg0 context.Context, arg1 rpc.ConfigurationCallback) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Watch", arg0, arg1)
}

// Watch indicates an expected call of Watch.
func (mr *MockConfigurationWatcherInterfaceMockRecorder) Watch(arg0, arg1 any) *ConfigurationWatcherInterfaceWatchCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockConfigurationWatcherInterface)(nil).Watch), arg0, arg1)
	return &ConfigurationWatcherInterfaceWatchCall{Call: call}
}

// ConfigurationWatcherInterfaceWatchCall wrap *gomock.Call
type ConfigurationWatcherInterfaceWatchCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConfigurationWatcherInterfaceWatchCall) Return() *ConfigurationWatcherInterfaceWatchCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConfigurationWatcherInterfaceWatchCall) Do(f func(context.Context, rpc.ConfigurationCallback)) *ConfigurationWatcherInterfaceWatchCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConfigurationWatcherInterfaceWatchCall) DoAndReturn(f func(context.Context, rpc.ConfigurationCallback)) *ConfigurationWatcherInterfaceWatchCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
