// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/reverse_tunnel (interfaces: TunnelHandler)

// Package mock_reverse_tunnel is a generated GoMock package.
package mock_reverse_tunnel

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	api "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/api"
	rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/reverse_tunnel/rpc"
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
