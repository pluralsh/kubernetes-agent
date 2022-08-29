package main

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd/agentk/agentkapp"
)

func main() {
	cmd.Run(agentkapp.NewCommand())
}
