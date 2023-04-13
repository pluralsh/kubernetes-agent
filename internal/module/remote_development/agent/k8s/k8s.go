package k8s

import "context"

type Client interface {
	NamespaceExists(ctx context.Context, name string) bool
	CreateNamespace(ctx context.Context, name string) error
	DeleteNamespace(ctx context.Context, name string) error
	Apply(ctx context.Context, namespace, config string) error
}
