package mock_kubernetes_api

//go:generate mockgen.sh -destination "rpc.go" -package "mock_kubernetes_api" "github.com/pluralsh/kuberentes-agent/pkg/module/kubernetes_api/rpc" "KubernetesApiClient,KubernetesApi_MakeRequestClient"
