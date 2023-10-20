package agentkapp

import (
	"github.com/pluralsh/kuberentes-agent/internal/module/modagent"
	"github.com/pluralsh/kuberentes-agent/internal/module/modshared"
	"github.com/pluralsh/kuberentes-agent/internal/tool/grpctool"
	"github.com/pluralsh/kuberentes-agent/internal/tool/logz"
	"go.uber.org/zap"
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
