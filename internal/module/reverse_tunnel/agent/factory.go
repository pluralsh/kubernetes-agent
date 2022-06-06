package agent

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/grpctool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/retry"
	"google.golang.org/grpc"
)

const (
	minIdleConnections = 2
	maxConnections     = 100
	maxIdleTime        = time.Minute
	// scaleUpStep defines how many new connections are started when there is not enough idle connections.
	scaleUpStep = 10

	connectionInitBackoff   = 10 * time.Second
	connectionMaxBackoff    = 5 * time.Minute
	connectionResetDuration = 10 * time.Minute
	connectionBackoffFactor = 2.0
	connectionJitter        = 1.0
)

type Factory struct {
	InternalServerConn grpc.ClientConnInterface
}

func (f *Factory) New(config *modagent.Config) (modagent.Module, error) {
	sv, err := grpctool.NewStreamVisitor(&rpc.ConnectResponse{})
	if err != nil {
		return nil, err
	}
	client := rpc.NewReverseTunnelClient(config.KasConn)
	pollConfig := retry.NewPollConfigFactory(0, retry.NewExponentialBackoffFactory(
		connectionInitBackoff,
		connectionMaxBackoff,
		connectionResetDuration,
		connectionBackoffFactor,
		connectionJitter,
	))
	return &module{
		server:             config.Server,
		minIdleConnections: minIdleConnections,
		maxConnections:     maxConnections,
		scaleUpStep:        scaleUpStep,
		maxIdleTime:        maxIdleTime,
		connectionFactory: func(descriptor *info.AgentDescriptor, onActive, onIdle func(c connectionInterface)) connectionInterface {
			return &connection{
				log:                config.Log,
				descriptor:         descriptor,
				client:             client,
				internalServerConn: f.InternalServerConn,
				streamVisitor:      sv,
				pollConfig:         pollConfig,
				onActive:           onActive,
				onIdle:             onIdle,
			}
		},
	}, nil
}

func (f *Factory) Name() string {
	return reverse_tunnel.ModuleName
}

func (f *Factory) UsesInternalServer() bool {
	return true
}
