// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/gitops/agent/chartops (interfaces: Helm)

// Package chartops is a generated GoMock package.
package chartops

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	chart "helm.sh/helm/v3/pkg/chart"
	release "helm.sh/helm/v3/pkg/release"
)

// MockHelm is a mock of Helm interface.
type MockHelm struct {
	ctrl     *gomock.Controller
	recorder *MockHelmMockRecorder
}

// MockHelmMockRecorder is the mock recorder for MockHelm.
type MockHelmMockRecorder struct {
	mock *MockHelm
}

// NewMockHelm creates a new mock instance.
func NewMockHelm(ctrl *gomock.Controller) *MockHelm {
	mock := &MockHelm{ctrl: ctrl}
	mock.recorder = &MockHelmMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHelm) EXPECT() *MockHelmMockRecorder {
	return m.recorder
}

// History mocks base method.
func (m *MockHelm) History(arg0 string) ([]*release.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "History", arg0)
	ret0, _ := ret[0].([]*release.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// History indicates an expected call of History.
func (mr *MockHelmMockRecorder) History(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "History", reflect.TypeOf((*MockHelm)(nil).History), arg0)
}

// Install mocks base method.
func (m *MockHelm) Install(arg0 context.Context, arg1 *chart.Chart, arg2 ChartValues, arg3 InstallConfig) (*release.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*release.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Install indicates an expected call of Install.
func (mr *MockHelmMockRecorder) Install(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockHelm)(nil).Install), arg0, arg1, arg2, arg3)
}

// Upgrade mocks base method.
func (m *MockHelm) Upgrade(arg0 context.Context, arg1 string, arg2 *chart.Chart, arg3 ChartValues, arg4 UpgradeConfig) (*release.Release, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upgrade", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*release.Release)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upgrade indicates an expected call of Upgrade.
func (mr *MockHelmMockRecorder) Upgrade(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upgrade", reflect.TypeOf((*MockHelm)(nil).Upgrade), arg0, arg1, arg2, arg3, arg4)
}
