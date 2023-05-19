package dnsresolver

import (
	"context"
	"net"
	"net/url"
	"strings"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/structerr"
)

// ValidateURL validates Gitaly address URL having dns scheme. The URL follows three forms:
// * dns://authority-port:authority-host/host:port
// * dns:///host:port
// * dns:host:port
// Either form, the real address is the URL's path
func ValidateURL(rawAddress string) error {
	if rawAddress == "" {
		return structerr.New("empty address")
	}

	uri, err := url.Parse(rawAddress)
	if err != nil {
		return structerr.New("fail to parse address: %w", err)
	}

	if uri.Scheme != "dns" {
		return structerr.New("unexpected scheme: %s", uri.Scheme)
	}

	path := uri.Path
	if path == "" {
		// When "//" part is stripped
		path = uri.Opaque
	}
	_, _, err = parseTarget(strings.TrimPrefix(path, "/"), "50051")
	return err
}

// parseTarget takes the user input target string and default port, returns formatted host and port info.
// This is a shameless copy of built-in gRPC dns resolver, because we don't want to have any
// inconsistency between our resolver and dns resolver.
// Source: https://github.com/grpc/grpc-go/blob/eeb9afa1f6b6388152955eeca8926e36ca94c768/internal/resolver/dns/dns_resolver.go#L378-L378
func parseTarget(target, defaultPort string) (string, string, error) {
	var err error

	if target == "" {
		return "", "", structerr.New("dns resolver: missing address")
	}

	if ip := net.ParseIP(target); ip != nil {
		// target is an IPv4 or IPv6(without brackets) address
		return target, defaultPort, nil
	}

	if host, port, err := net.SplitHostPort(target); err == nil {
		if port == "" {
			// If the port field is empty (target ends with colon), e.g. "[::1]:", this is an error.
			return "", "", structerr.New("dns resolver: missing port after port-separator colon")
		}
		if host == "" {
			host = "localhost"
		}
		return host, port, nil
	}

	host, port, err := net.SplitHostPort(target + ":" + defaultPort)
	if err == nil {
		return host, port, nil
	}
	return "", "", structerr.New("dns resolver: %w", err)
}

func tryParseIP(host, port string) (addr string, ok bool) {
	ip := net.ParseIP(host)
	if ip == nil {
		return "", false
	}
	return net.JoinHostPort(host, port), true
}

func findDNSLookup(authority string) (dnsLookuper, error) {
	if authority == "" {
		return net.DefaultResolver, nil
	}

	host, port, err := parseTarget(authority, defaultDNSNameserverPort)
	if err != nil {
		return nil, err
	}

	addr := net.JoinHostPort(host, port)
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			var dialer net.Dialer
			return dialer.DialContext(ctx, network, addr)
		},
	}, nil
}
