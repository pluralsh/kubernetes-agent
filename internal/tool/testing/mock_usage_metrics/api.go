// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics (interfaces: UsageTrackerInterface,Counter,UniqueCounter)
//
// Generated by this command:
//
//	mockgen -destination api.go gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics UsageTrackerInterface,Counter,UniqueCounter
//
// Package mock_usage_metrics is a generated GoMock package.
package mock_usage_metrics

import (
	reflect "reflect"

	usage_metrics "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/usage_metrics"
	gomock "go.uber.org/mock/gomock"
)

// MockUsageTrackerInterface is a mock of UsageTrackerInterface interface.
type MockUsageTrackerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUsageTrackerInterfaceMockRecorder
}

// MockUsageTrackerInterfaceMockRecorder is the mock recorder for MockUsageTrackerInterface.
type MockUsageTrackerInterfaceMockRecorder struct {
	mock *MockUsageTrackerInterface
}

// NewMockUsageTrackerInterface creates a new mock instance.
func NewMockUsageTrackerInterface(ctrl *gomock.Controller) *MockUsageTrackerInterface {
	mock := &MockUsageTrackerInterface{ctrl: ctrl}
	mock.recorder = &MockUsageTrackerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsageTrackerInterface) EXPECT() *MockUsageTrackerInterfaceMockRecorder {
	return m.recorder
}

// CloneUsageData mocks base method.
func (m *MockUsageTrackerInterface) CloneUsageData() *usage_metrics.UsageData {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloneUsageData")
	ret0, _ := ret[0].(*usage_metrics.UsageData)
	return ret0
}

// CloneUsageData indicates an expected call of CloneUsageData.
func (mr *MockUsageTrackerInterfaceMockRecorder) CloneUsageData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloneUsageData", reflect.TypeOf((*MockUsageTrackerInterface)(nil).CloneUsageData))
}

// RegisterCounter mocks base method.
func (m *MockUsageTrackerInterface) RegisterCounter(arg0 string) usage_metrics.Counter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCounter", arg0)
	ret0, _ := ret[0].(usage_metrics.Counter)
	return ret0
}

// RegisterCounter indicates an expected call of RegisterCounter.
func (mr *MockUsageTrackerInterfaceMockRecorder) RegisterCounter(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCounter", reflect.TypeOf((*MockUsageTrackerInterface)(nil).RegisterCounter), arg0)
}

// RegisterUniqueCounter mocks base method.
func (m *MockUsageTrackerInterface) RegisterUniqueCounter(arg0 string) usage_metrics.UniqueCounter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUniqueCounter", arg0)
	ret0, _ := ret[0].(usage_metrics.UniqueCounter)
	return ret0
}

// RegisterUniqueCounter indicates an expected call of RegisterUniqueCounter.
func (mr *MockUsageTrackerInterfaceMockRecorder) RegisterUniqueCounter(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUniqueCounter", reflect.TypeOf((*MockUsageTrackerInterface)(nil).RegisterUniqueCounter), arg0)
}

// Subtract mocks base method.
func (m *MockUsageTrackerInterface) Subtract(arg0 *usage_metrics.UsageData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subtract", arg0)
}

// Subtract indicates an expected call of Subtract.
func (mr *MockUsageTrackerInterfaceMockRecorder) Subtract(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subtract", reflect.TypeOf((*MockUsageTrackerInterface)(nil).Subtract), arg0)
}

// MockCounter is a mock of Counter interface.
type MockCounter struct {
	ctrl     *gomock.Controller
	recorder *MockCounterMockRecorder
}

// MockCounterMockRecorder is the mock recorder for MockCounter.
type MockCounterMockRecorder struct {
	mock *MockCounter
}

// NewMockCounter creates a new mock instance.
func NewMockCounter(ctrl *gomock.Controller) *MockCounter {
	mock := &MockCounter{ctrl: ctrl}
	mock.recorder = &MockCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCounter) EXPECT() *MockCounterMockRecorder {
	return m.recorder
}

// Inc mocks base method.
func (m *MockCounter) Inc() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Inc")
}

// Inc indicates an expected call of Inc.
func (mr *MockCounterMockRecorder) Inc() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Inc", reflect.TypeOf((*MockCounter)(nil).Inc))
}

// MockUniqueCounter is a mock of UniqueCounter interface.
type MockUniqueCounter struct {
	ctrl     *gomock.Controller
	recorder *MockUniqueCounterMockRecorder
}

// MockUniqueCounterMockRecorder is the mock recorder for MockUniqueCounter.
type MockUniqueCounterMockRecorder struct {
	mock *MockUniqueCounter
}

// NewMockUniqueCounter creates a new mock instance.
func NewMockUniqueCounter(ctrl *gomock.Controller) *MockUniqueCounter {
	mock := &MockUniqueCounter{ctrl: ctrl}
	mock.recorder = &MockUniqueCounterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUniqueCounter) EXPECT() *MockUniqueCounterMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockUniqueCounter) Add(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockUniqueCounterMockRecorder) Add(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockUniqueCounter)(nil).Add), arg0)
}
