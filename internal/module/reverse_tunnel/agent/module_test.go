package agent

import "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"

var (
	_ modagent.Module = &module{}
)
