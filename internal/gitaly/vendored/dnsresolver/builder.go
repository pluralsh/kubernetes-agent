package dnsresolver

import (
	"context"
	"strings"
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/backoff"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/structerr"
	"google.golang.org/grpc/resolver"
)

// Default DNS desc server port. This is a de-facto convention for both UDP and TCP.
const defaultDNSNameserverPort = "53"

// gRPC depends on the target's scheme to determine which resolver to use. Built-in DNS Resolver
// registers itself with "dns" scheme. We should use a different scheme for this resolver. However,
// Ruby, and other cares-based clients, don't support custom resolver. At GitLab, the gRPC target
// configuration is shared between components. To ensure the compatibility between clients, this
// resolver intentionally replaces the built-in resolver by itself.
// The client should use grpc.WithResolvers to inject Gitaly custom DNS resolver when resolving
// the target URL.
const dnsResolverScheme = "dns"

// BuilderConfig defines the configuration for customizing the builder.
type BuilderConfig struct {
	// RefreshRate determines the periodic refresh rate of the resolver. The resolver may issue
	// the resolver earlier if client connection demands
	RefreshRate time.Duration
	// Backoff defines the backoff strategy when the resolver fails to resolve or pushes new
	// state to client connection
	Backoff backoff.Strategy
	// DefaultGrpcPort sets the gRPC port if the target URL doesn't specify a target port
	DefaultGrpcPort string
	// authorityFinder is to inject a custom authority finder from the authority address in
	// the target URL. For example: dns://authority-host:authority-port/host:port
	authorityFinder func(authority string) (dnsLookuper, error)
}

// Builder is an object to build the resolver for a connection. A client connection uses the builder
// specified by grpc.WithResolvers dial option or the one fetched from global Resolver registry. The
// local option has higher precedence than the global one.
type Builder struct {
	opts *BuilderConfig
}

// NewBuilder creates a builder option with an input option
func NewBuilder(opts *BuilderConfig) *Builder {
	return &Builder{opts: opts}
}

// Scheme returns the scheme handled by this builder. Client connection queries the resolver based
// on the target URL scheme. This builder handles dns://*/* targets.
func (d *Builder) Scheme() string {
	return dnsResolverScheme
}

// Build returns a resolver that periodically resolves the input target. Each client connection
// maintains a resolver. It's a part of client connection's life cycle. The target follows
// gRPC desc resolution format (https://github.com/grpc/grpc/blob/master/doc/naming.md). As this
// builds a DNS resolver, we care about dns URL only: dns:[//authority/]host[:port]
// If the authority is missing (dns:host[:port]), it fallbacks to use OS resolver.
func (d *Builder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	path := target.URL.Path
	if path == "" {
		path = target.URL.Opaque
	}
	host, port, err := parseTarget(strings.TrimPrefix(path, "/"), d.opts.DefaultGrpcPort)
	if err != nil {
		return nil, structerr.New("building dns resolver: %w", err).WithMetadata("target", target.URL.String())
	}

	if addr, ok := tryParseIP(host, port); ok {
		// When the address is a static IP, we don't need this resolver anymore. Client
		// connection is responsible for handling network error in this case.
		_ = cc.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: addr}}})
		return &noopResolver{}, nil
	}

	authorityFinder := findDNSLookup
	if d.opts.authorityFinder != nil {
		authorityFinder = d.opts.authorityFinder
	}
	lookup, err := authorityFinder(target.URL.Host)
	if err != nil {
		return nil, structerr.New("finding DNS resolver: %w", err).WithMetadata("authority", target.URL.Host)
	}

	ctx, cancel := context.WithCancel(context.Background())
	dr := &dnsResolver{
		retry: d.opts.Backoff,

		ctx:         ctx,
		cancel:      cancel,
		host:        host,
		port:        port,
		cc:          cc,
		refreshRate: d.opts.RefreshRate,
		lookup:      lookup,
		reqs:        make(chan struct{}, 1),
	}

	dr.wg.Add(1)
	go dr.watch()

	return dr, nil
}
