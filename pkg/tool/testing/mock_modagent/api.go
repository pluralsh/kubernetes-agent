// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pluralsh/kuberentes-agent/pkg/module/modagent (interfaces: Api,Factory,Module)
//
// Generated by this command:
//
//	mockgen -typed -destination api.go -package mock_modagent github.com/pluralsh/kuberentes-agent/pkg/module/modagent Api,Factory,Module
//
// Package mock_modagent is a generated GoMock package.
package mock_modagent

import (
	context "context"
	url "net/url"
	reflect "reflect"

	agentcfg "github.com/pluralsh/kuberentes-agent/pkg/agentcfg"
	modagent "github.com/pluralsh/kuberentes-agent/pkg/module/modagent"
	modshared "github.com/pluralsh/kuberentes-agent/pkg/module/modshared"
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
func (mr *MockApiMockRecorder) GetAgentId(arg0 any) *ApiGetAgentIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAgentId", reflect.TypeOf((*MockApi)(nil).GetAgentId), arg0)
	return &ApiGetAgentIdCall{Call: call}
}

// ApiGetAgentIdCall wrap *gomock.Call
type ApiGetAgentIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiGetAgentIdCall) Return(arg0 int64, arg1 error) *ApiGetAgentIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiGetAgentIdCall) Do(f func(context.Context) (int64, error)) *ApiGetAgentIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiGetAgentIdCall) DoAndReturn(f func(context.Context) (int64, error)) *ApiGetAgentIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockApiMockRecorder) GetGitLabExternalUrl(arg0 any) *ApiGetGitLabExternalUrlCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGitLabExternalUrl", reflect.TypeOf((*MockApi)(nil).GetGitLabExternalUrl), arg0)
	return &ApiGetGitLabExternalUrlCall{Call: call}
}

// ApiGetGitLabExternalUrlCall wrap *gomock.Call
type ApiGetGitLabExternalUrlCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiGetGitLabExternalUrlCall) Return(arg0 url.URL, arg1 error) *ApiGetGitLabExternalUrlCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiGetGitLabExternalUrlCall) Do(f func(context.Context) (url.URL, error)) *ApiGetGitLabExternalUrlCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiGetGitLabExternalUrlCall) DoAndReturn(f func(context.Context) (url.URL, error)) *ApiGetGitLabExternalUrlCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// HandleProcessingError mocks base method.
func (m *MockApi) HandleProcessingError(arg0 context.Context, arg1 *zap.Logger, arg2 int64, arg3 string, arg4 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleProcessingError", arg0, arg1, arg2, arg3, arg4)
}

// HandleProcessingError indicates an expected call of HandleProcessingError.
func (mr *MockApiMockRecorder) HandleProcessingError(arg0, arg1, arg2, arg3, arg4 any) *ApiHandleProcessingErrorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleProcessingError", reflect.TypeOf((*MockApi)(nil).HandleProcessingError), arg0, arg1, arg2, arg3, arg4)
	return &ApiHandleProcessingErrorCall{Call: call}
}

// ApiHandleProcessingErrorCall wrap *gomock.Call
type ApiHandleProcessingErrorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiHandleProcessingErrorCall) Return() *ApiHandleProcessingErrorCall {
	c.Call = c.Call.Return()
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiHandleProcessingErrorCall) Do(f func(context.Context, *zap.Logger, int64, string, error)) *ApiHandleProcessingErrorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiHandleProcessingErrorCall) DoAndReturn(f func(context.Context, *zap.Logger, int64, string, error)) *ApiHandleProcessingErrorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockApiMockRecorder) TryGetAgentId() *ApiTryGetAgentIdCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryGetAgentId", reflect.TypeOf((*MockApi)(nil).TryGetAgentId))
	return &ApiTryGetAgentIdCall{Call: call}
}

// ApiTryGetAgentIdCall wrap *gomock.Call
type ApiTryGetAgentIdCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ApiTryGetAgentIdCall) Return(arg0 int64, arg1 bool) *ApiTryGetAgentIdCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ApiTryGetAgentIdCall) Do(f func() (int64, bool)) *ApiTryGetAgentIdCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ApiTryGetAgentIdCall) DoAndReturn(f func() (int64, bool)) *ApiTryGetAgentIdCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) IsProducingLeaderModules() *FactoryIsProducingLeaderModulesCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsProducingLeaderModules", reflect.TypeOf((*MockFactory)(nil).IsProducingLeaderModules))
	return &FactoryIsProducingLeaderModulesCall{Call: call}
}

// FactoryIsProducingLeaderModulesCall wrap *gomock.Call
type FactoryIsProducingLeaderModulesCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryIsProducingLeaderModulesCall) Return(arg0 bool) *FactoryIsProducingLeaderModulesCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryIsProducingLeaderModulesCall) Do(f func() bool) *FactoryIsProducingLeaderModulesCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryIsProducingLeaderModulesCall) DoAndReturn(f func() bool) *FactoryIsProducingLeaderModulesCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Name mocks base method.
func (m *MockFactory) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockFactoryMockRecorder) Name() *FactoryNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFactory)(nil).Name))
	return &FactoryNameCall{Call: call}
}

// FactoryNameCall wrap *gomock.Call
type FactoryNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryNameCall) Return(arg0 string) *FactoryNameCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryNameCall) Do(f func() string) *FactoryNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryNameCall) DoAndReturn(f func() string) *FactoryNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) New(arg0 any) *FactoryNewCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockFactory)(nil).New), arg0)
	return &FactoryNewCall{Call: call}
}

// FactoryNewCall wrap *gomock.Call
type FactoryNewCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryNewCall) Return(arg0 modagent.Module, arg1 error) *FactoryNewCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryNewCall) Do(f func(*modagent.Config) (modagent.Module, error)) *FactoryNewCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryNewCall) DoAndReturn(f func(*modagent.Config) (modagent.Module, error)) *FactoryNewCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// StartStopPhase mocks base method.
func (m *MockFactory) StartStopPhase() modshared.ModuleStartStopPhase {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartStopPhase")
	ret0, _ := ret[0].(modshared.ModuleStartStopPhase)
	return ret0
}

// StartStopPhase indicates an expected call of StartStopPhase.
func (mr *MockFactoryMockRecorder) StartStopPhase() *FactoryStartStopPhaseCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartStopPhase", reflect.TypeOf((*MockFactory)(nil).StartStopPhase))
	return &FactoryStartStopPhaseCall{Call: call}
}

// FactoryStartStopPhaseCall wrap *gomock.Call
type FactoryStartStopPhaseCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryStartStopPhaseCall) Return(arg0 modshared.ModuleStartStopPhase) *FactoryStartStopPhaseCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryStartStopPhaseCall) Do(f func() modshared.ModuleStartStopPhase) *FactoryStartStopPhaseCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryStartStopPhaseCall) DoAndReturn(f func() modshared.ModuleStartStopPhase) *FactoryStartStopPhaseCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockModuleMockRecorder) DefaultAndValidateConfiguration(arg0 any) *ModuleDefaultAndValidateConfigurationCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultAndValidateConfiguration", reflect.TypeOf((*MockModule)(nil).DefaultAndValidateConfiguration), arg0)
	return &ModuleDefaultAndValidateConfigurationCall{Call: call}
}

// ModuleDefaultAndValidateConfigurationCall wrap *gomock.Call
type ModuleDefaultAndValidateConfigurationCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ModuleDefaultAndValidateConfigurationCall) Return(arg0 error) *ModuleDefaultAndValidateConfigurationCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ModuleDefaultAndValidateConfigurationCall) Do(f func(*agentcfg.AgentConfiguration) error) *ModuleDefaultAndValidateConfigurationCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ModuleDefaultAndValidateConfigurationCall) DoAndReturn(f func(*agentcfg.AgentConfiguration) error) *ModuleDefaultAndValidateConfigurationCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Name mocks base method.
func (m *MockModule) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockModuleMockRecorder) Name() *ModuleNameCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockModule)(nil).Name))
	return &ModuleNameCall{Call: call}
}

// ModuleNameCall wrap *gomock.Call
type ModuleNameCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ModuleNameCall) Return(arg0 string) *ModuleNameCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ModuleNameCall) Do(f func() string) *ModuleNameCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ModuleNameCall) DoAndReturn(f func() string) *ModuleNameCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Run mocks base method.
func (m *MockModule) Run(arg0 context.Context, arg1 <-chan *agentcfg.AgentConfiguration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockModuleMockRecorder) Run(arg0, arg1 any) *ModuleRunCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockModule)(nil).Run), arg0, arg1)
	return &ModuleRunCall{Call: call}
}

// ModuleRunCall wrap *gomock.Call
type ModuleRunCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *ModuleRunCall) Return(arg0 error) *ModuleRunCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *ModuleRunCall) Do(f func(context.Context, <-chan *agentcfg.AgentConfiguration) error) *ModuleRunCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *ModuleRunCall) DoAndReturn(f func(context.Context, <-chan *agentcfg.AgentConfiguration) error) *ModuleRunCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}