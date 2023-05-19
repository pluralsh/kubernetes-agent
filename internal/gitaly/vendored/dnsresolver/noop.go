package dnsresolver

import "google.golang.org/grpc/resolver"

// noopResolver does nothing. It is used when a target is not resolvable.
type noopResolver struct{}

func (noopResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (noopResolver) Close() {}
