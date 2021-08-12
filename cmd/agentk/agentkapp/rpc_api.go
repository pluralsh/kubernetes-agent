package agentkapp

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modshared"
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

func (a *agentRpcApi) HandleSendError(log *zap.Logger, msg string, err error) error {
	return handleSendError(log, msg, err)
}