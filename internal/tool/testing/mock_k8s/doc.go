package mock_k8s

//go:generate go run github.com/golang/mock/mockgen -destination "resource.go" -package "mock_k8s" "k8s.io/cli-runtime/pkg/resource" "RESTClientGetter"
//go:generate go run github.com/golang/mock/mockgen -destination "meta.go" -package "mock_k8s" "k8s.io/apimachinery/pkg/api/meta" "ResettableRESTMapper"
