package manifestops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ modagent.LeaderModule = &module{}
	_ modagent.Factory      = &Factory{}
)

func TestIgnoresInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		config *agentcfg.AgentConfiguration
	}{
		{
			name: "duplicate project ids",
			config: &agentcfg.AgentConfiguration{
				Gitops: &agentcfg.GitopsCF{
					ManifestProjects: []*agentcfg.ManifestProjectCF{
						{
							Id: "project1",
						},
						{
							Id: "project1",
							Paths: []*agentcfg.PathCF{
								{
									Glob: "*.yaml",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &module{
				log: zaptest.NewLogger(t),
			}
			var wg wait.Group
			defer wg.Wait()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			cfg := make(chan *agentcfg.AgentConfiguration)
			wg.Start(func() {
				err := m.Run(ctx, cfg)
				assert.NoError(t, err)
			})
			require.NoError(t, m.DefaultAndValidateConfiguration(tc.config))
			cfg <- tc.config
			close(cfg)
			wg.Wait()
		})
	}
}
