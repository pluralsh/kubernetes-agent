// Code generated by MockGen. DO NOT EDIT.
// Source: reconcile_trigger.go
//
// Generated by this command:
//
//	mockgen -typed -source reconcile_trigger.go -destination mock_reconciler_trigger_for_test.go -package agent reconcileTrigger
//
// Package agent is a generated GoMock package.
package agent

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockreconcileTrigger is a mock of reconcileTrigger interface.
type MockreconcileTrigger struct {
	ctrl     *gomock.Controller
	recorder *MockreconcileTriggerMockRecorder
}

// MockreconcileTriggerMockRecorder is the mock recorder for MockreconcileTrigger.
type MockreconcileTriggerMockRecorder struct {
	mock *MockreconcileTrigger
}

// NewMockreconcileTrigger creates a new mock instance.
func NewMockreconcileTrigger(ctrl *gomock.Controller) *MockreconcileTrigger {
	mock := &MockreconcileTrigger{ctrl: ctrl}
	mock.recorder = &MockreconcileTriggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockreconcileTrigger) EXPECT() *MockreconcileTriggerMockRecorder {
	return m.recorder
}

// reconcile mocks base method.
func (m *MockreconcileTrigger) reconcile(ctx context.Context, webhookPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "reconcile", ctx, webhookPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// reconcile indicates an expected call of reconcile.
func (mr *MockreconcileTriggerMockRecorder) reconcile(ctx, webhookPath any) *reconcileTriggerreconcileCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "reconcile", reflect.TypeOf((*MockreconcileTrigger)(nil).reconcile), ctx, webhookPath)
	return &reconcileTriggerreconcileCall{Call: call}
}

// reconcileTriggerreconcileCall wrap *gomock.Call
type reconcileTriggerreconcileCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *reconcileTriggerreconcileCall) Return(arg0 error) *reconcileTriggerreconcileCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *reconcileTriggerreconcileCall) Do(f func(context.Context, string) error) *reconcileTriggerreconcileCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *reconcileTriggerreconcileCall) DoAndReturn(f func(context.Context, string) error) *reconcileTriggerreconcileCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
