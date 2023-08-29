// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel (interfaces: Registerer,Handler,FindHandle,Tunnel,PollingQuerier,TunnelFinder)

// Package mock_reverse_tunnel_tunnel is a generated GoMock package.
package mock_reverse_tunnel_tunnel

import (
	context "context"
	reflect "reflect"

	api "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/api"
	rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/rpc"
	tunnel "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

// MockRegisterer is a mock of Registerer interface.
type MockRegisterer struct {
	ctrl     *gomock.Controller
	recorder *MockRegistererMockRecorder
}

// MockRegistererMockRecorder is the mock recorder for MockRegisterer.
type MockRegistererMockRecorder struct {
	mock *MockRegisterer
}

// NewMockRegisterer creates a new mock instance.
func NewMockRegisterer(ctrl *gomock.Controller) *MockRegisterer {
	mock := &MockRegisterer{ctrl: ctrl}
	mock.recorder = &MockRegistererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRegisterer) EXPECT() *MockRegistererMockRecorder {
	return m.recorder
}

// RegisterTunnel mocks base method.
func (m *MockRegisterer) RegisterTunnel(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTunnel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterTunnel indicates an expected call of RegisterTunnel.
func (mr *MockRegistererMockRecorder) RegisterTunnel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTunnel", reflect.TypeOf((*MockRegisterer)(nil).RegisterTunnel), arg0, arg1)
}

// UnregisterTunnel mocks base method.
func (m *MockRegisterer) UnregisterTunnel(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterTunnel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterTunnel indicates an expected call of UnregisterTunnel.
func (mr *MockRegistererMockRecorder) UnregisterTunnel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterTunnel", reflect.TypeOf((*MockRegisterer)(nil).UnregisterTunnel), arg0, arg1)
}

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// HandleTunnel mocks base method.
func (m *MockHandler) HandleTunnel(arg0 context.Context, arg1 *api.AgentInfo, arg2 rpc.ReverseTunnel_ConnectServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleTunnel", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleTunnel indicates an expected call of HandleTunnel.
func (mr *MockHandlerMockRecorder) HandleTunnel(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleTunnel", reflect.TypeOf((*MockHandler)(nil).HandleTunnel), arg0, arg1, arg2)
}

// MockFindHandle is a mock of FindHandle interface.
type MockFindHandle struct {
	ctrl     *gomock.Controller
	recorder *MockFindHandleMockRecorder
}

// MockFindHandleMockRecorder is the mock recorder for MockFindHandle.
type MockFindHandleMockRecorder struct {
	mock *MockFindHandle
}

// NewMockFindHandle creates a new mock instance.
func NewMockFindHandle(ctrl *gomock.Controller) *MockFindHandle {
	mock := &MockFindHandle{ctrl: ctrl}
	mock.recorder = &MockFindHandleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFindHandle) EXPECT() *MockFindHandleMockRecorder {
	return m.recorder
}

// Done mocks base method.
func (m *MockFindHandle) Done() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Done")
}

// Done indicates an expected call of Done.
func (mr *MockFindHandleMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockFindHandle)(nil).Done))
}

// Get mocks base method.
func (m *MockFindHandle) Get(arg0 context.Context) (tunnel.Tunnel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(tunnel.Tunnel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockFindHandleMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockFindHandle)(nil).Get), arg0)
}

// MockTunnel is a mock of Tunnel interface.
type MockTunnel struct {
	ctrl     *gomock.Controller
	recorder *MockTunnelMockRecorder
}

// MockTunnelMockRecorder is the mock recorder for MockTunnel.
type MockTunnelMockRecorder struct {
	mock *MockTunnel
}

// NewMockTunnel creates a new mock instance.
func NewMockTunnel(ctrl *gomock.Controller) *MockTunnel {
	mock := &MockTunnel{ctrl: ctrl}
	mock.recorder = &MockTunnelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTunnel) EXPECT() *MockTunnelMockRecorder {
	return m.recorder
}

// Done mocks base method.
func (m *MockTunnel) Done() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Done")
}

// Done indicates an expected call of Done.
func (mr *MockTunnelMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockTunnel)(nil).Done))
}

// ForwardStream mocks base method.
func (m *MockTunnel) ForwardStream(arg0 *zap.Logger, arg1 tunnel.RpcApi, arg2 grpc.ServerStream, arg3 tunnel.DataCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForwardStream", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForwardStream indicates an expected call of ForwardStream.
func (mr *MockTunnelMockRecorder) ForwardStream(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForwardStream", reflect.TypeOf((*MockTunnel)(nil).ForwardStream), arg0, arg1, arg2, arg3)
}

// MockPollingQuerier is a mock of PollingQuerier interface.
type MockPollingQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockPollingQuerierMockRecorder
}

// MockPollingQuerierMockRecorder is the mock recorder for MockPollingQuerier.
type MockPollingQuerierMockRecorder struct {
	mock *MockPollingQuerier
}

// NewMockPollingQuerier creates a new mock instance.
func NewMockPollingQuerier(ctrl *gomock.Controller) *MockPollingQuerier {
	mock := &MockPollingQuerier{ctrl: ctrl}
	mock.recorder = &MockPollingQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPollingQuerier) EXPECT() *MockPollingQuerierMockRecorder {
	return m.recorder
}

// CachedKasUrlsByAgentId mocks base method.
func (m *MockPollingQuerier) CachedKasUrlsByAgentId(arg0 int64) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CachedKasUrlsByAgentId", arg0)
	ret0, _ := ret[0].([]string)
	return ret0
}

// CachedKasUrlsByAgentId indicates an expected call of CachedKasUrlsByAgentId.
func (mr *MockPollingQuerierMockRecorder) CachedKasUrlsByAgentId(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CachedKasUrlsByAgentId", reflect.TypeOf((*MockPollingQuerier)(nil).CachedKasUrlsByAgentId), arg0)
}

// PollKasUrlsByAgentId mocks base method.
func (m *MockPollingQuerier) PollKasUrlsByAgentId(arg0 context.Context, arg1 int64, arg2 tunnel.PollKasUrlsByAgentIdCallback) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PollKasUrlsByAgentId", arg0, arg1, arg2)
}

// PollKasUrlsByAgentId indicates an expected call of PollKasUrlsByAgentId.
func (mr *MockPollingQuerierMockRecorder) PollKasUrlsByAgentId(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PollKasUrlsByAgentId", reflect.TypeOf((*MockPollingQuerier)(nil).PollKasUrlsByAgentId), arg0, arg1, arg2)
}

// MockTunnelFinder is a mock of TunnelFinder interface.
type MockTunnelFinder struct {
	ctrl     *gomock.Controller
	recorder *MockTunnelFinderMockRecorder
}

// MockTunnelFinderMockRecorder is the mock recorder for MockTunnelFinder.
type MockTunnelFinderMockRecorder struct {
	mock *MockTunnelFinder
}

// NewMockTunnelFinder creates a new mock instance.
func NewMockTunnelFinder(ctrl *gomock.Controller) *MockTunnelFinder {
	mock := &MockTunnelFinder{ctrl: ctrl}
	mock.recorder = &MockTunnelFinderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTunnelFinder) EXPECT() *MockTunnelFinderMockRecorder {
	return m.recorder
}

// FindTunnel mocks base method.
func (m *MockTunnelFinder) FindTunnel(arg0 int64, arg1, arg2 string) (bool, tunnel.FindHandle) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTunnel", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(tunnel.FindHandle)
	return ret0, ret1
}

// FindTunnel indicates an expected call of FindTunnel.
func (mr *MockTunnelFinderMockRecorder) FindTunnel(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTunnel", reflect.TypeOf((*MockTunnelFinder)(nil).FindTunnel), arg0, arg1, arg2)
}
