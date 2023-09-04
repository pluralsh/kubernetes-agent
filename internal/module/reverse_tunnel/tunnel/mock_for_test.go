// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/reverse_tunnel/tunnel (interfaces: DataCallback,Querier,Tracker)

// Package tunnel is a generated GoMock package.
package tunnel

import (
	context "context"
	reflect "reflect"
	time "time"

	prototool "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	gomock "go.uber.org/mock/gomock"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

// MockDataCallback is a mock of DataCallback interface.
type MockDataCallback struct {
	ctrl     *gomock.Controller
	recorder *MockDataCallbackMockRecorder
}

// MockDataCallbackMockRecorder is the mock recorder for MockDataCallback.
type MockDataCallbackMockRecorder struct {
	mock *MockDataCallback
}

// NewMockDataCallback creates a new mock instance.
func NewMockDataCallback(ctrl *gomock.Controller) *MockDataCallback {
	mock := &MockDataCallback{ctrl: ctrl}
	mock.recorder = &MockDataCallbackMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataCallback) EXPECT() *MockDataCallbackMockRecorder {
	return m.recorder
}

// Error mocks base method.
func (m *MockDataCallback) Error(arg0 *status.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockDataCallbackMockRecorder) Error(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockDataCallback)(nil).Error), arg0)
}

// Header mocks base method.
func (m *MockDataCallback) Header(arg0 map[string]*prototool.Values) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Header indicates an expected call of Header.
func (mr *MockDataCallbackMockRecorder) Header(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDataCallback)(nil).Header), arg0)
}

// Message mocks base method.
func (m *MockDataCallback) Message(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Message", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Message indicates an expected call of Message.
func (mr *MockDataCallbackMockRecorder) Message(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Message", reflect.TypeOf((*MockDataCallback)(nil).Message), arg0)
}

// Trailer mocks base method.
func (m *MockDataCallback) Trailer(arg0 map[string]*prototool.Values) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockDataCallbackMockRecorder) Trailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockDataCallback)(nil).Trailer), arg0)
}

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// KasUrlsByAgentId mocks base method.
func (m *MockQuerier) KasUrlsByAgentId(arg0 context.Context, arg1 int64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KasUrlsByAgentId", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KasUrlsByAgentId indicates an expected call of KasUrlsByAgentId.
func (mr *MockQuerierMockRecorder) KasUrlsByAgentId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KasUrlsByAgentId", reflect.TypeOf((*MockQuerier)(nil).KasUrlsByAgentId), arg0, arg1)
}

// MockTracker is a mock of Tracker interface.
type MockTracker struct {
	ctrl     *gomock.Controller
	recorder *MockTrackerMockRecorder
}

// MockTrackerMockRecorder is the mock recorder for MockTracker.
type MockTrackerMockRecorder struct {
	mock *MockTracker
}

// NewMockTracker creates a new mock instance.
func NewMockTracker(ctrl *gomock.Controller) *MockTracker {
	mock := &MockTracker{ctrl: ctrl}
	mock.recorder = &MockTrackerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTracker) EXPECT() *MockTrackerMockRecorder {
	return m.recorder
}

// GC mocks base method.
func (m *MockTracker) GC() func(context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GC")
	ret0, _ := ret[0].(func(context.Context) (int, error))
	return ret0
}

// GC indicates an expected call of GC.
func (mr *MockTrackerMockRecorder) GC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GC", reflect.TypeOf((*MockTracker)(nil).GC))
}

// KasUrlsByAgentId mocks base method.
func (m *MockTracker) KasUrlsByAgentId(arg0 context.Context, arg1 int64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KasUrlsByAgentId", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KasUrlsByAgentId indicates an expected call of KasUrlsByAgentId.
func (mr *MockTrackerMockRecorder) KasUrlsByAgentId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KasUrlsByAgentId", reflect.TypeOf((*MockTracker)(nil).KasUrlsByAgentId), arg0, arg1)
}

// Refresh mocks base method.
func (m *MockTracker) Refresh(arg0 context.Context, arg1 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockTrackerMockRecorder) Refresh(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockTracker)(nil).Refresh), arg0, arg1)
}

// RegisterTunnel mocks base method.
func (m *MockTracker) RegisterTunnel(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterTunnel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterTunnel indicates an expected call of RegisterTunnel.
func (mr *MockTrackerMockRecorder) RegisterTunnel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterTunnel", reflect.TypeOf((*MockTracker)(nil).RegisterTunnel), arg0, arg1)
}

// UnregisterTunnel mocks base method.
func (m *MockTracker) UnregisterTunnel(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterTunnel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterTunnel indicates an expected call of UnregisterTunnel.
func (mr *MockTrackerMockRecorder) UnregisterTunnel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterTunnel", reflect.TypeOf((*MockTracker)(nil).UnregisterTunnel), arg0, arg1)
}
