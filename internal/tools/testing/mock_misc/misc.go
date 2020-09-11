// Code generated by MockGen. DO NOT EDIT.
// Source: io (interfaces: Closer)

// Package mock_misc is a generated GoMock package.
package mock_misc

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCloser is a mock of Closer interface
type MockCloser struct {
	ctrl     *gomock.Controller
	recorder *MockCloserMockRecorder
}

// MockCloserMockRecorder is the mock recorder for MockCloser
type MockCloserMockRecorder struct {
	mock *MockCloser
}

// NewMockCloser creates a new mock instance
func NewMockCloser(ctrl *gomock.Controller) *MockCloser {
	mock := &MockCloser{ctrl: ctrl}
	mock.recorder = &MockCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCloser) EXPECT() *MockCloserMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCloser)(nil).Close))
}
