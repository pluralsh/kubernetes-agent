// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/gitops/agent (interfaces: ApplierFactory,GitopsWorkerFactory,GitopsWorker,Applier)

// Package agent is a generated GoMock package.
package agent

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	agentcfg "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/pkg/agentcfg"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apply "sigs.k8s.io/cli-utils/pkg/apply"
	event "sigs.k8s.io/cli-utils/pkg/apply/event"
	inventory "sigs.k8s.io/cli-utils/pkg/inventory"
)

// MockApplierFactory is a mock of ApplierFactory interface.
type MockApplierFactory struct {
	ctrl     *gomock.Controller
	recorder *MockApplierFactoryMockRecorder
}

// MockApplierFactoryMockRecorder is the mock recorder for MockApplierFactory.
type MockApplierFactoryMockRecorder struct {
	mock *MockApplierFactory
}

// NewMockApplierFactory creates a new mock instance.
func NewMockApplierFactory(ctrl *gomock.Controller) *MockApplierFactory {
	mock := &MockApplierFactory{ctrl: ctrl}
	mock.recorder = &MockApplierFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplierFactory) EXPECT() *MockApplierFactoryMockRecorder {
	return m.recorder
}

// New mocks base method.
func (m *MockApplierFactory) New() Applier {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New")
	ret0, _ := ret[0].(Applier)
	return ret0
}

// New indicates an expected call of New.
func (mr *MockApplierFactoryMockRecorder) New() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockApplierFactory)(nil).New))
}

// MockGitopsWorkerFactory is a mock of GitopsWorkerFactory interface.
type MockGitopsWorkerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockGitopsWorkerFactoryMockRecorder
}

// MockGitopsWorkerFactoryMockRecorder is the mock recorder for MockGitopsWorkerFactory.
type MockGitopsWorkerFactoryMockRecorder struct {
	mock *MockGitopsWorkerFactory
}

// NewMockGitopsWorkerFactory creates a new mock instance.
func NewMockGitopsWorkerFactory(ctrl *gomock.Controller) *MockGitopsWorkerFactory {
	mock := &MockGitopsWorkerFactory{ctrl: ctrl}
	mock.recorder = &MockGitopsWorkerFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitopsWorkerFactory) EXPECT() *MockGitopsWorkerFactoryMockRecorder {
	return m.recorder
}

// New mocks base method.
func (m *MockGitopsWorkerFactory) New(arg0 int64, arg1 *agentcfg.ManifestProjectCF) GitopsWorker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0, arg1)
	ret0, _ := ret[0].(GitopsWorker)
	return ret0
}

// New indicates an expected call of New.
func (mr *MockGitopsWorkerFactoryMockRecorder) New(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockGitopsWorkerFactory)(nil).New), arg0, arg1)
}

// MockGitopsWorker is a mock of GitopsWorker interface.
type MockGitopsWorker struct {
	ctrl     *gomock.Controller
	recorder *MockGitopsWorkerMockRecorder
}

// MockGitopsWorkerMockRecorder is the mock recorder for MockGitopsWorker.
type MockGitopsWorkerMockRecorder struct {
	mock *MockGitopsWorker
}

// NewMockGitopsWorker creates a new mock instance.
func NewMockGitopsWorker(ctrl *gomock.Controller) *MockGitopsWorker {
	mock := &MockGitopsWorker{ctrl: ctrl}
	mock.recorder = &MockGitopsWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitopsWorker) EXPECT() *MockGitopsWorkerMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockGitopsWorker) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockGitopsWorkerMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockGitopsWorker)(nil).Run), arg0)
}

// MockApplier is a mock of Applier interface.
type MockApplier struct {
	ctrl     *gomock.Controller
	recorder *MockApplierMockRecorder
}

// MockApplierMockRecorder is the mock recorder for MockApplier.
type MockApplierMockRecorder struct {
	mock *MockApplier
}

// NewMockApplier creates a new mock instance.
func NewMockApplier(ctrl *gomock.Controller) *MockApplier {
	mock := &MockApplier{ctrl: ctrl}
	mock.recorder = &MockApplierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApplier) EXPECT() *MockApplierMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockApplier) Initialize() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockApplierMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockApplier)(nil).Initialize))
}

// Run mocks base method.
func (m *MockApplier) Run(arg0 context.Context, arg1 inventory.InventoryInfo, arg2 []*unstructured.Unstructured, arg3 apply.Options) <-chan event.Event {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(<-chan event.Event)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockApplierMockRecorder) Run(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockApplier)(nil).Run), arg0, arg1, arg2, arg3)
}
