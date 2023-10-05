package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

var (
	_ modagent.Module = &module{}
)

func TestModule_Run_NoStartWhenNotEnabled(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockProfilerRunner := NewMockprofilerRunner(ctrl)
	m := module{
		log:    zaptest.NewLogger(t),
		runner: mockProfilerRunner,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg := make(chan *agentcfg.AgentConfiguration, 1)

	cfg <- &agentcfg.AgentConfiguration{
		Observability: &agentcfg.ObservabilityCF{
			GoogleProfiler: &agentcfg.GoogleProfilerCF{
				Enabled: false,
			},
		},
	}
	close(cfg)

	// THEN
	mockProfilerRunner.EXPECT().start(gomock.Any()).Times(0)

	// WHEN
	err := m.Run(ctx, cfg)
	require.NoError(t, err)
}

func TestModule_Run_StartWhenEnabled(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockProfilerRunner := NewMockprofilerRunner(ctrl)
	m := module{
		log:    zaptest.NewLogger(t),
		runner: mockProfilerRunner,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg := make(chan *agentcfg.AgentConfiguration, 1)

	profilerCfg := &agentcfg.GoogleProfilerCF{
		Enabled: true,
	}
	cfg <- &agentcfg.AgentConfiguration{
		Observability: &agentcfg.ObservabilityCF{
			GoogleProfiler: profilerCfg,
		},
	}
	close(cfg)

	// THEN
	mockProfilerRunner.EXPECT().start(matcher.ProtoEq(t, profilerCfg)).Times(1)

	// WHEN
	err := m.Run(ctx, cfg)
	require.NoError(t, err)
}

func TestModule_Run_StartWhenEnabledAndNotStartedAgain(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockProfilerRunner := NewMockprofilerRunner(ctrl)
	m := module{
		log:    zaptest.NewLogger(t),
		runner: mockProfilerRunner,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg := make(chan *agentcfg.AgentConfiguration, 2)

	profilerCfg1 := &agentcfg.GoogleProfilerCF{
		Enabled: true,
	}
	profilerCfg2 := &agentcfg.GoogleProfilerCF{
		Enabled:      true,
		DebugLogging: true,
	}
	cfg <- &agentcfg.AgentConfiguration{
		Observability: &agentcfg.ObservabilityCF{
			GoogleProfiler: profilerCfg1,
		},
	}
	cfg <- &agentcfg.AgentConfiguration{
		Observability: &agentcfg.ObservabilityCF{
			GoogleProfiler: profilerCfg2,
		},
	}
	close(cfg)

	// THEN
	mockProfilerRunner.EXPECT().start(matcher.ProtoEq(t, profilerCfg1)).Times(1)
	mockProfilerRunner.EXPECT().start(matcher.ProtoEq(t, profilerCfg2)).Times(0)

	// WHEN
	err := m.Run(ctx, cfg)
	require.NoError(t, err)
}

func TestModule_Run_StartWhenEnabledAfterPreviousStartFailure(t *testing.T) {
	// GIVEN
	ctrl := gomock.NewController(t)
	mockProfilerRunner := NewMockprofilerRunner(ctrl)
	m := module{
		log:    zaptest.NewLogger(t),
		runner: mockProfilerRunner,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg := make(chan *agentcfg.AgentConfiguration, 2)

	profilerCfgFails := &agentcfg.GoogleProfilerCF{
		Enabled: true,
	}
	profilerCfgSucceeds := &agentcfg.GoogleProfilerCF{
		Enabled:      true,
		DebugLogging: true,
	}
	cfg <- &agentcfg.AgentConfiguration{
		Observability: &agentcfg.ObservabilityCF{
			GoogleProfiler: profilerCfgFails,
		},
	}
	cfg <- &agentcfg.AgentConfiguration{
		Observability: &agentcfg.ObservabilityCF{
			GoogleProfiler: profilerCfgSucceeds,
		},
	}
	close(cfg)

	// THEN
	mockProfilerRunner.EXPECT().start(matcher.ProtoEq(t, profilerCfgFails)).Return(errors.New("dummy error"))
	mockProfilerRunner.EXPECT().start(matcher.ProtoEq(t, profilerCfgSucceeds)).Times(1)

	// WHEN
	err := m.Run(ctx, cfg)
	require.NoError(t, err)
}
