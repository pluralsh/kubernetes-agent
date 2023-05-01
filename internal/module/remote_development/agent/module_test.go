package agent

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/util/wait"
)

type mockReconciler struct {
	timesCalled uint32
}

func (r *mockReconciler) Stop() {
	// do nothing
}

func (r *mockReconciler) Run(_ context.Context) error {
	r.timesCalled += 1
	return nil
}

func TestConfigChange(t *testing.T) {
	tests := []struct {
		description string
		timesCalled uint32
		configs     []*agentcfg.AgentConfiguration
	}{
		{
			description: "When remote dev is enabled in the config, does start reconciler",
			timesCalled: 1,
			configs: []*agentcfg.AgentConfiguration{
				{
					RemoteDevelopment: &agentcfg.RemoteCF{
						Enabled: true,
					},
				},
			},
		},
		{
			description: "When the config is updated, restarts the reconciler",
			timesCalled: 2,
			configs: []*agentcfg.AgentConfiguration{
				{
					RemoteDevelopment: &agentcfg.RemoteCF{
						Enabled: true,
						DnsZone: "one",
					},
				},
				{
					RemoteDevelopment: &agentcfg.RemoteCF{
						Enabled: true,
						DnsZone: "two",
					},
				},
			},
		},
		{
			description: "When the config is published multiple times without any changes",
			timesCalled: 1,
			configs: []*agentcfg.AgentConfiguration{
				{
					RemoteDevelopment: &agentcfg.RemoteCF{
						Enabled: true,
						DnsZone: "one",
					},
				},
				{
					RemoteDevelopment: &agentcfg.RemoteCF{
						Enabled: true,
						DnsZone: "one",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			configChannel := make(chan *agentcfg.AgentConfiguration)
			mock := &mockReconciler{}

			mod := module{
				log: zaptest.NewLogger(t),
				api: newMockApi(t),
				reconcilerFactory: func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevReconciler, error) {
					return mock, nil
				},
			}

			ctx := context.Background()

			wg := wait.Group{}

			// publish configs asynchronously
			wg.StartWithContext(ctx, func(ctx context.Context) {
				publishInterval := 50 * time.Millisecond

				if tt.configs != nil {
					for _, cfg := range tt.configs {
						// populate the test config with defaults if missing
						// this must be explicitly done in tests where module's Run() is invoked directly
						err := mod.DefaultAndValidateConfiguration(cfg)
						require.NoError(t, err)

						configChannel <- cfg
						time.Sleep(publishInterval)
					}
				}
				close(configChannel)
			})

			err := mod.Run(ctx, configChannel)
			wg.Wait()
			require.NoError(t, err)
			require.Equal(t, tt.timesCalled, mock.timesCalled)
		})
	}
}
