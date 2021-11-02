// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool (interfaces: InboundGrpcToOutboundHttpStream,PoolConn)

// Package mock_rpc is a generated GoMock package.
package mock_rpc

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpctool "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/grpctool"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockInboundGrpcToOutboundHttpStream is a mock of InboundGrpcToOutboundHttpStream interface.
type MockInboundGrpcToOutboundHttpStream struct {
	ctrl     *gomock.Controller
	recorder *MockInboundGrpcToOutboundHttpStreamMockRecorder
}

// MockInboundGrpcToOutboundHttpStreamMockRecorder is the mock recorder for MockInboundGrpcToOutboundHttpStream.
type MockInboundGrpcToOutboundHttpStreamMockRecorder struct {
	mock *MockInboundGrpcToOutboundHttpStream
}

// NewMockInboundGrpcToOutboundHttpStream creates a new mock instance.
func NewMockInboundGrpcToOutboundHttpStream(ctrl *gomock.Controller) *MockInboundGrpcToOutboundHttpStream {
	mock := &MockInboundGrpcToOutboundHttpStream{ctrl: ctrl}
	mock.recorder = &MockInboundGrpcToOutboundHttpStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInboundGrpcToOutboundHttpStream) EXPECT() *MockInboundGrpcToOutboundHttpStreamMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) Send(arg0 *grpctool.HttpResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockInboundGrpcToOutboundHttpStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockInboundGrpcToOutboundHttpStreamMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockInboundGrpcToOutboundHttpStream)(nil).SetTrailer), arg0)
}

// MockPoolConn is a mock of PoolConn interface.
type MockPoolConn struct {
	ctrl     *gomock.Controller
	recorder *MockPoolConnMockRecorder
}

// MockPoolConnMockRecorder is the mock recorder for MockPoolConn.
type MockPoolConnMockRecorder struct {
	mock *MockPoolConn
}

// NewMockPoolConn creates a new mock instance.
func NewMockPoolConn(ctrl *gomock.Controller) *MockPoolConn {
	mock := &MockPoolConn{ctrl: ctrl}
	mock.recorder = &MockPoolConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPoolConn) EXPECT() *MockPoolConnMockRecorder {
	return m.recorder
}

// Done mocks base method.
func (m *MockPoolConn) Done() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Done")
}

// Done indicates an expected call of Done.
func (mr *MockPoolConnMockRecorder) Done() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockPoolConn)(nil).Done))
}

// Invoke mocks base method.
func (m *MockPoolConn) Invoke(arg0 context.Context, arg1 string, arg2, arg3 interface{}, arg4 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Invoke", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Invoke indicates an expected call of Invoke.
func (mr *MockPoolConnMockRecorder) Invoke(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockPoolConn)(nil).Invoke), varargs...)
}

// NewStream mocks base method.
func (m *MockPoolConn) NewStream(arg0 context.Context, arg1 *grpc.StreamDesc, arg2 string, arg3 ...grpc.CallOption) (grpc.ClientStream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewStream", varargs...)
	ret0, _ := ret[0].(grpc.ClientStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream.
func (mr *MockPoolConnMockRecorder) NewStream(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*MockPoolConn)(nil).NewStream), varargs...)
}
