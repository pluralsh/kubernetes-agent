// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/kubernetes_api/rpc (interfaces: KubernetesApiClient,KubernetesApi_MakeRequestClient)
//
// Generated by this command:
//
//	mockgen -typed -destination rpc.go -package mock_kubernetes_api gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/kubernetes_api/rpc KubernetesApiClient,KubernetesApi_MakeRequestClient
//
// Package mock_kubernetes_api is a generated GoMock package.
package mock_kubernetes_api

import (
	context "context"
	reflect "reflect"

	rpc "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/kubernetes_api/rpc"
	grpctool "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/grpctool"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockKubernetesApiClient is a mock of KubernetesApiClient interface.
type MockKubernetesApiClient struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesApiClientMockRecorder
}

// MockKubernetesApiClientMockRecorder is the mock recorder for MockKubernetesApiClient.
type MockKubernetesApiClientMockRecorder struct {
	mock *MockKubernetesApiClient
}

// NewMockKubernetesApiClient creates a new mock instance.
func NewMockKubernetesApiClient(ctrl *gomock.Controller) *MockKubernetesApiClient {
	mock := &MockKubernetesApiClient{ctrl: ctrl}
	mock.recorder = &MockKubernetesApiClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetesApiClient) EXPECT() *MockKubernetesApiClientMockRecorder {
	return m.recorder
}

// MakeRequest mocks base method.
func (m *MockKubernetesApiClient) MakeRequest(arg0 context.Context, arg1 ...grpc.CallOption) (rpc.KubernetesApi_MakeRequestClient, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MakeRequest", varargs...)
	ret0, _ := ret[0].(rpc.KubernetesApi_MakeRequestClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeRequest indicates an expected call of MakeRequest.
func (mr *MockKubernetesApiClientMockRecorder) MakeRequest(arg0 any, arg1 ...any) *KubernetesApiClientMakeRequestCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeRequest", reflect.TypeOf((*MockKubernetesApiClient)(nil).MakeRequest), varargs...)
	return &KubernetesApiClientMakeRequestCall{Call: call}
}

// KubernetesApiClientMakeRequestCall wrap *gomock.Call
type KubernetesApiClientMakeRequestCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApiClientMakeRequestCall) Return(arg0 rpc.KubernetesApi_MakeRequestClient, arg1 error) *KubernetesApiClientMakeRequestCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApiClientMakeRequestCall) Do(f func(context.Context, ...grpc.CallOption) (rpc.KubernetesApi_MakeRequestClient, error)) *KubernetesApiClientMakeRequestCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApiClientMakeRequestCall) DoAndReturn(f func(context.Context, ...grpc.CallOption) (rpc.KubernetesApi_MakeRequestClient, error)) *KubernetesApiClientMakeRequestCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockKubernetesApi_MakeRequestClient is a mock of KubernetesApi_MakeRequestClient interface.
type MockKubernetesApi_MakeRequestClient struct {
	ctrl     *gomock.Controller
	recorder *MockKubernetesApi_MakeRequestClientMockRecorder
}

// MockKubernetesApi_MakeRequestClientMockRecorder is the mock recorder for MockKubernetesApi_MakeRequestClient.
type MockKubernetesApi_MakeRequestClientMockRecorder struct {
	mock *MockKubernetesApi_MakeRequestClient
}

// NewMockKubernetesApi_MakeRequestClient creates a new mock instance.
func NewMockKubernetesApi_MakeRequestClient(ctrl *gomock.Controller) *MockKubernetesApi_MakeRequestClient {
	mock := &MockKubernetesApi_MakeRequestClient{ctrl: ctrl}
	mock.recorder = &MockKubernetesApi_MakeRequestClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKubernetesApi_MakeRequestClient) EXPECT() *MockKubernetesApi_MakeRequestClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) CloseSend() *KubernetesApi_MakeRequestClientCloseSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).CloseSend))
	return &KubernetesApi_MakeRequestClientCloseSendCall{Call: call}
}

// KubernetesApi_MakeRequestClientCloseSendCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientCloseSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientCloseSendCall) Return(arg0 error) *KubernetesApi_MakeRequestClientCloseSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientCloseSendCall) Do(f func() error) *KubernetesApi_MakeRequestClientCloseSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientCloseSendCall) DoAndReturn(f func() error) *KubernetesApi_MakeRequestClientCloseSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Context mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) Context() *KubernetesApi_MakeRequestClientContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).Context))
	return &KubernetesApi_MakeRequestClientContextCall{Call: call}
}

// KubernetesApi_MakeRequestClientContextCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientContextCall) Return(arg0 context.Context) *KubernetesApi_MakeRequestClientContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientContextCall) Do(f func() context.Context) *KubernetesApi_MakeRequestClientContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientContextCall) DoAndReturn(f func() context.Context) *KubernetesApi_MakeRequestClientContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Header mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) Header() *KubernetesApi_MakeRequestClientHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).Header))
	return &KubernetesApi_MakeRequestClientHeaderCall{Call: call}
}

// KubernetesApi_MakeRequestClientHeaderCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientHeaderCall) Return(arg0 metadata.MD, arg1 error) *KubernetesApi_MakeRequestClientHeaderCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientHeaderCall) Do(f func() (metadata.MD, error)) *KubernetesApi_MakeRequestClientHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientHeaderCall) DoAndReturn(f func() (metadata.MD, error)) *KubernetesApi_MakeRequestClientHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Recv mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) Recv() (*grpctool.HttpResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*grpctool.HttpResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) Recv() *KubernetesApi_MakeRequestClientRecvCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).Recv))
	return &KubernetesApi_MakeRequestClientRecvCall{Call: call}
}

// KubernetesApi_MakeRequestClientRecvCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientRecvCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientRecvCall) Return(arg0 *grpctool.HttpResponse, arg1 error) *KubernetesApi_MakeRequestClientRecvCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientRecvCall) Do(f func() (*grpctool.HttpResponse, error)) *KubernetesApi_MakeRequestClientRecvCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientRecvCall) DoAndReturn(f func() (*grpctool.HttpResponse, error)) *KubernetesApi_MakeRequestClientRecvCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RecvMsg mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) RecvMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) RecvMsg(arg0 any) *KubernetesApi_MakeRequestClientRecvMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).RecvMsg), arg0)
	return &KubernetesApi_MakeRequestClientRecvMsgCall{Call: call}
}

// KubernetesApi_MakeRequestClientRecvMsgCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientRecvMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientRecvMsgCall) Return(arg0 error) *KubernetesApi_MakeRequestClientRecvMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientRecvMsgCall) Do(f func(any) error) *KubernetesApi_MakeRequestClientRecvMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientRecvMsgCall) DoAndReturn(f func(any) error) *KubernetesApi_MakeRequestClientRecvMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Send mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) Send(arg0 *grpctool.HttpRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) Send(arg0 any) *KubernetesApi_MakeRequestClientSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).Send), arg0)
	return &KubernetesApi_MakeRequestClientSendCall{Call: call}
}

// KubernetesApi_MakeRequestClientSendCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientSendCall) Return(arg0 error) *KubernetesApi_MakeRequestClientSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientSendCall) Do(f func(*grpctool.HttpRequest) error) *KubernetesApi_MakeRequestClientSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientSendCall) DoAndReturn(f func(*grpctool.HttpRequest) error) *KubernetesApi_MakeRequestClientSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendMsg mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) SendMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) SendMsg(arg0 any) *KubernetesApi_MakeRequestClientSendMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).SendMsg), arg0)
	return &KubernetesApi_MakeRequestClientSendMsgCall{Call: call}
}

// KubernetesApi_MakeRequestClientSendMsgCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientSendMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientSendMsgCall) Return(arg0 error) *KubernetesApi_MakeRequestClientSendMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientSendMsgCall) Do(f func(any) error) *KubernetesApi_MakeRequestClientSendMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientSendMsgCall) DoAndReturn(f func(any) error) *KubernetesApi_MakeRequestClientSendMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Trailer mocks base method.
func (m *MockKubernetesApi_MakeRequestClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockKubernetesApi_MakeRequestClientMockRecorder) Trailer() *KubernetesApi_MakeRequestClientTrailerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockKubernetesApi_MakeRequestClient)(nil).Trailer))
	return &KubernetesApi_MakeRequestClientTrailerCall{Call: call}
}

// KubernetesApi_MakeRequestClientTrailerCall wrap *gomock.Call
type KubernetesApi_MakeRequestClientTrailerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *KubernetesApi_MakeRequestClientTrailerCall) Return(arg0 metadata.MD) *KubernetesApi_MakeRequestClientTrailerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *KubernetesApi_MakeRequestClientTrailerCall) Do(f func() metadata.MD) *KubernetesApi_MakeRequestClientTrailerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *KubernetesApi_MakeRequestClientTrailerCall) DoAndReturn(f func() metadata.MD) *KubernetesApi_MakeRequestClientTrailerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
