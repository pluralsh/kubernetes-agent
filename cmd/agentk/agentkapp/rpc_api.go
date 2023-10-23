package agentkapp

import (
	"go.uber.org/zap"

	"github.com/pluralsh/kuberentes-agent/pkg/module/modagent"
	"github.com/pluralsh/kuberentes-agent/pkg/module/modshared"
	"github.com/pluralsh/kuberentes-agent/pkg/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/pkg/tool/logz"
)

var (
	_ modagent.RpcApi = (*agentRpcApi)(nil)
)

type agentRpcApi struct {
	modshared.RpcApiStub
}

func (a *agentRpcApi) HandleProcessingError(log *zap.Logger, agentId int64, msg string, err error) {
	handleProcessingError(a.StreamCtx, log, agentId, msg, err)
}

func (a *agentRpcApi) HandleIoError(log *zap.Logger, msg string, err error) error {
	// The problem is almost certainly with the client's connection.
	// Still log it on Debug.
	log.Debug(msg, logz.Error(err))
	return grpctool.HandleIoError(msg, err)
}
