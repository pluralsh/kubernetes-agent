// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/kubectl/pkg/cmd/util (interfaces: Factory)

// Package mock_k8s is a generated GoMock package.
package mock_k8s

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	meta "k8s.io/apimachinery/pkg/api/meta"
	resource "k8s.io/cli-runtime/pkg/resource"
	discovery "k8s.io/client-go/discovery"
	dynamic "k8s.io/client-go/dynamic"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	openapi "k8s.io/kubectl/pkg/util/openapi"
	validation "k8s.io/kubectl/pkg/validation"
)

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

// ClientForMapping mocks base method.
func (m *MockFactory) ClientForMapping(arg0 *meta.RESTMapping) (resource.RESTClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClientForMapping", arg0)
	ret0, _ := ret[0].(resource.RESTClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClientForMapping indicates an expected call of ClientForMapping.
func (mr *MockFactoryMockRecorder) ClientForMapping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientForMapping", reflect.TypeOf((*MockFactory)(nil).ClientForMapping), arg0)
}

// DynamicClient mocks base method.
func (m *MockFactory) DynamicClient() (dynamic.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DynamicClient")
	ret0, _ := ret[0].(dynamic.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DynamicClient indicates an expected call of DynamicClient.
func (mr *MockFactoryMockRecorder) DynamicClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DynamicClient", reflect.TypeOf((*MockFactory)(nil).DynamicClient))
}

// KubernetesClientSet mocks base method.
func (m *MockFactory) KubernetesClientSet() (*kubernetes.Clientset, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KubernetesClientSet")
	ret0, _ := ret[0].(*kubernetes.Clientset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KubernetesClientSet indicates an expected call of KubernetesClientSet.
func (mr *MockFactoryMockRecorder) KubernetesClientSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KubernetesClientSet", reflect.TypeOf((*MockFactory)(nil).KubernetesClientSet))
}

// NewBuilder mocks base method.
func (m *MockFactory) NewBuilder() *resource.Builder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBuilder")
	ret0, _ := ret[0].(*resource.Builder)
	return ret0
}

// NewBuilder indicates an expected call of NewBuilder.
func (mr *MockFactoryMockRecorder) NewBuilder() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBuilder", reflect.TypeOf((*MockFactory)(nil).NewBuilder))
}

// OpenAPIGetter mocks base method.
func (m *MockFactory) OpenAPIGetter() discovery.OpenAPISchemaInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenAPIGetter")
	ret0, _ := ret[0].(discovery.OpenAPISchemaInterface)
	return ret0
}

// OpenAPIGetter indicates an expected call of OpenAPIGetter.
func (mr *MockFactoryMockRecorder) OpenAPIGetter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAPIGetter", reflect.TypeOf((*MockFactory)(nil).OpenAPIGetter))
}

// OpenAPISchema mocks base method.
func (m *MockFactory) OpenAPISchema() (openapi.Resources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenAPISchema")
	ret0, _ := ret[0].(openapi.Resources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenAPISchema indicates an expected call of OpenAPISchema.
func (mr *MockFactoryMockRecorder) OpenAPISchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAPISchema", reflect.TypeOf((*MockFactory)(nil).OpenAPISchema))
}

// RESTClient mocks base method.
func (m *MockFactory) RESTClient() (*rest.RESTClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RESTClient")
	ret0, _ := ret[0].(*rest.RESTClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RESTClient indicates an expected call of RESTClient.
func (mr *MockFactoryMockRecorder) RESTClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RESTClient", reflect.TypeOf((*MockFactory)(nil).RESTClient))
}

// ToDiscoveryClient mocks base method.
func (m *MockFactory) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToDiscoveryClient")
	ret0, _ := ret[0].(discovery.CachedDiscoveryInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToDiscoveryClient indicates an expected call of ToDiscoveryClient.
func (mr *MockFactoryMockRecorder) ToDiscoveryClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToDiscoveryClient", reflect.TypeOf((*MockFactory)(nil).ToDiscoveryClient))
}

// ToRESTConfig mocks base method.
func (m *MockFactory) ToRESTConfig() (*rest.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToRESTConfig")
	ret0, _ := ret[0].(*rest.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToRESTConfig indicates an expected call of ToRESTConfig.
func (mr *MockFactoryMockRecorder) ToRESTConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRESTConfig", reflect.TypeOf((*MockFactory)(nil).ToRESTConfig))
}

// ToRESTMapper mocks base method.
func (m *MockFactory) ToRESTMapper() (meta.RESTMapper, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToRESTMapper")
	ret0, _ := ret[0].(meta.RESTMapper)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToRESTMapper indicates an expected call of ToRESTMapper.
func (mr *MockFactoryMockRecorder) ToRESTMapper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRESTMapper", reflect.TypeOf((*MockFactory)(nil).ToRESTMapper))
}

// ToRawKubeConfigLoader mocks base method.
func (m *MockFactory) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToRawKubeConfigLoader")
	ret0, _ := ret[0].(clientcmd.ClientConfig)
	return ret0
}

// ToRawKubeConfigLoader indicates an expected call of ToRawKubeConfigLoader.
func (mr *MockFactoryMockRecorder) ToRawKubeConfigLoader() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRawKubeConfigLoader", reflect.TypeOf((*MockFactory)(nil).ToRawKubeConfigLoader))
}

// UnstructuredClientForMapping mocks base method.
func (m *MockFactory) UnstructuredClientForMapping(arg0 *meta.RESTMapping) (resource.RESTClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnstructuredClientForMapping", arg0)
	ret0, _ := ret[0].(resource.RESTClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnstructuredClientForMapping indicates an expected call of UnstructuredClientForMapping.
func (mr *MockFactoryMockRecorder) UnstructuredClientForMapping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnstructuredClientForMapping", reflect.TypeOf((*MockFactory)(nil).UnstructuredClientForMapping), arg0)
}

// Validator mocks base method.
func (m *MockFactory) Validator(arg0 string, arg1 *resource.QueryParamVerifier) (validation.Schema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validator", arg0, arg1)
	ret0, _ := ret[0].(validation.Schema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validator indicates an expected call of Validator.
func (mr *MockFactoryMockRecorder) Validator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validator", reflect.TypeOf((*MockFactory)(nil).Validator), arg0, arg1)
}