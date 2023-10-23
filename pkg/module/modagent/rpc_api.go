package modagent

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"google.golang.org/grpc"
)

type rpcApiKeyType int

const (
	rpcApiKey rpcApiKeyType = iota
)

type RpcApiFactory func(ctx context.Context, method string) RpcApi

func InjectRpcApi(ctx context.Context, rpcApi RpcApi) context.Context {
	return context.WithValue(ctx, rpcApiKey, rpcApi)
}

func RpcApiFromContext(ctx context.Context) RpcApi {
	rpcApi, ok := ctx.Value(rpcApiKey).(RpcApi)
	if !ok {
		// This is a programmer error, so panic.
		panic("modagent.RPCAPI not attached to context. Make sure you are using interceptors")
	}
	return rpcApi
}

// UnaryRpcApiInterceptor returns a new unary server interceptor that augments connection context with a RpcApi.
func UnaryRpcApiInterceptor(factory RpcApiFactory) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(augmentContextWithRpcApi(ctx, info.FullMethod, factory), req)
	}
}

// StreamRpcApiInterceptor returns a new stream server interceptor that augments connection context with a RpcApi.
func StreamRpcApiInterceptor(factory RpcApiFactory) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := grpc_middleware.WrapServerStream(ss)
		wrapper.WrappedContext = augmentContextWithRpcApi(wrapper.WrappedContext, info.FullMethod, factory)
		return handler(srv, wrapper)
	}
}

func augmentContextWithRpcApi(ctx context.Context, method string, factory RpcApiFactory) context.Context {
	return InjectRpcApi(ctx, factory(ctx, method))
}
