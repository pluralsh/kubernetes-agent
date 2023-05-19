package client

import (
	"context"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/backoff"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/dnsresolver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/internal_client"
	"gitlab.com/gitlab-org/gitaly/v16/internal/grpc/sidechannel"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// DefaultDialOpts hold the default DialOptions for connection to Gitaly over UNIX-socket
var DefaultDialOpts = []grpc.DialOption{}

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

// Dial calls DialContext with the provided arguments and context.Background. Refer to DialContext's documentation
// for details.
func Dial(rawAddress string, connOpts []grpc.DialOption) (*grpc.ClientConn, error) {
	return DialContext(context.Background(), rawAddress, connOpts)
}

// DialSidechannel configures the dialer to establish a Gitaly
// backchannel connection instead of a regular gRPC connection. It also
// injects sr as a sidechannel registry, so that Gitaly can establish
// sidechannels back to the client.
func DialSidechannel(ctx context.Context, rawAddress string, sr *SidechannelRegistry, connOpts []grpc.DialOption) (*grpc.ClientConn, error) {
	clientHandshaker := sidechannel.NewClientHandshaker(sr.logger, sr.registry)
	return client.Dial(ctx, rawAddress, connOpts, clientHandshaker)
}

// FailOnNonTempDialError helps to identify if remote listener is ready to accept new connections.
func FailOnNonTempDialError() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithBlock(),
		grpc.FailOnNonTempDialError(true),
	}
}

// HealthCheckDialer uses provided dialer as an actual dialer, but issues a health check request to the remote
// to verify the connection was set properly and could be used with no issues.
func HealthCheckDialer(base Dialer) Dialer {
	return func(ctx context.Context, address string, dialOptions []grpc.DialOption) (*grpc.ClientConn, error) {
		cc, err := base(ctx, address, dialOptions)
		if err != nil {
			return nil, err
		}

		if _, err := healthpb.NewHealthClient(cc).Check(ctx, &healthpb.HealthCheckRequest{}); err != nil {
			_ = cc.Close()
			return nil, err
		}

		return cc, nil
	}
}

// DNSResolverBuilderConfig exposes the DNS resolver builder option. It is used to build Gitaly
// custom DNS resolver.
type DNSResolverBuilderConfig dnsresolver.BuilderConfig

// DefaultDNSResolverBuilderConfig returns the default options for building DNS resolver.
func DefaultDNSResolverBuilderConfig() *DNSResolverBuilderConfig {
	return &DNSResolverBuilderConfig{
		RefreshRate:     5 * time.Minute,
		Logger:          logrus.StandardLogger(),
		Backoff:         backoff.NewDefaultExponential(rand.New(rand.NewSource(time.Now().UnixNano()))),
		DefaultGrpcPort: "443",
	}
}

// WithGitalyDNSResolver defines a gRPC dial option for injecting Gitaly's custom DNS resolver. This
// resolver watches for the changes of target URL periodically and update the target subchannels
// accordingly.
func WithGitalyDNSResolver(opts *DNSResolverBuilderConfig) grpc.DialOption {
	return grpc.WithResolvers(dnsresolver.NewBuilder((*dnsresolver.BuilderConfig)(opts)))
}
