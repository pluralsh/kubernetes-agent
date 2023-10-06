package agent

import (
	"context"

	"cloud.google.com/go/profiler"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/google_profiler"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/prototool"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/proto"
)

type module struct {
	log    *zap.Logger
	runner profilerRunner
}

type profilerRunner interface {
	start(cfg *agentcfg.GoogleProfilerCF) error
}

type googleProfilerRunner struct {
	service string
	version string
}

func (m *module) Run(_ context.Context, cfg <-chan *agentcfg.AgentConfiguration) error {
	// NOTE: The Google profiler cannot be stopped or restarted. Therefore, the first config
	// received that enables the profiler is to last to take effect. To stop or restart the profiler
	// the agent must be restarted.
	// See https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/469#note_1585658837

	// Wait for the first config that enables and successfully starts the profiler
	var configInUse *agentcfg.GoogleProfilerCF
	for config := range cfg {
		profilerCfg := config.Observability.GoogleProfiler
		if !profilerCfg.Enabled {
			continue
		}

		if err := m.runner.start(profilerCfg); err != nil {
			m.log.Error("Failed to start profiler", logz.Error(err))
			continue
		}

		configInUse = profilerCfg
		m.log.Info("Started profiler. Changes to the observability.google_profiler config won't affect this agent until it's restarted")
		break
	}

	for config := range cfg {
		if !proto.Equal(config.Observability.GoogleProfiler, configInUse) {
			m.log.Warn("The observability.google_profiler config has changed, but a previous configuration of the profiler is already running. Please restart the agent for it to have any effect")
		}
	}
	return nil
}

func (r *googleProfilerRunner) start(cfg *agentcfg.GoogleProfilerCF) error {
	config := profiler.Config{
		Service:        r.service,
		ServiceVersion: r.version,
		DebugLogging:   cfg.DebugLogging,
		MutexProfiling: true, // like in LabKit
		ProjectID:      cfg.ProjectId,
	}
	var opts []option.ClientOption
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}
	return profiler.Start(config, opts...)
}

func (m *module) DefaultAndValidateConfiguration(cfg *agentcfg.AgentConfiguration) error {
	prototool.NotNil(&cfg.Observability)
	prototool.NotNil(&cfg.Observability.GoogleProfiler)
	return nil
}

func (m *module) Name() string {
	return google_profiler.ModuleName
}
