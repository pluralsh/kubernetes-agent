package server

import (
	"github.com/pluralsh/kuberentes-agent/internal/module/modserver"
)

var (
	_ modserver.Module  = &module{}
	_ modserver.Factory = &Factory{}
)
