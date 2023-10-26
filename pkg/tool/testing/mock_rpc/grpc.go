// Code generated by MockGen. DO NOT EDIT.
// Source: google.golang.org/grpc (interfaces: ServerStream,ClientStream,ClientConnInterface,ServerTransportStream)
//
// Generated by this command:
//
//	mockgen -typed -destination grpc.go -package mock_rpc google.golang.org/grpc ServerStream,ClientStream,ClientConnInterface,ServerTransportStream
//
// Package mock_rpc is a generated GoMock package.
package mock_rpc

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockServerStream is a mock of ServerStream interface.
type MockServerStream struct {
	ctrl     *gomock.Controller
	recorder *MockServerStreamMockRecorder
}

// MockServerStreamMockRecorder is the mock recorder for MockServerStream.
type MockServerStreamMockRecorder struct {
	mock *MockServerStream
}

// NewMockServerStream creates a new mock instance.
func NewMockServerStream(ctrl *gomock.Controller) *MockServerStream {
	mock := &MockServerStream{ctrl: ctrl}
	mock.recorder = &MockServerStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerStream) EXPECT() *MockServerStreamMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockServerStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockServerStreamMockRecorder) Context() *ServerStreamContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockServerStream)(nil).Context))
	return &ServerStreamContextCall{Call: call}
}

// ServerStreamContextCall wrap *gomock.Call
type ServerStreamContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerStreamContextCall) Return(arg0 context.Context) *ServerStreamContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerStreamContextCall) Do(f func() context.Context) *ServerStreamContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerStreamContextCall) DoAndReturn(f func() context.Context) *ServerStreamContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RecvMsg mocks base method.
func (m *MockServerStream) RecvMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockServerStreamMockRecorder) RecvMsg(arg0 any) *ServerStreamRecvMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockServerStream)(nil).RecvMsg), arg0)
	return &ServerStreamRecvMsgCall{Call: call}
}

// ServerStreamRecvMsgCall wrap *gomock.Call
type ServerStreamRecvMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerStreamRecvMsgCall) Return(arg0 error) *ServerStreamRecvMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerStreamRecvMsgCall) Do(f func(any) error) *ServerStreamRecvMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerStreamRecvMsgCall) DoAndReturn(f func(any) error) *ServerStreamRecvMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendHeader mocks base method.
func (m *MockServerStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockServerStreamMockRecorder) SendHeader(arg0 any) *ServerStreamSendHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockServerStream)(nil).SendHeader), arg0)
	return &ServerStreamSendHeaderCall{Call: call}
}

// ServerStreamSendHeaderCall wrap *gomock.Call
type ServerStreamSendHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerStreamSendHeaderCall) Return(arg0 error) *ServerStreamSendHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerStreamSendHeaderCall) Do(f func(metadata.MD) error) *ServerStreamSendHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerStreamSendHeaderCall) DoAndReturn(f func(metadata.MD) error) *ServerStreamSendHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendMsg mocks base method.
func (m *MockServerStream) SendMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockServerStreamMockRecorder) SendMsg(arg0 any) *ServerStreamSendMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockServerStream)(nil).SendMsg), arg0)
	return &ServerStreamSendMsgCall{Call: call}
}

// ServerStreamSendMsgCall wrap *gomock.Call
type ServerStreamSendMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerStreamSendMsgCall) Return(arg0 error) *ServerStreamSendMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerStreamSendMsgCall) Do(f func(any) error) *ServerStreamSendMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerStreamSendMsgCall) DoAndReturn(f func(any) error) *ServerStreamSendMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetHeader mocks base method.
func (m *MockServerStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockServerStreamMockRecorder) SetHeader(arg0 any) *ServerStreamSetHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockServerStream)(nil).SetHeader), arg0)
	return &ServerStreamSetHeaderCall{Call: call}
}

// ServerStreamSetHeaderCall wrap *gomock.Call
type ServerStreamSetHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerStreamSetHeaderCall) Return(arg0 error) *ServerStreamSetHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerStreamSetHeaderCall) Do(f func(metadata.MD) error) *ServerStreamSetHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerStreamSetHeaderCall) DoAndReturn(f func(metadata.MD) error) *ServerStreamSetHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetTrailer mocks base method.
func (m *MockServerStream) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockServerStreamMockRecorder) SetTrailer(arg0 any) *ServerStreamSetTrailerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockServerStream)(nil).SetTrailer), arg0)
	return &ServerStreamSetTrailerCall{Call: call}
}

// ServerStreamSetTrailerCall wrap *gomock.Call
type ServerStreamSetTrailerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerStreamSetTrailerCall) Return() *ServerStreamSetTrailerCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerStreamSetTrailerCall) Do(f func(metadata.MD)) *ServerStreamSetTrailerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerStreamSetTrailerCall) DoAndReturn(f func(metadata.MD)) *ServerStreamSetTrailerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockClientStream is a mock of ClientStream interface.
type MockClientStream struct {
	ctrl     *gomock.Controller
	recorder *MockClientStreamMockRecorder
}

// MockClientStreamMockRecorder is the mock recorder for MockClientStream.
type MockClientStreamMockRecorder struct {
	mock *MockClientStream
}

// NewMockClientStream creates a new mock instance.
func NewMockClientStream(ctrl *gomock.Controller) *MockClientStream {
	mock := &MockClientStream{ctrl: ctrl}
	mock.recorder = &MockClientStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientStream) EXPECT() *MockClientStreamMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockClientStream) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockClientStreamMockRecorder) CloseSend() *ClientStreamCloseSendCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockClientStream)(nil).CloseSend))
	return &ClientStreamCloseSendCall{Call: call}
}

// ClientStreamCloseSendCall wrap *gomock.Call
type ClientStreamCloseSendCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientStreamCloseSendCall) Return(arg0 error) *ClientStreamCloseSendCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientStreamCloseSendCall) Do(f func() error) *ClientStreamCloseSendCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientStreamCloseSendCall) DoAndReturn(f func() error) *ClientStreamCloseSendCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Context mocks base method.
func (m *MockClientStream) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockClientStreamMockRecorder) Context() *ClientStreamContextCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockClientStream)(nil).Context))
	return &ClientStreamContextCall{Call: call}
}

// ClientStreamContextCall wrap *gomock.Call
type ClientStreamContextCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientStreamContextCall) Return(arg0 context.Context) *ClientStreamContextCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientStreamContextCall) Do(f func() context.Context) *ClientStreamContextCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientStreamContextCall) DoAndReturn(f func() context.Context) *ClientStreamContextCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Header mocks base method.
func (m *MockClientStream) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockClientStreamMockRecorder) Header() *ClientStreamHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockClientStream)(nil).Header))
	return &ClientStreamHeaderCall{Call: call}
}

// ClientStreamHeaderCall wrap *gomock.Call
type ClientStreamHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientStreamHeaderCall) Return(arg0 metadata.MD, arg1 error) *ClientStreamHeaderCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientStreamHeaderCall) Do(f func() (metadata.MD, error)) *ClientStreamHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientStreamHeaderCall) DoAndReturn(f func() (metadata.MD, error)) *ClientStreamHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RecvMsg mocks base method.
func (m *MockClientStream) RecvMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockClientStreamMockRecorder) RecvMsg(arg0 any) *ClientStreamRecvMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockClientStream)(nil).RecvMsg), arg0)
	return &ClientStreamRecvMsgCall{Call: call}
}

// ClientStreamRecvMsgCall wrap *gomock.Call
type ClientStreamRecvMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientStreamRecvMsgCall) Return(arg0 error) *ClientStreamRecvMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientStreamRecvMsgCall) Do(f func(any) error) *ClientStreamRecvMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientStreamRecvMsgCall) DoAndReturn(f func(any) error) *ClientStreamRecvMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendMsg mocks base method.
func (m *MockClientStream) SendMsg(arg0 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockClientStreamMockRecorder) SendMsg(arg0 any) *ClientStreamSendMsgCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockClientStream)(nil).SendMsg), arg0)
	return &ClientStreamSendMsgCall{Call: call}
}

// ClientStreamSendMsgCall wrap *gomock.Call
type ClientStreamSendMsgCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientStreamSendMsgCall) Return(arg0 error) *ClientStreamSendMsgCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientStreamSendMsgCall) Do(f func(any) error) *ClientStreamSendMsgCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientStreamSendMsgCall) DoAndReturn(f func(any) error) *ClientStreamSendMsgCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Trailer mocks base method.
func (m *MockClientStream) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockClientStreamMockRecorder) Trailer() *ClientStreamTrailerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockClientStream)(nil).Trailer))
	return &ClientStreamTrailerCall{Call: call}
}

// ClientStreamTrailerCall wrap *gomock.Call
type ClientStreamTrailerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientStreamTrailerCall) Return(arg0 metadata.MD) *ClientStreamTrailerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientStreamTrailerCall) Do(f func() metadata.MD) *ClientStreamTrailerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientStreamTrailerCall) DoAndReturn(f func() metadata.MD) *ClientStreamTrailerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockClientConnInterface is a mock of ClientConnInterface interface.
type MockClientConnInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientConnInterfaceMockRecorder
}

// MockClientConnInterfaceMockRecorder is the mock recorder for MockClientConnInterface.
type MockClientConnInterfaceMockRecorder struct {
	mock *MockClientConnInterface
}

// NewMockClientConnInterface creates a new mock instance.
func NewMockClientConnInterface(ctrl *gomock.Controller) *MockClientConnInterface {
	mock := &MockClientConnInterface{ctrl: ctrl}
	mock.recorder = &MockClientConnInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientConnInterface) EXPECT() *MockClientConnInterfaceMockRecorder {
	return m.recorder
}

// Invoke mocks base method.
func (m *MockClientConnInterface) Invoke(arg0 context.Context, arg1 string, arg2, arg3 any, arg4 ...grpc.CallOption) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Invoke", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Invoke indicates an expected call of Invoke.
func (mr *MockClientConnInterfaceMockRecorder) Invoke(arg0, arg1, arg2, arg3 any, arg4 ...any) *ClientConnInterfaceInvokeCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2, arg3}, arg4...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invoke", reflect.TypeOf((*MockClientConnInterface)(nil).Invoke), varargs...)
	return &ClientConnInterfaceInvokeCall{Call: call}
}

// ClientConnInterfaceInvokeCall wrap *gomock.Call
type ClientConnInterfaceInvokeCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientConnInterfaceInvokeCall) Return(arg0 error) *ClientConnInterfaceInvokeCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientConnInterfaceInvokeCall) Do(f func(context.Context, string, any, any, ...grpc.CallOption) error) *ClientConnInterfaceInvokeCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientConnInterfaceInvokeCall) DoAndReturn(f func(context.Context, string, any, any, ...grpc.CallOption) error) *ClientConnInterfaceInvokeCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// NewStream mocks base method.
func (m *MockClientConnInterface) NewStream(arg0 context.Context, arg1 *grpc.StreamDesc, arg2 string, arg3 ...grpc.CallOption) (grpc.ClientStream, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewStream", varargs...)
	ret0, _ := ret[0].(grpc.ClientStream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream.
func (mr *MockClientConnInterfaceMockRecorder) NewStream(arg0, arg1, arg2 any, arg3 ...any) *ClientConnInterfaceNewStreamCall {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2}, arg3...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*MockClientConnInterface)(nil).NewStream), varargs...)
	return &ClientConnInterfaceNewStreamCall{Call: call}
}

// ClientConnInterfaceNewStreamCall wrap *gomock.Call
type ClientConnInterfaceNewStreamCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ClientConnInterfaceNewStreamCall) Return(arg0 grpc.ClientStream, arg1 error) *ClientConnInterfaceNewStreamCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ClientConnInterfaceNewStreamCall) Do(f func(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error)) *ClientConnInterfaceNewStreamCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ClientConnInterfaceNewStreamCall) DoAndReturn(f func(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error)) *ClientConnInterfaceNewStreamCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// MockServerTransportStream is a mock of ServerTransportStream interface.
type MockServerTransportStream struct {
	ctrl     *gomock.Controller
	recorder *MockServerTransportStreamMockRecorder
}

// MockServerTransportStreamMockRecorder is the mock recorder for MockServerTransportStream.
type MockServerTransportStreamMockRecorder struct {
	mock *MockServerTransportStream
}

// NewMockServerTransportStream creates a new mock instance.
func NewMockServerTransportStream(ctrl *gomock.Controller) *MockServerTransportStream {
	mock := &MockServerTransportStream{ctrl: ctrl}
	mock.recorder = &MockServerTransportStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServerTransportStream) EXPECT() *MockServerTransportStreamMockRecorder {
	return m.recorder
}

// Method mocks base method.
func (m *MockServerTransportStream) Method() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Method")
	ret0, _ := ret[0].(string)
	return ret0
}

// Method indicates an expected call of Method.
func (mr *MockServerTransportStreamMockRecorder) Method() *ServerTransportStreamMethodCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Method", reflect.TypeOf((*MockServerTransportStream)(nil).Method))
	return &ServerTransportStreamMethodCall{Call: call}
}

// ServerTransportStreamMethodCall wrap *gomock.Call
type ServerTransportStreamMethodCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerTransportStreamMethodCall) Return(arg0 string) *ServerTransportStreamMethodCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerTransportStreamMethodCall) Do(f func() string) *ServerTransportStreamMethodCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerTransportStreamMethodCall) DoAndReturn(f func() string) *ServerTransportStreamMethodCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SendHeader mocks base method.
func (m *MockServerTransportStream) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockServerTransportStreamMockRecorder) SendHeader(arg0 any) *ServerTransportStreamSendHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockServerTransportStream)(nil).SendHeader), arg0)
	return &ServerTransportStreamSendHeaderCall{Call: call}
}

// ServerTransportStreamSendHeaderCall wrap *gomock.Call
type ServerTransportStreamSendHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerTransportStreamSendHeaderCall) Return(arg0 error) *ServerTransportStreamSendHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerTransportStreamSendHeaderCall) Do(f func(metadata.MD) error) *ServerTransportStreamSendHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerTransportStreamSendHeaderCall) DoAndReturn(f func(metadata.MD) error) *ServerTransportStreamSendHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetHeader mocks base method.
func (m *MockServerTransportStream) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockServerTransportStreamMockRecorder) SetHeader(arg0 any) *ServerTransportStreamSetHeaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockServerTransportStream)(nil).SetHeader), arg0)
	return &ServerTransportStreamSetHeaderCall{Call: call}
}

// ServerTransportStreamSetHeaderCall wrap *gomock.Call
type ServerTransportStreamSetHeaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerTransportStreamSetHeaderCall) Return(arg0 error) *ServerTransportStreamSetHeaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerTransportStreamSetHeaderCall) Do(f func(metadata.MD) error) *ServerTransportStreamSetHeaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerTransportStreamSetHeaderCall) DoAndReturn(f func(metadata.MD) error) *ServerTransportStreamSetHeaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetTrailer mocks base method.
func (m *MockServerTransportStream) SetTrailer(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTrailer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockServerTransportStreamMockRecorder) SetTrailer(arg0 any) *ServerTransportStreamSetTrailerCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockServerTransportStream)(nil).SetTrailer), arg0)
	return &ServerTransportStreamSetTrailerCall{Call: call}
}

// ServerTransportStreamSetTrailerCall wrap *gomock.Call
type ServerTransportStreamSetTrailerCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ServerTransportStreamSetTrailerCall) Return(arg0 error) *ServerTransportStreamSetTrailerCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ServerTransportStreamSetTrailerCall) Do(f func(metadata.MD) error) *ServerTransportStreamSetTrailerCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ServerTransportStreamSetTrailerCall) DoAndReturn(f func(metadata.MD) error) *ServerTransportStreamSetTrailerCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}