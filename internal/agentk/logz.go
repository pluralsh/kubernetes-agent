package agentk

import (
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tools/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/pkg/agentcfg"
	"go.uber.org/zap"
)

func engineSyncResult(syncResult string) zap.Field {
	return zap.String(logz.EngineSyncResult, syncResult)
}

func engineResourceKey(resourceKey kube.ResourceKey) zap.Field {
	return zap.Stringer(logz.EngineResourceKey, &resourceKey)
}

func agentConfig(config *agentcfg.AgentConfiguration) zap.Field {
	return zap.Reflect(logz.AgentConfig, config)
}
