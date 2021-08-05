package grpctool

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"google.golang.org/grpc"
)

type rpcApiKeyType int

const (
	rpcApiKey rpcApiKeyType = iota
)

type RpcApiFactory func(ctx context.Context, method string) modserver.RpcApi

func InjectRpcApi(ctx context.Context, rpcApi modserver.RpcApi) context.Context {
	return context.WithValue(ctx, rpcApiKey, rpcApi)
}

func RpcApiFromContext(ctx context.Context) modserver.RpcApi {
	rpcApi, ok := ctx.Value(rpcApiKey).(modserver.RpcApi)
	if !ok {
		// This is a programmer error, so panic.
		panic("modserver.RPCAPI not attached to context. Make sure you are using interceptors")
	}
	return rpcApi
}

// UnaryServerRpcApiInterceptor returns a new unary server interceptor that augments connection context with a modserver.RpcApi.
func UnaryServerRpcApiInterceptor(factory RpcApiFactory) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(augmentContextWithRpcApi(ctx, info.FullMethod, factory), req)
	}
}

// StreamServerRpcApiInterceptor returns a new stream server interceptor that augments connection context with a modserver.RpcApi.
func StreamServerRpcApiInterceptor(factory RpcApiFactory) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := grpc_middleware.WrapServerStream(ss)
		wrapper.WrappedContext = augmentContextWithRpcApi(wrapper.WrappedContext, info.FullMethod, factory)
		return handler(srv, wrapper)
	}
}

func augmentContextWithRpcApi(ctx context.Context, method string, factory RpcApiFactory) context.Context {
	return InjectRpcApi(ctx, factory(ctx, method))
}
