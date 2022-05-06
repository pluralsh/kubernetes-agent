package main

import (
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/cmd/kas/kasapp"
)

func main() {
	cmd.Run(kasapp.NewCommand())
}
