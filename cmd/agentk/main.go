package main

import (
	"github.com/pluralsh/kuberentes-agent/cmd"
	"github.com/pluralsh/kuberentes-agent/cmd/agentk/agentkapp"
)

func main() {
	cmd.Run(agentkapp.NewCommand())
}
