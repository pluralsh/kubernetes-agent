package main

import (
	"github.com/pluralsh/kuberentes-agent/cmd"
	"github.com/pluralsh/kuberentes-agent/cmd/kas/kasapp"
)

func main() {
	cmd.Run(kasapp.NewCommand())
}
