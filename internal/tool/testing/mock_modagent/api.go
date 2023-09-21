// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent (interfaces: Api,Factory,Module)
//
// Generated by this command:
//
//	mockgen -destination api.go -package mock_modagent gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent Api,Factory,Module
//
// Package mock_modagent is a generated GoMock package.
package mock_modagent

import (
	context "context"
	url "net/url"
	reflect "reflect"

	modagent "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	modshared "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modshared"
	agentcfg "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockApi is a mock of Api interface.
type MockApi struct {
	ctrl     *gomock.Controller
	recorder *MockApiMockRecorder
}

// MockApiMockRecorder is the mock recorder for MockApi.
type MockApiMockRecorder struct {
	mock *MockApi
}

// NewMockApi creates a new mock instance.
func NewMockApi(ctrl *gomock.Controller) *MockApi {
	mock := &MockApi{ctrl: ctrl}
	mock.recorder = &MockApiMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApi) EXPECT() *MockApiMockRecorder {
	return m.recorder
}

// GetAgentId mocks base method.
func (m *MockApi) GetAgentId(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAgentId", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAgentId indicates an expected call of GetAgentId.
func (mr *MockApiMockRecorder) GetAgentId(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgentId", reflect.TypeOf((*MockApi)(nil).GetAgentId), arg0)
}

// GetGitLabExternalUrl mocks base method.
func (m *MockApi) GetGitLabExternalUrl(arg0 context.Context) (url.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGitLabExternalUrl", arg0)
	ret0, _ := ret[0].(url.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGitLabExternalUrl indicates an expected call of GetGitLabExternalUrl.
func (mr *MockApiMockRecorder) GetGitLabExternalUrl(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitLabExternalUrl", reflect.TypeOf((*MockApi)(nil).GetGitLabExternalUrl), arg0)
}

// HandleProcessingError mocks base method.
func (m *MockApi) HandleProcessingError(arg0 context.Context, arg1 *zap.Logger, arg2 int64, arg3 string, arg4 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleProcessingError", arg0, arg1, arg2, arg3, arg4)
}

// HandleProcessingError indicates an expected call of HandleProcessingError.
func (mr *MockApiMockRecorder) HandleProcessingError(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProcessingError", reflect.TypeOf((*MockApi)(nil).HandleProcessingError), arg0, arg1, arg2, arg3, arg4)
}

// MakeGitLabRequest mocks base method.
func (m *MockApi) MakeGitLabRequest(arg0 context.Context, arg1 string, arg2 ...modagent.GitLabRequestOption) (*modagent.GitLabResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "MakeGitLabRequest", varargs...)
	ret0, _ := ret[0].(*modagent.GitLabResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MakeGitLabRequest indicates an expected call of MakeGitLabRequest.
func (mr *MockApiMockRecorder) MakeGitLabRequest(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeGitLabRequest", reflect.TypeOf((*MockApi)(nil).MakeGitLabRequest), varargs...)
}

// TryGetAgentId mocks base method.
func (m *MockApi) TryGetAgentId() (int64, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryGetAgentId")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// TryGetAgentId indicates an expected call of TryGetAgentId.
func (mr *MockApiMockRecorder) TryGetAgentId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryGetAgentId", reflect.TypeOf((*MockApi)(nil).TryGetAgentId))
}

// MockFactory is a mock of Factory interface.
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory.
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance.
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// IsProducingLeaderModules mocks base method.
func (m *MockFactory) IsProducingLeaderModules() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsProducingLeaderModules")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsProducingLeaderModules indicates an expected call of IsProducingLeaderModules.
func (mr *MockFactoryMockRecorder) IsProducingLeaderModules() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsProducingLeaderModules", reflect.TypeOf((*MockFactory)(nil).IsProducingLeaderModules))
}

// Name mocks base method.
func (m *MockFactory) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockFactoryMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFactory)(nil).Name))
}

// New mocks base method.
func (m *MockFactory) New(arg0 *modagent.Config) (modagent.Module, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0)
	ret0, _ := ret[0].(modagent.Module)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockFactoryMockRecorder) New(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockFactory)(nil).New), arg0)
}

// StartStopPhase mocks base method.
func (m *MockFactory) StartStopPhase() modshared.ModuleStartStopPhase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartStopPhase")
	ret0, _ := ret[0].(modshared.ModuleStartStopPhase)
	return ret0
}

// StartStopPhase indicates an expected call of StartStopPhase.
func (mr *MockFactoryMockRecorder) StartStopPhase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartStopPhase", reflect.TypeOf((*MockFactory)(nil).StartStopPhase))
}

// MockModule is a mock of Module interface.
type MockModule struct {
	ctrl     *gomock.Controller
	recorder *MockModuleMockRecorder
}

// MockModuleMockRecorder is the mock recorder for MockModule.
type MockModuleMockRecorder struct {
	mock *MockModule
}

// NewMockModule creates a new mock instance.
func NewMockModule(ctrl *gomock.Controller) *MockModule {
	mock := &MockModule{ctrl: ctrl}
	mock.recorder = &MockModuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockModule) EXPECT() *MockModuleMockRecorder {
	return m.recorder
}

// DefaultAndValidateConfiguration mocks base method.
func (m *MockModule) DefaultAndValidateConfiguration(arg0 *agentcfg.AgentConfiguration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultAndValidateConfiguration", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DefaultAndValidateConfiguration indicates an expected call of DefaultAndValidateConfiguration.
func (mr *MockModuleMockRecorder) DefaultAndValidateConfiguration(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultAndValidateConfiguration", reflect.TypeOf((*MockModule)(nil).DefaultAndValidateConfiguration), arg0)
}

// Name mocks base method.
func (m *MockModule) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockModuleMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockModule)(nil).Name))
}

// Run mocks base method.
func (m *MockModule) Run(arg0 context.Context, arg1 <-chan *agentcfg.AgentConfiguration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockModuleMockRecorder) Run(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockModule)(nil).Run), arg0, arg1)
}
