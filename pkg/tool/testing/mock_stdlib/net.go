// Code generated by MockGen. DO NOT EDIT.
// Source: net (interfaces: Conn)
//
// Generated by this command:
//
//	mockgen -typed -destination net.go -package mock_stdlib net Conn
//
// Package mock_stdlib is a generated GoMock package.
package mock_stdlib

import (
	net "net"
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockConn is a mock of Conn interface.
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn.
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance.
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnMockRecorder) Close() *ConnCloseCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
	return &ConnCloseCall{Call: call}
}

// ConnCloseCall wrap *gomock.Call
type ConnCloseCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnCloseCall) Return(arg0 error) *ConnCloseCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnCloseCall) Do(f func() error) *ConnCloseCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnCloseCall) DoAndReturn(f func() error) *ConnCloseCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// LocalAddr mocks base method.
func (m *MockConn) LocalAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocalAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// LocalAddr indicates an expected call of LocalAddr.
func (mr *MockConnMockRecorder) LocalAddr() *ConnLocalAddrCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocalAddr", reflect.TypeOf((*MockConn)(nil).LocalAddr))
	return &ConnLocalAddrCall{Call: call}
}

// ConnLocalAddrCall wrap *gomock.Call
type ConnLocalAddrCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnLocalAddrCall) Return(arg0 net.Addr) *ConnLocalAddrCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnLocalAddrCall) Do(f func() net.Addr) *ConnLocalAddrCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnLocalAddrCall) DoAndReturn(f func() net.Addr) *ConnLocalAddrCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Read mocks base method.
func (m *MockConn) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockConnMockRecorder) Read(arg0 any) *ConnReadCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockConn)(nil).Read), arg0)
	return &ConnReadCall{Call: call}
}

// ConnReadCall wrap *gomock.Call
type ConnReadCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnReadCall) Return(arg0 int, arg1 error) *ConnReadCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnReadCall) Do(f func([]byte) (int, error)) *ConnReadCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnReadCall) DoAndReturn(f func([]byte) (int, error)) *ConnReadCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// RemoteAddr mocks base method.
func (m *MockConn) RemoteAddr() net.Addr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoteAddr")
	ret0, _ := ret[0].(net.Addr)
	return ret0
}

// RemoteAddr indicates an expected call of RemoteAddr.
func (mr *MockConnMockRecorder) RemoteAddr() *ConnRemoteAddrCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoteAddr", reflect.TypeOf((*MockConn)(nil).RemoteAddr))
	return &ConnRemoteAddrCall{Call: call}
}

// ConnRemoteAddrCall wrap *gomock.Call
type ConnRemoteAddrCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnRemoteAddrCall) Return(arg0 net.Addr) *ConnRemoteAddrCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnRemoteAddrCall) Do(f func() net.Addr) *ConnRemoteAddrCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnRemoteAddrCall) DoAndReturn(f func() net.Addr) *ConnRemoteAddrCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetDeadline mocks base method.
func (m *MockConn) SetDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDeadline indicates an expected call of SetDeadline.
func (mr *MockConnMockRecorder) SetDeadline(arg0 any) *ConnSetDeadlineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDeadline", reflect.TypeOf((*MockConn)(nil).SetDeadline), arg0)
	return &ConnSetDeadlineCall{Call: call}
}

// ConnSetDeadlineCall wrap *gomock.Call
type ConnSetDeadlineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnSetDeadlineCall) Return(arg0 error) *ConnSetDeadlineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnSetDeadlineCall) Do(f func(time.Time) error) *ConnSetDeadlineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnSetDeadlineCall) DoAndReturn(f func(time.Time) error) *ConnSetDeadlineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetReadDeadline mocks base method.
func (m *MockConn) SetReadDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetReadDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetReadDeadline indicates an expected call of SetReadDeadline.
func (mr *MockConnMockRecorder) SetReadDeadline(arg0 any) *ConnSetReadDeadlineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetReadDeadline", reflect.TypeOf((*MockConn)(nil).SetReadDeadline), arg0)
	return &ConnSetReadDeadlineCall{Call: call}
}

// ConnSetReadDeadlineCall wrap *gomock.Call
type ConnSetReadDeadlineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnSetReadDeadlineCall) Return(arg0 error) *ConnSetReadDeadlineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnSetReadDeadlineCall) Do(f func(time.Time) error) *ConnSetReadDeadlineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnSetReadDeadlineCall) DoAndReturn(f func(time.Time) error) *ConnSetReadDeadlineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// SetWriteDeadline mocks base method.
func (m *MockConn) SetWriteDeadline(arg0 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWriteDeadline", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetWriteDeadline indicates an expected call of SetWriteDeadline.
func (mr *MockConnMockRecorder) SetWriteDeadline(arg0 any) *ConnSetWriteDeadlineCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWriteDeadline", reflect.TypeOf((*MockConn)(nil).SetWriteDeadline), arg0)
	return &ConnSetWriteDeadlineCall{Call: call}
}

// ConnSetWriteDeadlineCall wrap *gomock.Call
type ConnSetWriteDeadlineCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnSetWriteDeadlineCall) Return(arg0 error) *ConnSetWriteDeadlineCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnSetWriteDeadlineCall) Do(f func(time.Time) error) *ConnSetWriteDeadlineCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnSetWriteDeadlineCall) DoAndReturn(f func(time.Time) error) *ConnSetWriteDeadlineCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Write mocks base method.
func (m *MockConn) Write(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockConnMockRecorder) Write(arg0 any) *ConnWriteCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockConn)(nil).Write), arg0)
	return &ConnWriteCall{Call: call}
}

// ConnWriteCall wrap *gomock.Call
type ConnWriteCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ConnWriteCall) Return(arg0 int, arg1 error) *ConnWriteCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ConnWriteCall) Do(f func([]byte) (int, error)) *ConnWriteCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ConnWriteCall) DoAndReturn(f func([]byte) (int, error)) *ConnWriteCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}