// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1 (interfaces: ApiextensionsV1Interface,CustomResourceDefinitionInterface)
//
// Generated by this command:
//
//	mockgen -destination apiextensionclient_v1.go -package mock_k8s k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1 ApiextensionsV1Interface,CustomResourceDefinitionInterface
//
// Package mock_k8s is a generated GoMock package.
package mock_k8s

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v10 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v11 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MockApiextensionsV1Interface is a mock of ApiextensionsV1Interface interface.
type MockApiextensionsV1Interface struct {
	ctrl     *gomock.Controller
	recorder *MockApiextensionsV1InterfaceMockRecorder
}

// MockApiextensionsV1InterfaceMockRecorder is the mock recorder for MockApiextensionsV1Interface.
type MockApiextensionsV1InterfaceMockRecorder struct {
	mock *MockApiextensionsV1Interface
}

// NewMockApiextensionsV1Interface creates a new mock instance.
func NewMockApiextensionsV1Interface(ctrl *gomock.Controller) *MockApiextensionsV1Interface {
	mock := &MockApiextensionsV1Interface{ctrl: ctrl}
	mock.recorder = &MockApiextensionsV1InterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApiextensionsV1Interface) EXPECT() *MockApiextensionsV1InterfaceMockRecorder {
	return m.recorder
}

// CustomResourceDefinitions mocks base method.
func (m *MockApiextensionsV1Interface) CustomResourceDefinitions() v10.CustomResourceDefinitionInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CustomResourceDefinitions")
	ret0, _ := ret[0].(v10.CustomResourceDefinitionInterface)
	return ret0
}

// CustomResourceDefinitions indicates an expected call of CustomResourceDefinitions.
func (mr *MockApiextensionsV1InterfaceMockRecorder) CustomResourceDefinitions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CustomResourceDefinitions", reflect.TypeOf((*MockApiextensionsV1Interface)(nil).CustomResourceDefinitions))
}

// RESTClient mocks base method.
func (m *MockApiextensionsV1Interface) RESTClient() rest.Interface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RESTClient")
	ret0, _ := ret[0].(rest.Interface)
	return ret0
}

// RESTClient indicates an expected call of RESTClient.
func (mr *MockApiextensionsV1InterfaceMockRecorder) RESTClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RESTClient", reflect.TypeOf((*MockApiextensionsV1Interface)(nil).RESTClient))
}

// MockCustomResourceDefinitionInterface is a mock of CustomResourceDefinitionInterface interface.
type MockCustomResourceDefinitionInterface struct {
	ctrl     *gomock.Controller
	recorder *MockCustomResourceDefinitionInterfaceMockRecorder
}

// MockCustomResourceDefinitionInterfaceMockRecorder is the mock recorder for MockCustomResourceDefinitionInterface.
type MockCustomResourceDefinitionInterfaceMockRecorder struct {
	mock *MockCustomResourceDefinitionInterface
}

// NewMockCustomResourceDefinitionInterface creates a new mock instance.
func NewMockCustomResourceDefinitionInterface(ctrl *gomock.Controller) *MockCustomResourceDefinitionInterface {
	mock := &MockCustomResourceDefinitionInterface{ctrl: ctrl}
	mock.recorder = &MockCustomResourceDefinitionInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomResourceDefinitionInterface) EXPECT() *MockCustomResourceDefinitionInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCustomResourceDefinitionInterface) Create(arg0 context.Context, arg1 *v1.CustomResourceDefinition, arg2 v11.CreateOptions) (*v1.CustomResourceDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.CustomResourceDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) Create(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).Create), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockCustomResourceDefinitionInterface) Delete(arg0 context.Context, arg1 string, arg2 v11.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) Delete(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).Delete), arg0, arg1, arg2)
}

// DeleteCollection mocks base method.
func (m *MockCustomResourceDefinitionInterface) DeleteCollection(arg0 context.Context, arg1 v11.DeleteOptions, arg2 v11.ListOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCollection", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCollection indicates an expected call of DeleteCollection.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) DeleteCollection(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCollection", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).DeleteCollection), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *MockCustomResourceDefinitionInterface) Get(arg0 context.Context, arg1 string, arg2 v11.GetOptions) (*v1.CustomResourceDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.CustomResourceDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) Get(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).Get), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockCustomResourceDefinitionInterface) List(arg0 context.Context, arg1 v11.ListOptions) (*v1.CustomResourceDefinitionList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*v1.CustomResourceDefinitionList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).List), arg0, arg1)
}

// Patch mocks base method.
func (m *MockCustomResourceDefinitionInterface) Patch(arg0 context.Context, arg1 string, arg2 types.PatchType, arg3 []byte, arg4 v11.PatchOptions, arg5 ...string) (*v1.CustomResourceDefinition, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1, arg2, arg3, arg4}
	for _, a := range arg5 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(*v1.CustomResourceDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) Patch(arg0, arg1, arg2, arg3, arg4 any, arg5 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1, arg2, arg3, arg4}, arg5...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).Patch), varargs...)
}

// Update mocks base method.
func (m *MockCustomResourceDefinitionInterface) Update(arg0 context.Context, arg1 *v1.CustomResourceDefinition, arg2 v11.UpdateOptions) (*v1.CustomResourceDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.CustomResourceDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) Update(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).Update), arg0, arg1, arg2)
}

// UpdateStatus mocks base method.
func (m *MockCustomResourceDefinitionInterface) UpdateStatus(arg0 context.Context, arg1 *v1.CustomResourceDefinition, arg2 v11.UpdateOptions) (*v1.CustomResourceDefinition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1.CustomResourceDefinition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) UpdateStatus(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).UpdateStatus), arg0, arg1, arg2)
}

// Watch mocks base method.
func (m *MockCustomResourceDefinitionInterface) Watch(arg0 context.Context, arg1 v11.ListOptions) (watch.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0, arg1)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch.
func (mr *MockCustomResourceDefinitionInterfaceMockRecorder) Watch(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockCustomResourceDefinitionInterface)(nil).Watch), arg0, arg1)
}
