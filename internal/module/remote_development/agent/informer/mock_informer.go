package informer

import (
	"context"
)

type MockInformer struct {
	Resources map[string]*ParsedWorkspace
}

func NewMockInformer() *MockInformer {
	return &MockInformer{
		Resources: make(map[string]*ParsedWorkspace),
	}
}

func (i *MockInformer) Start(ctx context.Context) error {
	return nil
}

func (i *MockInformer) List() []*ParsedWorkspace {
	result := make([]*ParsedWorkspace, 0, len(i.Resources))
	for _, pw := range i.Resources {
		result = append(result, pw)
	}
	return result
}
