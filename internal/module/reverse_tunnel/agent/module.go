package agent

import (
	"context"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
)

type module struct {
	server            *grpc.Server
	numConnections    int
	connectionFactory func(*info.AgentDescriptor) connectionInterface // helps testing
}

func (m *module) Run(ctx context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	descriptor := m.agentDescriptor()
	var wg wait.Group
	defer wg.Wait()
	for i := 0; i < m.numConnections; i++ {
		conn := m.connectionFactory(descriptor)
		wg.StartWithContext(ctx, conn.Run)
	}
	<-ctx.Done()
	return nil
}

func (m *module) DefaultAndValidateConfiguration(config *agentcfg.AgentConfiguration) error {
	return nil
}

func (m *module) Name() string {
	return reverse_tunnel.ModuleName
}

func (m *module) agentDescriptor() *info.AgentDescriptor {
	serverInfo := m.server.GetServiceInfo()
	services := make([]*info.Service, 0, len(serverInfo))
	for svcName, svcInfo := range serverInfo {
		methods := make([]*info.Method, 0, len(svcInfo.Methods))
		for _, mInfo := range svcInfo.Methods {
			methods = append(methods, &info.Method{
				Name: mInfo.Name,
			})
		}
		services = append(services, &info.Service{
			Name:    svcName,
			Methods: methods,
		})
	}
	return &info.AgentDescriptor{
		Services: services,
	}
}
