// Package mock_stdlib contains Go standard library mocks
package mock_stdlib

import "net/http"

//go:generate mockgen.sh -destination "net_http_custom.go" -package "mock_stdlib" "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_stdlib" "ResponseWriterFlusher"

//go:generate mockgen.sh -destination "net.go" -package "mock_stdlib" "net" "Conn"

//go:generate mockgen.sh -destination "net_http.go" -package "mock_stdlib" "net/http" "RoundTripper"

type ResponseWriterFlusher interface {
	http.ResponseWriter
	http.Flusher
	http.Hijacker
}
