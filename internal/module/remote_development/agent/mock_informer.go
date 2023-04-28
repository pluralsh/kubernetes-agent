package agent

import (
	"context"
)

type mockInformer struct {
	Resources map[string]*parsedWorkspace
}

func newMockInformer() *mockInformer {
	return &mockInformer{
		Resources: make(map[string]*parsedWorkspace),
	}
}

func (i *mockInformer) Start(ctx context.Context) error {
	return nil
}

func (i *mockInformer) List() []*parsedWorkspace {
	result := make([]*parsedWorkspace, 0, len(i.Resources))
	for _, pw := range i.Resources {
		result = append(result, pw)
	}
	return result
}

func (i *mockInformer) Stop() {
	// do nothing
}
