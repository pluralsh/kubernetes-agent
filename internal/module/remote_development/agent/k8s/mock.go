package k8s

import (
	"context"
)

type MockClient struct {
	NamespaceStore map[string]struct{}
	ApplyRecorder  string
}

func NewMockClient() *MockClient {
	return &MockClient{
		NamespaceStore: map[string]struct{}{},
	}
}

func (m *MockClient) NamespaceExists(ctx context.Context, name string) bool {
	_, ok := m.NamespaceStore[name]
	return ok
}

func (m *MockClient) CreateNamespace(ctx context.Context, name string) error {
	m.NamespaceStore[name] = struct{}{}
	return nil
}

func (m *MockClient) DeleteNamespace(ctx context.Context, name string) error {
	delete(m.NamespaceStore, name)
	return nil
}

func (m *MockClient) Apply(ctx context.Context, namespace, config string) error {
	m.ApplyRecorder = config
	return nil
}
