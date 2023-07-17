// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd/kas/kasapp (interfaces: SentryHub)

// Package kasapp is a generated GoMock package.
package kasapp

import (
	reflect "reflect"

	sentry "github.com/getsentry/sentry-go"
	gomock "go.uber.org/mock/gomock"
)

// MockSentryHub is a mock of SentryHub interface.
type MockSentryHub struct {
	ctrl     *gomock.Controller
	recorder *MockSentryHubMockRecorder
}

// MockSentryHubMockRecorder is the mock recorder for MockSentryHub.
type MockSentryHubMockRecorder struct {
	mock *MockSentryHub
}

// NewMockSentryHub creates a new mock instance.
func NewMockSentryHub(ctrl *gomock.Controller) *MockSentryHub {
	mock := &MockSentryHub{ctrl: ctrl}
	mock.recorder = &MockSentryHubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSentryHub) EXPECT() *MockSentryHubMockRecorder {
	return m.recorder
}

// CaptureEvent mocks base method.
func (m *MockSentryHub) CaptureEvent(arg0 *sentry.Event) *sentry.EventID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CaptureEvent", arg0)
	ret0, _ := ret[0].(*sentry.EventID)
	return ret0
}

// CaptureEvent indicates an expected call of CaptureEvent.
func (mr *MockSentryHubMockRecorder) CaptureEvent(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CaptureEvent", reflect.TypeOf((*MockSentryHub)(nil).CaptureEvent), arg0)
}
