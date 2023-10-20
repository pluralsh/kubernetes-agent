package server

import (
	"context"
	"net"

	"github.com/pluralsh/kuberentes-agent/internal/module/kubernetes_api"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"go.uber.org/zap"
)

type module struct {
	log      *zap.Logger
	proxy    kubernetesApiProxy
	listener func() (net.Listener, error)
}

func (m *module) Run(ctx context.Context) error {
	lis, err := m.listener()
	if err != nil {
		return err
	}
	// Error is ignored because kubernetesApiProxy.Run() closes the listener and
	// a second close always produces an error.
	defer lis.Close() // nolint:errcheck,gosec

	m.log.Info("Kubernetes API endpoint is up",
		logz.NetNetworkFromAddr(lis.Addr()),
		logz.NetAddressFromAddr(lis.Addr()),
	)
	return m.proxy.Run(ctx, lis)
}

func (m *module) Name() string {
	return kubernetes_api.ModuleName
}

type nopModule struct {
}

func (m nopModule) Run(ctx context.Context) error {
	return nil
}

func (m nopModule) Name() string {
	return kubernetes_api.ModuleName
}
