package agent

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"go.uber.org/zap/zaptest"
)

type mockWorker struct {
	timesCalled uint32
}

func (w *mockWorker) Run(ctx context.Context) error {
	w.timesCalled += 1
	return nil
}

func TestConfigChange(t *testing.T) {
	tests := []struct {
		description string
		timesCalled uint32
		configs     []*agentcfg.AgentConfiguration
	}{
		{
			description: "When no config exists, does not start worker",
			timesCalled: 0,
			configs:     nil,
		},
		{
			description: "When remote dev is disabled in the config, does not start worker",
			timesCalled: 0,
			configs: []*agentcfg.AgentConfiguration{
				{
					RemoteDevelopment: &agentcfg.RemoteCF{
						Enabled: false,
					},
				},
			},
		},
		{
			description: "When remote dev is enabled in the config, does start worker",
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
			description: "When the config is updated, restarts the worker",
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
			description: "When the config is first enabled then disabled, calls the worker run method once",
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
						Enabled: false,
						DnsZone: "two",
					},
				},
			},
		},
		{
			description: "When remote dev config is not changed, does not restart the worker",
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
			configChannel := make(chan *agentcfg.AgentConfiguration, len(tt.configs))
			mock := &mockWorker{}
			ctrl := gomock.NewController(t)

			mod := module{
				log: zaptest.NewLogger(t),
				api: mock_modagent.NewMockApi(ctrl),
				workerFactory: func(ctx context.Context, cfg *agentcfg.RemoteCF) (remoteDevWorker, error) {
					return mock, nil
				},
			}

			ctx := context.Background()

			if tt.configs != nil {
				for _, c := range tt.configs {
					configChannel <- c
				}
			}
			close(configChannel)

			err := mod.Run(ctx, configChannel)
			require.NoError(t, err)
			require.Equal(t, tt.timesCalled, mock.timesCalled)
		})
	}
}
