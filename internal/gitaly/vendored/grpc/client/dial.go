package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/dnsresolver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/gitalypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/encoding/protojson"
)

type connectionType int

const (
	invalidConnection connectionType = iota
	tcpConnection
	tlsConnection
	unixConnection
	dnsConnection
)

func getConnectionType(rawAddress string) connectionType {
	u, err := url.Parse(rawAddress)
	if err != nil {
		return invalidConnection
	}

	switch u.Scheme {
	case "tls":
		return tlsConnection
	case "unix":
		return unixConnection
	case "tcp":
		return tcpConnection
	case "dns":
		return dnsConnection
	default:
		return invalidConnection
	}
}

// Handshaker is an interface that allows for wrapping the transport credentials
// with a custom handshake.
type Handshaker interface {
	// ClientHandshake wraps the provided credentials and returns new credentials.
	ClientHandshake(credentials.TransportCredentials) credentials.TransportCredentials
}

type dialConfig struct {
	handshaker Handshaker
	grpcOpts   []grpc.DialOption
	creds      credentials.TransportCredentials
}

// DialOption is an option that can be passed to Dial.
type DialOption func(*dialConfig)

// WithHandshaker sets up the given handshaker so that it's passed as the transport credentials
// which would be otherwise set. The transport credentials returned by handshaker are then set
// instead.
func WithHandshaker(handshaker Handshaker) DialOption {
	return func(cfg *dialConfig) {
		cfg.handshaker = handshaker
	}
}

// WithGrpcOptions will set up the given gRPC dial options so that they will be used when calling
// `grpc.DialContext()`.
func WithGrpcOptions(opts []grpc.DialOption) DialOption {
	return func(cfg *dialConfig) {
		cfg.grpcOpts = opts
	}
}

// WithTransportCredentials sets up the given credentials. By default, non-TLS connections will use
// insecure credentials whereas TLS connections will use the x509 system certificate pool. This
// option allows callers to override these defaults.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return func(cfg *dialConfig) {
		cfg.creds = creds
	}
}

// Dial dials a Gitaly node serving at the given address. Dial is used by the public 'client' package
// and the expected behavior is mostly documented there.
func Dial(ctx context.Context, rawAddress string, opts ...DialOption) (*grpc.ClientConn, error) {
	var dialCfg dialConfig
	for _, opt := range opts {
		opt(&dialCfg)
	}

	// copy to avoid potentially mutating the backing array of the passed slice
	connOpts := cloneOpts(dialCfg.grpcOpts)

	var canonicalAddress string
	var err error
	var secure bool

	switch getConnectionType(rawAddress) {
	case invalidConnection:
		return nil, fmt.Errorf("invalid connection string: %q", rawAddress)
	case tlsConnection:
		canonicalAddress, err = extractHostFromRemoteURL(rawAddress) // Ensure the form: "host:port" ...
		if err != nil {
			return nil, fmt.Errorf("failed to extract host for 'tls' connection: %w", err)
		}

		secure = true
	case tcpConnection:
		canonicalAddress, err = extractHostFromRemoteURL(rawAddress) // Ensure the form: "host:port" ...
		if err != nil {
			return nil, fmt.Errorf("failed to extract host for 'tcp' connection: %w", err)
		}
	case dnsConnection:
		err = dnsresolver.ValidateURL(rawAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to parse target for 'dns' connection: %w", err)
		}
		canonicalAddress = rawAddress // DNS Resolver will handle this

	case unixConnection:
		canonicalAddress = rawAddress // This will be overridden by the custom dialer...
		connOpts = append(
			connOpts,
			// Use a custom dialer to ensure that we don't experience
			// issues in environments that have proxy configurations
			// https://gitlab.com/gitlab-org/gitaly/merge_requests/1072#note_140408512
			grpc.WithContextDialer(func(ctx context.Context, addr string) (conn net.Conn, err error) {
				path, err := extractPathFromSocketURL(addr)
				if err != nil {
					return nil, fmt.Errorf("failed to extract host for 'unix' connection: %w", err)
				}

				d := net.Dialer{}
				return d.DialContext(ctx, "unix", path)
			}),
		)
	}

	transportCredentials := dialCfg.creds
	if transportCredentials == nil {
		if !secure {
			transportCredentials = insecure.NewCredentials()
		} else {
			certPool, err := x509.SystemCertPool() // nolint: govet
			if err != nil {
				return nil, fmt.Errorf("failed to get system certificat pool for 'tls' connection: %w", err)
			}

			transportCredentials = credentials.NewTLS(&tls.Config{
				RootCAs:    certPool,
				MinVersion: tls.VersionTLS12,
			})
		}
	}

	if dialCfg.handshaker != nil {
		transportCredentials = dialCfg.handshaker.ClientHandshake(transportCredentials)
	}

	connOpts = append(connOpts,
		grpc.WithTransportCredentials(transportCredentials),
		// grpc.KeepaliveParams must be specified at least as large as what is allowed by the
		// server-side grpc.KeepaliveEnforcementPolicy
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                20 * time.Second,
			PermitWithoutStream: true,
		}),
		// grpc.WithDisableServiceConfig ignores the service config provided by resolvers
		// when they resolve the target. gRPC provides this feature to inject service
		// config from external sources (DNS TXT record, for example). Gitaly doesn't need
		// this feature. When we implement a custom client-side load balancer, this feature
		// can even break the balancer. So, we should better disable it.
		// For more information, please visit
		// - https://github.com/grpc/proposal/blob/master/A2-service-configs-in-dns.md
		grpc.WithDisableServiceConfig(),
		// grpc.WithDefaultServiceConfig sets the recommended client-side load balancing
		// configuration to client dial. By default, gRPC clients don't support client-side load
		// balancing. After the connection to a host is established for the first time, that
		// client always sticks to that host. In all Gitaly clients, the connection is cached
		// somehow, usually one connection per host. It means they always stick to the same
		// host until the process restarts. This is not a problem in pure Gitaly environment.
		// In a cluster with more than one Praefect node, this behavior may cause serious
		// workload skew, especially after a fail-over event.
		//
		// This option configures the load balancing strategy to `round_robin`. This is a
		// built-in strategy grpc-go provides. When combining with service discovery via DNS,
		// a client can distribute its requests to all discovered nodes in a round-robin
		// fashion. The client can detect the connectivity changes, such as a node goes
		// down/up again. It evicts subsequent requests accordingly.
		//
		// For more information:
		// https://gitlab.com/groups/gitlab-org/-/epics/8971#note_1207008162
		grpc.WithDefaultServiceConfig(defaultServiceConfig()),
	)

	conn, err := grpc.DialContext(ctx, canonicalAddress, connOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %q connection: %w", canonicalAddress, err)
	}

	return conn, nil
}

func cloneOpts(opts []grpc.DialOption) []grpc.DialOption {
	clone := make([]grpc.DialOption, len(opts))
	copy(clone, opts)
	return clone
}

func defaultServiceConfig() string {
	serviceConfig := &gitalypb.ServiceConfig{
		LoadBalancingConfig: []*gitalypb.LoadBalancingConfig{{
			Policy: &gitalypb.LoadBalancingConfig_RoundRobin{},
		}},
	}
	configJSON, err := protojson.Marshal(serviceConfig)
	if err != nil {
		panic("fail to convert service config from protobuf to json")
	}

	return string(configJSON)
}
