// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel (interfaces: TunnelHandler,TunnelFinder,Tunnel)

// Package mock_reverse_tunnel is a generated GoMock package.
package mock_reverse_tunnel

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	api "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	reverse_tunnel "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel"
	rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/rpc"
	zap "go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

// MockTunnelHandler is a mock of TunnelHandler interface.
type MockTunnelHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTunnelHandlerMockRecorder
}

// MockTunnelHandlerMockRecorder is the mock recorder for MockTunnelHandler.
type MockTunnelHandlerMockRecorder struct {
	mock *MockTunnelHandler
}

// NewMockTunnelHandler creates a new mock instance.
func NewMockTunnelHandler(ctrl *gomock.Controller) *MockTunnelHandler {
	mock := &MockTunnelHandler{ctrl: ctrl}
	mock.recorder = &MockTunnelHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTunnelHandler) EXPECT() *MockTunnelHandlerMockRecorder {
	return m.recorder
}

// HandleTunnel mocks base method.
func (m *MockTunnelHandler) HandleTunnel(arg0 context.Context, arg1 *api.AgentInfo, arg2 rpc.ReverseTunnel_ConnectServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandleTunnel", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HandleTunnel indicates an expected call of HandleTunnel.
func (mr *MockTunnelHandlerMockRecorder) HandleTunnel(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleTunnel", reflect.TypeOf((*MockTunnelHandler)(nil).HandleTunnel), arg0, arg1, arg2)
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
func (m *MockTunnelFinder) FindTunnel(arg0 context.Context, arg1 int64, arg2, arg3 string) (reverse_tunnel.Tunnel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTunnel", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(reverse_tunnel.Tunnel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTunnel indicates an expected call of FindTunnel.
func (mr *MockTunnelFinderMockRecorder) FindTunnel(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTunnel", reflect.TypeOf((*MockTunnelFinder)(nil).FindTunnel), arg0, arg1, arg2, arg3)
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
func (m *MockTunnel) ForwardStream(arg0 *zap.Logger, arg1 reverse_tunnel.RpcApi, arg2 grpc.ServerStream, arg3 reverse_tunnel.TunnelDataCallback) error {
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
