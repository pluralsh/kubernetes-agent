package grpctool

import (
	"context"

	"github.com/pluralsh/kuberentes-agent/internal/api"
	"google.golang.org/grpc/credentials"
)

const (
	MetadataAuthorization = "authorization"
)

func NewTokenCredentials(token api.AgentToken, insecure bool) credentials.PerRPCCredentials {
	return &tokenCredentials{
		metadata: map[string]string{
			MetadataAuthorization: "Bearer " + string(token),
		},
		insecure: insecure,
	}
}

type tokenCredentials struct {
	metadata map[string]string
	insecure bool
}

func (t *tokenCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return t.metadata, nil
}

func (t *tokenCredentials) RequireTransportSecurity() bool {
	return !t.insecure
}
