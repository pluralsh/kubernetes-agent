// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc (interfaces: AgentRegistrarClient)
//
// Generated by this command:
//
//	mockgen -typed -destination agent_registrar.go -package mock_agent_registrar gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc AgentRegistrarClient
//
// Package mock_agent_registrar is a generated GoMock package.
package mock_agent_registrar

import (
	context "context"
	reflect "reflect"

	rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/agent_registrar/rpc"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAgentRegistrarClient is a mock of AgentRegistrarClient interface.
type MockAgentRegistrarClient struct {
	ctrl     *gomock.Controller
	recorder *MockAgentRegistrarClientMockRecorder
}

// MockAgentRegistrarClientMockRecorder is the mock recorder for MockAgentRegistrarClient.
type MockAgentRegistrarClientMockRecorder struct {
	mock *MockAgentRegistrarClient
}

// NewMockAgentRegistrarClient creates a new mock instance.
func NewMockAgentRegistrarClient(ctrl *gomock.Controller) *MockAgentRegistrarClient {
	mock := &MockAgentRegistrarClient{ctrl: ctrl}
	mock.recorder = &MockAgentRegistrarClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAgentRegistrarClient) EXPECT() *MockAgentRegistrarClientMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockAgentRegistrarClient) Register(arg0 context.Context, arg1 *rpc.RegisterRequest, arg2 ...grpc.CallOption) (*rpc.RegisterResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Register", varargs...)
	ret0, _ := ret[0].(*rpc.RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAgentRegistrarClientMockRecorder) Register(arg0, arg1 any, arg2 ...any) *AgentRegistrarClientRegisterCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAgentRegistrarClient)(nil).Register), varargs...)
	return &AgentRegistrarClientRegisterCall{Call: call}
}

// AgentRegistrarClientRegisterCall wrap *gomock.Call
type AgentRegistrarClientRegisterCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *AgentRegistrarClientRegisterCall) Return(arg0 *rpc.RegisterResponse, arg1 error) *AgentRegistrarClientRegisterCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *AgentRegistrarClientRegisterCall) Do(f func(context.Context, *rpc.RegisterRequest, ...grpc.CallOption) (*rpc.RegisterResponse, error)) *AgentRegistrarClientRegisterCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *AgentRegistrarClientRegisterCall) DoAndReturn(f func(context.Context, *rpc.RegisterRequest, ...grpc.CallOption) (*rpc.RegisterResponse, error)) *AgentRegistrarClientRegisterCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
