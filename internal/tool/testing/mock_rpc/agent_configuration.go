// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration/rpc (interfaces: AgentConfigurationClient,AgentConfiguration_GetConfigurationClient,AgentConfiguration_GetConfigurationServer,ConfigurationWatcherInterface)

// Package mock_rpc is a generated GoMock package.
package mock_rpc

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration/rpc"
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
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfiguration", varargs...)
	ret0, _ := ret[0].(rpc.AgentConfiguration_GetConfigurationClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfiguration indicates an expected call of GetConfiguration.
func (mr *MockAgentConfigurationClientMockRecorder) GetConfiguration(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfiguration", reflect.TypeOf((*MockAgentConfigurationClient)(nil).GetConfiguration), varargs...)
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
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Context))
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
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Header))
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
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockAgentConfiguration_GetConfigurationClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockAgentConfiguration_GetConfigurationClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationClient)(nil).Trailer))
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
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) Send(arg0 *rpc.ConfigurationResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockAgentConfiguration_GetConfigurationServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockAgentConfiguration_GetConfigurationServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockAgentConfiguration_GetConfigurationServer)(nil).SetTrailer), arg0)
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
func (mr *MockConfigurationWatcherInterfaceMockRecorder) Watch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockConfigurationWatcherInterface)(nil).Watch), arg0, arg1)
}
