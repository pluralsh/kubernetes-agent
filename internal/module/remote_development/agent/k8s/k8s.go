package k8s

import "context"

type Client interface {
	NamespaceExists(ctx context.Context, name string) (bool, error)
	CreateNamespace(ctx context.Context, name string) error
	DeleteNamespace(ctx context.Context, name string) error
	Apply(ctx context.Context, config string) <-chan error
}
