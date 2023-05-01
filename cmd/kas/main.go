package main

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/cmd/kas/kasapp"
)

func main() {
	cmd.Run(kasapp.NewCommand())
}
