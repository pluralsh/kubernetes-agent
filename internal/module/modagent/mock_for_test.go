// Code generated by MockGen. DO NOT EDIT.
// Source: worker_manager.go

// Package modagent is a generated GoMock package.
package modagent

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	agentcfg "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	proto "google.golang.org/protobuf/proto"
)

// MockWorkSource is a mock of WorkSource interface.
type MockWorkSource[C proto.Message] struct {
	ctrl     *gomock.Controller
	recorder *MockWorkSourceMockRecorder[C]
}

// MockWorkSourceMockRecorder is the mock recorder for MockWorkSource.
type MockWorkSourceMockRecorder[C proto.Message] struct {
	mock *MockWorkSource[C]
}

// NewMockWorkSource creates a new mock instance.
func NewMockWorkSource[C proto.Message](ctrl *gomock.Controller) *MockWorkSource[C] {
	mock := &MockWorkSource[C]{ctrl: ctrl}
	mock.recorder = &MockWorkSourceMockRecorder[C]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkSource[C]) EXPECT() *MockWorkSourceMockRecorder[C] {
	return m.recorder
}

// Configuration mocks base method.
func (m *MockWorkSource[C]) Configuration() C {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Configuration")
	ret0, _ := ret[0].(C)
	return ret0
}

// Configuration indicates an expected call of Configuration.
func (mr *MockWorkSourceMockRecorder[C]) Configuration() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configuration", reflect.TypeOf((*MockWorkSource[C])(nil).Configuration))
}

// ID mocks base method.
func (m *MockWorkSource[C]) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID.
func (mr *MockWorkSourceMockRecorder[C]) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockWorkSource[C])(nil).ID))
}

// MockWorkerFactory is a mock of WorkerFactory interface.
type MockWorkerFactory[C proto.Message] struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerFactoryMockRecorder[C]
}

// MockWorkerFactoryMockRecorder is the mock recorder for MockWorkerFactory.
type MockWorkerFactoryMockRecorder[C proto.Message] struct {
	mock *MockWorkerFactory[C]
}

// NewMockWorkerFactory creates a new mock instance.
func NewMockWorkerFactory[C proto.Message](ctrl *gomock.Controller) *MockWorkerFactory[C] {
	mock := &MockWorkerFactory[C]{ctrl: ctrl}
	mock.recorder = &MockWorkerFactoryMockRecorder[C]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorkerFactory[C]) EXPECT() *MockWorkerFactoryMockRecorder[C] {
	return m.recorder
}

// New mocks base method.
func (m *MockWorkerFactory[C]) New(agentId int64, source WorkSource[C]) Worker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", agentId, source)
	ret0, _ := ret[0].(Worker)
	return ret0
}

// New indicates an expected call of New.
func (mr *MockWorkerFactoryMockRecorder[C]) New(agentId, source interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockWorkerFactory[C])(nil).New), agentId, source)
}

// SourcesFromConfiguration mocks base method.
func (m *MockWorkerFactory[C]) SourcesFromConfiguration(arg0 *agentcfg.AgentConfiguration) []WorkSource[C] {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SourcesFromConfiguration", arg0)
	ret0, _ := ret[0].([]WorkSource[C])
	return ret0
}

// SourcesFromConfiguration indicates an expected call of SourcesFromConfiguration.
func (mr *MockWorkerFactoryMockRecorder[C]) SourcesFromConfiguration(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SourcesFromConfiguration", reflect.TypeOf((*MockWorkerFactory[C])(nil).SourcesFromConfiguration), arg0)
}

// MockWorker is a mock of Worker interface.
type MockWorker struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerMockRecorder
}

// MockWorkerMockRecorder is the mock recorder for MockWorker.
type MockWorkerMockRecorder struct {
	mock *MockWorker
}

// NewMockWorker creates a new mock instance.
func NewMockWorker(ctrl *gomock.Controller) *MockWorker {
	mock := &MockWorker{ctrl: ctrl}
	mock.recorder = &MockWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorker) EXPECT() *MockWorkerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockWorker) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockWorkerMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockWorker)(nil).Run), arg0)
}
