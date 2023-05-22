package client

import (
	"context"
	"math/rand"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/backoff"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/dnsresolver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/internal_client"
	"google.golang.org/grpc"
)

// DialContext dials the Gitaly at the given address with the provided options. Valid address formats are
// 'unix:<socket path>' for Unix sockets, 'tcp://<host:port>' for insecure TCP connections and 'tls://<host:port>'
// for TCP+TLS connections.
//
// The returned ClientConns are configured with tracing and correlation id interceptors to ensure they are propagated
// correctly. They're also configured to send Keepalives with settings matching what Gitaly expects.
//
// connOpts should not contain `grpc.WithInsecure` as DialContext determines whether it is needed or not from the
// scheme. `grpc.TransportCredentials` should not be provided either as those are handled internally as well.
func DialContext(ctx context.Context, rawAddress string, connOpts []grpc.DialOption) (*grpc.ClientConn, error) {
	return client.Dial(ctx, rawAddress, connOpts, nil)
}

// DNSResolverBuilderConfig exposes the DNS resolver builder option. It is used to build Gitaly
// custom DNS resolver.
type DNSResolverBuilderConfig dnsresolver.BuilderConfig

// DefaultDNSResolverBuilderConfig returns the default options for building DNS resolver.
func DefaultDNSResolverBuilderConfig() *DNSResolverBuilderConfig {
	return &DNSResolverBuilderConfig{
		RefreshRate:     5 * time.Minute,
		Backoff:         backoff.NewDefaultExponential(rand.New(rand.NewSource(time.Now().UnixNano()))), // nolint: gosec
		DefaultGrpcPort: "443",
	}
}

// WithGitalyDNSResolver defines a gRPC dial option for injecting Gitaly's custom DNS resolver. This
// resolver watches for the changes of target URL periodically and update the target subchannels
// accordingly.
func WithGitalyDNSResolver(opts *DNSResolverBuilderConfig) grpc.DialOption {
	return grpc.WithResolvers(dnsresolver.NewBuilder((*dnsresolver.BuilderConfig)(opts)))
}
