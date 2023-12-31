// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/kubectl/pkg/cmd/util (interfaces: Factory)
//
// Generated by this command:
//
//	mockgen -typed -destination kubectl_cmd_util.go -package mock_k8s k8s.io/kubectl/pkg/cmd/util Factory
//
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
	openapi "k8s.io/client-go/openapi"
	rest "k8s.io/client-go/rest"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	openapi0 "k8s.io/kubectl/pkg/util/openapi"
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
func (mr *MockFactoryMockRecorder) ClientForMapping(arg0 any) *FactoryClientForMappingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClientForMapping", reflect.TypeOf((*MockFactory)(nil).ClientForMapping), arg0)
	return &FactoryClientForMappingCall{Call: call}
}

// FactoryClientForMappingCall wrap *gomock.Call
type FactoryClientForMappingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryClientForMappingCall) Return(arg0 resource.RESTClient, arg1 error) *FactoryClientForMappingCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryClientForMappingCall) Do(f func(*meta.RESTMapping) (resource.RESTClient, error)) *FactoryClientForMappingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryClientForMappingCall) DoAndReturn(f func(*meta.RESTMapping) (resource.RESTClient, error)) *FactoryClientForMappingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) DynamicClient() *FactoryDynamicClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DynamicClient", reflect.TypeOf((*MockFactory)(nil).DynamicClient))
	return &FactoryDynamicClientCall{Call: call}
}

// FactoryDynamicClientCall wrap *gomock.Call
type FactoryDynamicClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryDynamicClientCall) Return(arg0 dynamic.Interface, arg1 error) *FactoryDynamicClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryDynamicClientCall) Do(f func() (dynamic.Interface, error)) *FactoryDynamicClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryDynamicClientCall) DoAndReturn(f func() (dynamic.Interface, error)) *FactoryDynamicClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) KubernetesClientSet() *FactoryKubernetesClientSetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KubernetesClientSet", reflect.TypeOf((*MockFactory)(nil).KubernetesClientSet))
	return &FactoryKubernetesClientSetCall{Call: call}
}

// FactoryKubernetesClientSetCall wrap *gomock.Call
type FactoryKubernetesClientSetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryKubernetesClientSetCall) Return(arg0 *kubernetes.Clientset, arg1 error) *FactoryKubernetesClientSetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryKubernetesClientSetCall) Do(f func() (*kubernetes.Clientset, error)) *FactoryKubernetesClientSetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryKubernetesClientSetCall) DoAndReturn(f func() (*kubernetes.Clientset, error)) *FactoryKubernetesClientSetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// NewBuilder mocks base method.
func (m *MockFactory) NewBuilder() *resource.Builder {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewBuilder")
	ret0, _ := ret[0].(*resource.Builder)
	return ret0
}

// NewBuilder indicates an expected call of NewBuilder.
func (mr *MockFactoryMockRecorder) NewBuilder() *FactoryNewBuilderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewBuilder", reflect.TypeOf((*MockFactory)(nil).NewBuilder))
	return &FactoryNewBuilderCall{Call: call}
}

// FactoryNewBuilderCall wrap *gomock.Call
type FactoryNewBuilderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryNewBuilderCall) Return(arg0 *resource.Builder) *FactoryNewBuilderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryNewBuilderCall) Do(f func() *resource.Builder) *FactoryNewBuilderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryNewBuilderCall) DoAndReturn(f func() *resource.Builder) *FactoryNewBuilderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OpenAPISchema mocks base method.
func (m *MockFactory) OpenAPISchema() (openapi0.Resources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenAPISchema")
	ret0, _ := ret[0].(openapi0.Resources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenAPISchema indicates an expected call of OpenAPISchema.
func (mr *MockFactoryMockRecorder) OpenAPISchema() *FactoryOpenAPISchemaCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAPISchema", reflect.TypeOf((*MockFactory)(nil).OpenAPISchema))
	return &FactoryOpenAPISchemaCall{Call: call}
}

// FactoryOpenAPISchemaCall wrap *gomock.Call
type FactoryOpenAPISchemaCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryOpenAPISchemaCall) Return(arg0 openapi0.Resources, arg1 error) *FactoryOpenAPISchemaCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryOpenAPISchemaCall) Do(f func() (openapi0.Resources, error)) *FactoryOpenAPISchemaCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryOpenAPISchemaCall) DoAndReturn(f func() (openapi0.Resources, error)) *FactoryOpenAPISchemaCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// OpenAPIV3Client mocks base method.
func (m *MockFactory) OpenAPIV3Client() (openapi.Client, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenAPIV3Client")
	ret0, _ := ret[0].(openapi.Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenAPIV3Client indicates an expected call of OpenAPIV3Client.
func (mr *MockFactoryMockRecorder) OpenAPIV3Client() *FactoryOpenAPIV3ClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenAPIV3Client", reflect.TypeOf((*MockFactory)(nil).OpenAPIV3Client))
	return &FactoryOpenAPIV3ClientCall{Call: call}
}

// FactoryOpenAPIV3ClientCall wrap *gomock.Call
type FactoryOpenAPIV3ClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryOpenAPIV3ClientCall) Return(arg0 openapi.Client, arg1 error) *FactoryOpenAPIV3ClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryOpenAPIV3ClientCall) Do(f func() (openapi.Client, error)) *FactoryOpenAPIV3ClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryOpenAPIV3ClientCall) DoAndReturn(f func() (openapi.Client, error)) *FactoryOpenAPIV3ClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) RESTClient() *FactoryRESTClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RESTClient", reflect.TypeOf((*MockFactory)(nil).RESTClient))
	return &FactoryRESTClientCall{Call: call}
}

// FactoryRESTClientCall wrap *gomock.Call
type FactoryRESTClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryRESTClientCall) Return(arg0 *rest.RESTClient, arg1 error) *FactoryRESTClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryRESTClientCall) Do(f func() (*rest.RESTClient, error)) *FactoryRESTClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryRESTClientCall) DoAndReturn(f func() (*rest.RESTClient, error)) *FactoryRESTClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) ToDiscoveryClient() *FactoryToDiscoveryClientCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToDiscoveryClient", reflect.TypeOf((*MockFactory)(nil).ToDiscoveryClient))
	return &FactoryToDiscoveryClientCall{Call: call}
}

// FactoryToDiscoveryClientCall wrap *gomock.Call
type FactoryToDiscoveryClientCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryToDiscoveryClientCall) Return(arg0 discovery.CachedDiscoveryInterface, arg1 error) *FactoryToDiscoveryClientCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryToDiscoveryClientCall) Do(f func() (discovery.CachedDiscoveryInterface, error)) *FactoryToDiscoveryClientCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryToDiscoveryClientCall) DoAndReturn(f func() (discovery.CachedDiscoveryInterface, error)) *FactoryToDiscoveryClientCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) ToRESTConfig() *FactoryToRESTConfigCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRESTConfig", reflect.TypeOf((*MockFactory)(nil).ToRESTConfig))
	return &FactoryToRESTConfigCall{Call: call}
}

// FactoryToRESTConfigCall wrap *gomock.Call
type FactoryToRESTConfigCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryToRESTConfigCall) Return(arg0 *rest.Config, arg1 error) *FactoryToRESTConfigCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryToRESTConfigCall) Do(f func() (*rest.Config, error)) *FactoryToRESTConfigCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryToRESTConfigCall) DoAndReturn(f func() (*rest.Config, error)) *FactoryToRESTConfigCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) ToRESTMapper() *FactoryToRESTMapperCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRESTMapper", reflect.TypeOf((*MockFactory)(nil).ToRESTMapper))
	return &FactoryToRESTMapperCall{Call: call}
}

// FactoryToRESTMapperCall wrap *gomock.Call
type FactoryToRESTMapperCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryToRESTMapperCall) Return(arg0 meta.RESTMapper, arg1 error) *FactoryToRESTMapperCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryToRESTMapperCall) Do(f func() (meta.RESTMapper, error)) *FactoryToRESTMapperCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryToRESTMapperCall) DoAndReturn(f func() (meta.RESTMapper, error)) *FactoryToRESTMapperCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ToRawKubeConfigLoader mocks base method.
func (m *MockFactory) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToRawKubeConfigLoader")
	ret0, _ := ret[0].(clientcmd.ClientConfig)
	return ret0
}

// ToRawKubeConfigLoader indicates an expected call of ToRawKubeConfigLoader.
func (mr *MockFactoryMockRecorder) ToRawKubeConfigLoader() *FactoryToRawKubeConfigLoaderCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToRawKubeConfigLoader", reflect.TypeOf((*MockFactory)(nil).ToRawKubeConfigLoader))
	return &FactoryToRawKubeConfigLoaderCall{Call: call}
}

// FactoryToRawKubeConfigLoaderCall wrap *gomock.Call
type FactoryToRawKubeConfigLoaderCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryToRawKubeConfigLoaderCall) Return(arg0 clientcmd.ClientConfig) *FactoryToRawKubeConfigLoaderCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryToRawKubeConfigLoaderCall) Do(f func() clientcmd.ClientConfig) *FactoryToRawKubeConfigLoaderCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryToRawKubeConfigLoaderCall) DoAndReturn(f func() clientcmd.ClientConfig) *FactoryToRawKubeConfigLoaderCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
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
func (mr *MockFactoryMockRecorder) UnstructuredClientForMapping(arg0 any) *FactoryUnstructuredClientForMappingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnstructuredClientForMapping", reflect.TypeOf((*MockFactory)(nil).UnstructuredClientForMapping), arg0)
	return &FactoryUnstructuredClientForMappingCall{Call: call}
}

// FactoryUnstructuredClientForMappingCall wrap *gomock.Call
type FactoryUnstructuredClientForMappingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryUnstructuredClientForMappingCall) Return(arg0 resource.RESTClient, arg1 error) *FactoryUnstructuredClientForMappingCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryUnstructuredClientForMappingCall) Do(f func(*meta.RESTMapping) (resource.RESTClient, error)) *FactoryUnstructuredClientForMappingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryUnstructuredClientForMappingCall) DoAndReturn(f func(*meta.RESTMapping) (resource.RESTClient, error)) *FactoryUnstructuredClientForMappingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Validator mocks base method.
func (m *MockFactory) Validator(arg0 string) (validation.Schema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validator", arg0)
	ret0, _ := ret[0].(validation.Schema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validator indicates an expected call of Validator.
func (mr *MockFactoryMockRecorder) Validator(arg0 any) *FactoryValidatorCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validator", reflect.TypeOf((*MockFactory)(nil).Validator), arg0)
	return &FactoryValidatorCall{Call: call}
}

// FactoryValidatorCall wrap *gomock.Call
type FactoryValidatorCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *FactoryValidatorCall) Return(arg0 validation.Schema, arg1 error) *FactoryValidatorCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *FactoryValidatorCall) Do(f func(string) (validation.Schema, error)) *FactoryValidatorCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *FactoryValidatorCall) DoAndReturn(f func(string) (validation.Schema, error)) *FactoryValidatorCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
