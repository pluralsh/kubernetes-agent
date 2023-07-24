package server

import (
	"time"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/kascfg"
)

const (
	defaultAgentConfigurationPollPeriod               = 5 * time.Minute
	defaultAgentConfigurationMaxConfigurationFileSize = 128 * 1024
)

func ApplyDefaults(config *kascfg.ConfigurationFile) {
	prototool.NotNil(&config.Agent)
	prototool.NotNil(&config.Agent.Configuration)
	prototool.NotNil(&config.Agent.Listen)

	c := config.Agent.Configuration
	prototool.Duration(&c.PollPeriod, defaultAgentConfigurationPollPeriod)
	prototool.Uint32(&c.MaxConfigurationFileSize, defaultAgentConfigurationMaxConfigurationFileSize)
}
