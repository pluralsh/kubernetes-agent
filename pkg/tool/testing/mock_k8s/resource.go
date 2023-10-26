// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/cli-runtime/pkg/resource (interfaces: RESTClientGetter)
//
// Generated by this command:
//
//	mockgen -typed -destination resource.go -package mock_k8s k8s.io/cli-runtime/pkg/resource RESTClientGetter
//
// Package mock_k8s is a generated GoMock package.
package mock_k8s

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	meta "k8s.io/apimachinery/pkg/api/meta"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
)

// MockRESTClientGetter is a mock of RESTClientGetter interface.
type MockRESTClientGetter struct {
	ctrl     *gomock.Controller
	recorder *MockRESTClientGetterMockRecorder
}

// MockRESTClientGetterMockRecorder is the mock recorder for MockRESTClientGetter.
type MockRESTClientGetterMockRecorder struct {
	mock *MockRESTClientGetter
}

// NewMockRESTClientGetter creates a new mock instance.
func NewMockRESTClientGetter(ctrl *gomock.Controller) *MockRESTClientGetter {
	mock := &MockRESTClientGetter{ctrl: ctrl}
	mock.recorder = &MockRESTClientGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRESTClientGetter) EXPECT() *MockRESTClientGetterMockRecorder {
	return m.recorder
}

// ToDiscoveryClient mocks base method.
func (m *MockRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToDiscoveryClient")
	ret0, _ := ret[0].(discovery.CachedDiscoveryInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToDiscoveryClient indicates an expected call of ToDiscoveryClient.
func (mr *MockRESTClientGetterMockRecorder) ToDiscoveryClient() *RESTClientGetterToDiscoveryClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToDiscoveryClient", reflect.TypeOf((*MockRESTClientGetter)(nil).ToDiscoveryClient))
	return &RESTClientGetterToDiscoveryClientCall{Call: call}
}

// RESTClientGetterToDiscoveryClientCall wrap *gomock.Call
type RESTClientGetterToDiscoveryClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RESTClientGetterToDiscoveryClientCall) Return(arg0 discovery.CachedDiscoveryInterface, arg1 error) *RESTClientGetterToDiscoveryClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RESTClientGetterToDiscoveryClientCall) Do(f func() (discovery.CachedDiscoveryInterface, error)) *RESTClientGetterToDiscoveryClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RESTClientGetterToDiscoveryClientCall) DoAndReturn(f func() (discovery.CachedDiscoveryInterface, error)) *RESTClientGetterToDiscoveryClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ToRESTConfig mocks base method.
func (m *MockRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToRESTConfig")
	ret0, _ := ret[0].(*rest.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToRESTConfig indicates an expected call of ToRESTConfig.
func (mr *MockRESTClientGetterMockRecorder) ToRESTConfig() *RESTClientGetterToRESTConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRESTConfig", reflect.TypeOf((*MockRESTClientGetter)(nil).ToRESTConfig))
	return &RESTClientGetterToRESTConfigCall{Call: call}
}

// RESTClientGetterToRESTConfigCall wrap *gomock.Call
type RESTClientGetterToRESTConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RESTClientGetterToRESTConfigCall) Return(arg0 *rest.Config, arg1 error) *RESTClientGetterToRESTConfigCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RESTClientGetterToRESTConfigCall) Do(f func() (*rest.Config, error)) *RESTClientGetterToRESTConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RESTClientGetterToRESTConfigCall) DoAndReturn(f func() (*rest.Config, error)) *RESTClientGetterToRESTConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ToRESTMapper mocks base method.
func (m *MockRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToRESTMapper")
	ret0, _ := ret[0].(meta.RESTMapper)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToRESTMapper indicates an expected call of ToRESTMapper.
func (mr *MockRESTClientGetterMockRecorder) ToRESTMapper() *RESTClientGetterToRESTMapperCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRESTMapper", reflect.TypeOf((*MockRESTClientGetter)(nil).ToRESTMapper))
	return &RESTClientGetterToRESTMapperCall{Call: call}
}

// RESTClientGetterToRESTMapperCall wrap *gomock.Call
type RESTClientGetterToRESTMapperCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *RESTClientGetterToRESTMapperCall) Return(arg0 meta.RESTMapper, arg1 error) *RESTClientGetterToRESTMapperCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *RESTClientGetterToRESTMapperCall) Do(f func() (meta.RESTMapper, error)) *RESTClientGetterToRESTMapperCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *RESTClientGetterToRESTMapperCall) DoAndReturn(f func() (meta.RESTMapper, error)) *RESTClientGetterToRESTMapperCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}