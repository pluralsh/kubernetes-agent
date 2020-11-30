// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/usage_metrics (interfaces: UsageTrackerInterface,Counter)

// Package mock_usage_metrics is a generated GoMock package.
package mock_usage_metrics

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	usage_metrics "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/usage_metrics"
)

// MockUsageTrackerInterface is a mock of UsageTrackerInterface interface
type MockUsageTrackerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUsageTrackerInterfaceMockRecorder
}

// MockUsageTrackerInterfaceMockRecorder is the mock recorder for MockUsageTrackerInterface
type MockUsageTrackerInterfaceMockRecorder struct {
	mock *MockUsageTrackerInterface
}

// NewMockUsageTrackerInterface creates a new mock instance
func NewMockUsageTrackerInterface(ctrl *gomock.Controller) *MockUsageTrackerInterface {
	mock := &MockUsageTrackerInterface{ctrl: ctrl}
	mock.recorder = &MockUsageTrackerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsageTrackerInterface) EXPECT() *MockUsageTrackerInterfaceMockRecorder {
	return m.recorder
}

// CloneUsageData mocks base method
func (m *MockUsageTrackerInterface) CloneUsageData() (*usage_metrics.UsageData, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloneUsageData")
	ret0, _ := ret[0].(*usage_metrics.UsageData)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// CloneUsageData indicates an expected call of CloneUsageData
func (mr *MockUsageTrackerInterfaceMockRecorder) CloneUsageData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloneUsageData", reflect.TypeOf((*MockUsageTrackerInterface)(nil).CloneUsageData))
}

// RegisterCounter mocks base method
func (m *MockUsageTrackerInterface) RegisterCounter(arg0 string) usage_metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCounter", arg0)
	ret0, _ := ret[0].(usage_metrics.Counter)
	return ret0
}

// RegisterCounter indicates an expected call of RegisterCounter
func (mr *MockUsageTrackerInterfaceMockRecorder) RegisterCounter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCounter", reflect.TypeOf((*MockUsageTrackerInterface)(nil).RegisterCounter), arg0)
}

// Subtract mocks base method
func (m *MockUsageTrackerInterface) Subtract(arg0 *usage_metrics.UsageData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subtract", arg0)
}

// Subtract indicates an expected call of Subtract
func (mr *MockUsageTrackerInterfaceMockRecorder) Subtract(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subtract", reflect.TypeOf((*MockUsageTrackerInterface)(nil).Subtract), arg0)
}

// MockCounter is a mock of Counter interface
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Inc mocks base method
func (m *MockCounter) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc
func (mr *MockCounterMockRecorder) Inc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
}