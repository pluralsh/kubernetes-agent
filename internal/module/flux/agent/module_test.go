package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
)

func TestModule_DefaultAndValidateConfiguration_WithoutFluxConfig(t *testing.T) {
	// GIVEN
	m := &module{}
	cfg := &agentcfg.AgentConfiguration{}

	// WHEN
	err := m.DefaultAndValidateConfiguration(cfg)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, defaultServiceApiBaseUrl, cfg.Flux.WebhookReceiverUrl)
}

func TestModule_DefaultAndValidateConfiguration_WithoutWebhookReceiverUrlConfig(t *testing.T) {
	// GIVEN
	m := &module{}
	cfg := &agentcfg.AgentConfiguration{
		Flux: &agentcfg.FluxCF{},
	}

	// WHEN
	err := m.DefaultAndValidateConfiguration(cfg)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, defaultServiceApiBaseUrl, cfg.Flux.WebhookReceiverUrl)
}

func TestModule_DefaultAndValidateConfiguration_WithWebhookReceiverUrlConfig(t *testing.T) {
	// GIVEN
	m := &module{}
	cfg := &agentcfg.AgentConfiguration{
		Flux: &agentcfg.FluxCF{
			WebhookReceiverUrl: "https://example.com",
		},
	}

	// WHEN
	err := m.DefaultAndValidateConfiguration(cfg)

	// THEN
	require.NoError(t, err)
	assert.Equal(t, "https://example.com", cfg.Flux.WebhookReceiverUrl)
}
