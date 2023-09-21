// Code generated by MockGen. DO NOT EDIT.
// Source: ../../redistool/expiring_hash_api_set_builder.go
//
// Generated by this command:
//
//	mockgen -typed -source ../../redistool/expiring_hash_api_set_builder.go -destination expiring_hash_api_set_builder.go -package mock_redis
//
// Package mock_redis is a generated GoMock package.
package mock_redis

import (
	context "context"
	reflect "reflect"
	time "time"

	redistool "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/redistool"
	gomock "go.uber.org/mock/gomock"
)

// MockSetBuilder is a mock of SetBuilder interface.
type MockSetBuilder[K1 any, K2 any] struct {
	ctrl     *gomock.Controller
	recorder *MockSetBuilderMockRecorder[K1, K2]
}

// MockSetBuilderMockRecorder is the mock recorder for MockSetBuilder.
type MockSetBuilderMockRecorder[K1 any, K2 any] struct {
	mock *MockSetBuilder[K1, K2]
}

// NewMockSetBuilder creates a new mock instance.
func NewMockSetBuilder[K1 any, K2 any](ctrl *gomock.Controller) *MockSetBuilder[K1, K2] {
	mock := &MockSetBuilder[K1, K2]{ctrl: ctrl}
	mock.recorder = &MockSetBuilderMockRecorder[K1, K2]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSetBuilder[K1, K2]) EXPECT() *MockSetBuilderMockRecorder[K1, K2] {
	return m.recorder
}

// Do mocks base method.
func (m *MockSetBuilder[K1, K2]) Do(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Do indicates an expected call of Do.
func (mr *MockSetBuilderMockRecorder[K1, K2]) Do(arg0 any) *SetBuilderDoCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockSetBuilder[K1, K2])(nil).Do), arg0)
	return &SetBuilderDoCall[K1, K2]{Call: call}
}

// SetBuilderDoCall wrap *gomock.Call
type SetBuilderDoCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SetBuilderDoCall[K1, K2]) Return(arg0 error) *SetBuilderDoCall[K1, K2] {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SetBuilderDoCall[K1, K2]) Do(f func(context.Context) error) *SetBuilderDoCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SetBuilderDoCall[K1, K2]) DoAndReturn(f func(context.Context) error) *SetBuilderDoCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Set mocks base method.
func (m *MockSetBuilder[K1, K2]) Set(key K1, ttl time.Duration, kvs ...redistool.BuilderKV[K2]) {
	m.ctrl.T.Helper()
	varargs := []any{key, ttl}
	for _, a := range kvs {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Set", varargs...)
}

// Set indicates an expected call of Set.
func (mr *MockSetBuilderMockRecorder[K1, K2]) Set(key, ttl any, kvs ...any) *SetBuilderSetCall[K1, K2] {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{key, ttl}, kvs...)
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockSetBuilder[K1, K2])(nil).Set), varargs...)
	return &SetBuilderSetCall[K1, K2]{Call: call}
}

// SetBuilderSetCall wrap *gomock.Call
type SetBuilderSetCall[K1 any, K2 any] struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *SetBuilderSetCall[K1, K2]) Return() *SetBuilderSetCall[K1, K2] {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *SetBuilderSetCall[K1, K2]) Do(f func(K1, time.Duration, ...redistool.BuilderKV[K2])) *SetBuilderSetCall[K1, K2] {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *SetBuilderSetCall[K1, K2]) DoAndReturn(f func(K1, time.Duration, ...redistool.BuilderKV[K2])) *SetBuilderSetCall[K1, K2] {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
