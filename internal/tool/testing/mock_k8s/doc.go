package mock_k8s

//go:generate go run github.com/golang/mock/mockgen -destination "resource.go" -package "mock_k8s" "k8s.io/cli-runtime/pkg/resource" "RESTClientGetter"
//go:generate go run github.com/golang/mock/mockgen -destination "meta.go" -package "mock_k8s" "k8s.io/apimachinery/pkg/api/meta" "ResettableRESTMapper"
//go:generate go run github.com/golang/mock/mockgen -destination "cache.go" -package "mock_k8s" "k8s.io/client-go/tools/cache" "Indexer,GenericLister,GenericNamespaceLister"
//go:generate go run github.com/golang/mock/mockgen -destination "dynamic.go" -package "mock_k8s" "k8s.io/client-go/dynamic" "NamespaceableResourceInterface,ResourceInterface"
//go:generate go run github.com/golang/mock/mockgen -destination "core_v1.go" -package "mock_k8s" "k8s.io/client-go/kubernetes/typed/core/v1" "CoreV1Interface,SecretInterface"
//go:generate go run github.com/golang/mock/mockgen -destination "workqueue.go" -package "mock_k8s" -mock_names "RateLimitingInterface=MockRateLimitingWorkqueue" "k8s.io/client-go/util/workqueue" "RateLimitingInterface"
//go:generate go run github.com/golang/mock/mockgen -destination "apiextensionclient_v1.go" -package "mock_k8s" "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1" "ApiextensionsV1Interface,CustomResourceDefinitionInterface"
