package k8s

import (
	"context"
)

type MockClient struct {
	NamespaceStore map[string]struct{}
	ApplyRecorder  string
	MockError      error
}

func NewMockClient() *MockClient {
	return &MockClient{
		NamespaceStore: map[string]struct{}{},
	}
}

func (m *MockClient) NamespaceExists(_ context.Context, name string) (bool, error) {
	_, ok := m.NamespaceStore[name]
	return ok, nil
}

func (m *MockClient) CreateNamespace(_ context.Context, name string) error {
	m.NamespaceStore[name] = struct{}{}
	return nil
}

func (m *MockClient) DeleteNamespace(_ context.Context, name string) error {
	delete(m.NamespaceStore, name)
	return nil
}

func (m *MockClient) Apply(_ context.Context, config string) <-chan error {
	m.ApplyRecorder = config

	if m.MockError != nil {
		errorCh := make(chan error)
		go func() {
			defer close(errorCh)
			errorCh <- m.MockError
			m.MockError = nil
		}()

		return errorCh
	}

	return nil
}
