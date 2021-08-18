package modserver

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type agentRpcApiKeyType int

const (
	agentRpcApiKey agentRpcApiKeyType = iota
)

type AgentRpcApiFactory func(ctx context.Context, fullMethodName string) (AgentRpcApi, error)

func InjectAgentRpcApi(ctx context.Context, rpcApi AgentRpcApi) context.Context {
	return context.WithValue(ctx, agentRpcApiKey, rpcApi)
}

func AgentRpcApiFromContext(ctx context.Context) AgentRpcApi {
	rpcApi, ok := ctx.Value(agentRpcApiKey).(AgentRpcApi)
	if !ok {
		// This is a programmer error, so panic.
		panic("modserver.AgentRpcApi not attached to context. Make sure you are using interceptors")
	}
	return rpcApi
}

// UnaryAgentRpcApiInterceptor returns a new unary server interceptor that augments connection context with a AgentRpcApi.
func UnaryAgentRpcApiInterceptor(factory AgentRpcApiFactory) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		rpcApi, err := factory(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(InjectAgentRpcApi(ctx, rpcApi), req)
	}
}

// StreamAgentRpcApiInterceptor returns a new stream server interceptor that augments connection context with a AgentRpcApi.
func StreamAgentRpcApiInterceptor(factory AgentRpcApiFactory) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := grpc_middleware.WrapServerStream(ss)
		rpcApi, err := factory(wrapper.WrappedContext, info.FullMethod)
		if err != nil {
			return err
		}
		wrapper.WrappedContext = InjectAgentRpcApi(wrapper.WrappedContext, rpcApi)
		return handler(srv, wrapper)
	}
}
