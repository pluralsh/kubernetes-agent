package tracker

import (
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/reverse_tunnel/info"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
)

func TestValidation_Valid(t *testing.T) {
	tests := []testhelpers.ValidTestcase{
		{
			Name: "minimal/grpc",
			Valid: &TunnelInfo{
				AgentDescriptor: &info.AgentDescriptor{},
				KasUrl:          "grpc://1.1.1.1:10",
			},
		},
		{
			Name: "grpcs",
			Valid: &TunnelInfo{
				AgentDescriptor: &info.AgentDescriptor{},
				KasUrl:          "grpcs://1.1.1.1:10",
			},
		},
	}
	testhelpers.AssertValid(t, tests)
}

func TestValidation_Invalid(t *testing.T) {
	tests := []testhelpers.InvalidTestcase{
		{
			ErrString: "invalid TunnelInfo.AgentDescriptor: value is required; invalid TunnelInfo.KasUrl: value length must be at least 1 bytes",
			Invalid:   &TunnelInfo{},
		},
		{
			ErrString: `invalid TunnelInfo.KasUrl: value does not match regex pattern "(?:^$|^grpcs?://)"`,
			Invalid: &TunnelInfo{
				AgentDescriptor: &info.AgentDescriptor{},
				KasUrl:          "tcp://1.1.1.1:12",
			},
		},
	}
	testhelpers.AssertInvalid(t, tests)
}
