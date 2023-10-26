// Code generated by MockGen. DO NOT EDIT.
// Source: ../../redistool/expiring_hash.go
//
// Generated by this command:
//
//	mockgen -typed -source ../../redistool/expiring_hash.go -destination expiring_hash.go -package mock_redis
//
// Package mock_redis is a generated GoMock package.
package mock_redis

import (
	context "context"
	reflect "reflect"
	time "time"

	redistool "github.com/pluralsh/kuberentes-agent/pkg/tool/redistool"
	gomock "go.uber.org/mock/gomock"
)

// MockExpiringHash is a mock of ExpiringHash interface.
type MockExpiringHash[K1 any, K2 any] struct {
	ctrl     *gomock.Controller
	recorder *MockExpiringHashMockRecorder[K1, K2]
}

// MockExpiringHashMockRecorder is the mock recorder for MockExpiringHash.
type MockExpiringHashMockRecorder[K1 any, K2 any] struct {
	mock *MockExpiringHash[K1, K2]
}

// NewMockExpiringHash creates a new mock instance.
func NewMockExpiringHash[K1 any, K2 any](ctrl *gomock.Controller) *MockExpiringHash[K1, K2] {
	mock := &MockExpiringHash[K1, K2]{ctrl: ctrl}
	mock.recorder = &MockExpiringHashMockRecorder[K1, K2]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExpiringHash[K1, K2]) EXPECT() *MockExpiringHashMockRecorder[K1, K2] {
	return m.recorder
}

// Clear mocks base method.
func (m *MockExpiringHash[K1, K2]) Clear(arg0 context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clear", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Clear indicates an expected call of Clear.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Clear(arg0 any) *ExpiringHashClearCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clear", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Clear), arg0)
	return &ExpiringHashClearCall[K1, K2]{Call: call}
}

// ExpiringHashClearCall wrap *gomock.Call
type ExpiringHashClearCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashClearCall[K1, K2]) Return(arg0 int, arg1 error) *ExpiringHashClearCall[K1, K2] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashClearCall[K1, K2]) Do(f func(context.Context) (int, error)) *ExpiringHashClearCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashClearCall[K1, K2]) DoAndReturn(f func(context.Context) (int, error)) *ExpiringHashClearCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Forget mocks base method.
func (m *MockExpiringHash[K1, K2]) Forget(key K1, hashKey K2) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Forget", key, hashKey)
}

// Forget indicates an expected call of Forget.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Forget(key, hashKey any) *ExpiringHashForgetCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Forget", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Forget), key, hashKey)
	return &ExpiringHashForgetCall[K1, K2]{Call: call}
}

// ExpiringHashForgetCall wrap *gomock.Call
type ExpiringHashForgetCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashForgetCall[K1, K2]) Return() *ExpiringHashForgetCall[K1, K2] {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashForgetCall[K1, K2]) Do(f func(K1, K2)) *ExpiringHashForgetCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashForgetCall[K1, K2]) DoAndReturn(f func(K1, K2)) *ExpiringHashForgetCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// GC mocks base method.
func (m *MockExpiringHash[K1, K2]) GC() func(context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GC")
	ret0, _ := ret[0].(func(context.Context) (int, error))
	return ret0
}

// GC indicates an expected call of GC.
func (mr *MockExpiringHashMockRecorder[K1, K2]) GC() *ExpiringHashGCCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GC", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).GC))
	return &ExpiringHashGCCall[K1, K2]{Call: call}
}

// ExpiringHashGCCall wrap *gomock.Call
type ExpiringHashGCCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashGCCall[K1, K2]) Return(arg0 func(context.Context) (int, error)) *ExpiringHashGCCall[K1, K2] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashGCCall[K1, K2]) Do(f func() func(context.Context) (int, error)) *ExpiringHashGCCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashGCCall[K1, K2]) DoAndReturn(f func() func(context.Context) (int, error)) *ExpiringHashGCCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Len mocks base method.
func (m *MockExpiringHash[K1, K2]) Len(ctx context.Context, key K1) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Len indicates an expected call of Len.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Len(ctx, key any) *ExpiringHashLenCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Len), ctx, key)
	return &ExpiringHashLenCall[K1, K2]{Call: call}
}

// ExpiringHashLenCall wrap *gomock.Call
type ExpiringHashLenCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashLenCall[K1, K2]) Return(arg0 int64, arg1 error) *ExpiringHashLenCall[K1, K2] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashLenCall[K1, K2]) Do(f func(context.Context, K1) (int64, error)) *ExpiringHashLenCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashLenCall[K1, K2]) DoAndReturn(f func(context.Context, K1) (int64, error)) *ExpiringHashLenCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Refresh mocks base method.
func (m *MockExpiringHash[K1, K2]) Refresh(ctx context.Context, nextRefresh time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, nextRefresh)
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Refresh(ctx, nextRefresh any) *ExpiringHashRefreshCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Refresh), ctx, nextRefresh)
	return &ExpiringHashRefreshCall[K1, K2]{Call: call}
}

// ExpiringHashRefreshCall wrap *gomock.Call
type ExpiringHashRefreshCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashRefreshCall[K1, K2]) Return(arg0 error) *ExpiringHashRefreshCall[K1, K2] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashRefreshCall[K1, K2]) Do(f func(context.Context, time.Time) error) *ExpiringHashRefreshCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashRefreshCall[K1, K2]) DoAndReturn(f func(context.Context, time.Time) error) *ExpiringHashRefreshCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Scan mocks base method.
func (m *MockExpiringHash[K1, K2]) Scan(ctx context.Context, key K1, cb redistool.ScanCallback) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", ctx, key, cb)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scan indicates an expected call of Scan.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Scan(ctx, key, cb any) *ExpiringHashScanCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Scan), ctx, key, cb)
	return &ExpiringHashScanCall[K1, K2]{Call: call}
}

// ExpiringHashScanCall wrap *gomock.Call
type ExpiringHashScanCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashScanCall[K1, K2]) Return(arg0 int, arg1 error) *ExpiringHashScanCall[K1, K2] {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashScanCall[K1, K2]) Do(f func(context.Context, K1, redistool.ScanCallback) (int, error)) *ExpiringHashScanCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashScanCall[K1, K2]) DoAndReturn(f func(context.Context, K1, redistool.ScanCallback) (int, error)) *ExpiringHashScanCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Set mocks base method.
func (m *MockExpiringHash[K1, K2]) Set(ctx context.Context, key K1, hashKey K2, value []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, hashKey, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Set(ctx, key, hashKey, value any) *ExpiringHashSetCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Set), ctx, key, hashKey, value)
	return &ExpiringHashSetCall[K1, K2]{Call: call}
}

// ExpiringHashSetCall wrap *gomock.Call
type ExpiringHashSetCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashSetCall[K1, K2]) Return(arg0 error) *ExpiringHashSetCall[K1, K2] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashSetCall[K1, K2]) Do(f func(context.Context, K1, K2, []byte) error) *ExpiringHashSetCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashSetCall[K1, K2]) DoAndReturn(f func(context.Context, K1, K2, []byte) error) *ExpiringHashSetCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Unset mocks base method.
func (m *MockExpiringHash[K1, K2]) Unset(ctx context.Context, key K1, hashKey K2) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unset", ctx, key, hashKey)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unset indicates an expected call of Unset.
func (mr *MockExpiringHashMockRecorder[K1, K2]) Unset(ctx, key, hashKey any) *ExpiringHashUnsetCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unset", reflect.TypeOf((*MockExpiringHash[K1, K2])(nil).Unset), ctx, key, hashKey)
	return &ExpiringHashUnsetCall[K1, K2]{Call: call}
}

// ExpiringHashUnsetCall wrap *gomock.Call
type ExpiringHashUnsetCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ExpiringHashUnsetCall[K1, K2]) Return(arg0 error) *ExpiringHashUnsetCall[K1, K2] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ExpiringHashUnsetCall[K1, K2]) Do(f func(context.Context, K1, K2) error) *ExpiringHashUnsetCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ExpiringHashUnsetCall[K1, K2]) DoAndReturn(f func(context.Context, K1, K2) error) *ExpiringHashUnsetCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}