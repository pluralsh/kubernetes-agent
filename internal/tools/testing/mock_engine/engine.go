// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/argoproj/gitops-engine/pkg/engine (interfaces: GitOpsEngine)

// Package mock_engine is a generated GoMock package.
package mock_engine

import (
	context "context"
	reflect "reflect"

	cache "github.com/argoproj/gitops-engine/pkg/cache"
	engine "github.com/argoproj/gitops-engine/pkg/engine"
	sync "github.com/argoproj/gitops-engine/pkg/sync"
	common "github.com/argoproj/gitops-engine/pkg/sync/common"
	gomock "github.com/golang/mock/gomock"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// MockGitOpsEngine is a mock of GitOpsEngine interface
type MockGitOpsEngine struct {
	ctrl     *gomock.Controller
	recorder *MockGitOpsEngineMockRecorder
}

// MockGitOpsEngineMockRecorder is the mock recorder for MockGitOpsEngine
type MockGitOpsEngineMockRecorder struct {
	mock *MockGitOpsEngine
}

// NewMockGitOpsEngine creates a new mock instance
func NewMockGitOpsEngine(ctrl *gomock.Controller) *MockGitOpsEngine {
	mock := &MockGitOpsEngine{ctrl: ctrl}
	mock.recorder = &MockGitOpsEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGitOpsEngine) EXPECT() *MockGitOpsEngineMockRecorder {
	return m.recorder
}

// Run mocks base method
func (m *MockGitOpsEngine) Run() (engine.StopFunc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(engine.StopFunc)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Run indicates an expected call of Run
func (mr *MockGitOpsEngineMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockGitOpsEngine)(nil).Run))
}

// Sync mocks base method
func (m *MockGitOpsEngine) Sync(arg0 context.Context, arg1 []*unstructured.Unstructured, arg2 func(*cache.Resource) bool, arg3, arg4 string, arg5 ...sync.SyncOpt) ([]common.ResourceSyncResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3, arg4}
	for _, a := range arg5 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Sync", varargs...)
	ret0, _ := ret[0].([]common.ResourceSyncResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Sync indicates an expected call of Sync
func (mr *MockGitOpsEngineMockRecorder) Sync(arg0, arg1, arg2, arg3, arg4 interface{}, arg5 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3, arg4}, arg5...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sync", reflect.TypeOf((*MockGitOpsEngine)(nil).Sync), varargs...)
}
