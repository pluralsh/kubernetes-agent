package api

import (
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
)

func TestValidation_Valid(t *testing.T) {
	tests := []testhelpers.ValidTestcase{
		{
			Name:  "empty Configuration",
			Valid: &Configuration{},
		},
		{
			Name: "minimal1 Configuration",
			Valid: &Configuration{
				DefaultNamespace: "def",
			},
		},
		{
			Name: "minimal2 Configuration",
			Valid: &Configuration{
				DefaultNamespace: "def",
				AccessAs: &agentcfg.CiAccessAsCF{
					As: &agentcfg.CiAccessAsCF_Agent{
						Agent: &agentcfg.CiAccessAsAgentCF{},
					},
				},
			},
		},
		{
			Name: "minimal AllowedAgent",
			Valid: &AllowedAgent{
				ConfigProject: &ConfigProject{},
			},
		},
		{
			Name:  "minimal ConfigProject",
			Valid: &ConfigProject{},
		},
		{
			Name:  "minimal Pipeline",
			Valid: &Pipeline{},
		},
		{
			Name:  "minimal Project",
			Valid: &Project{},
		},
		{
			Name: "Project with groups",
			Valid: &Project{
				Groups: []*Group{
					{},
				},
			},
		},
		{
			Name:  "minimal Group",
			Valid: &Group{},
		},
		{
			Name:  "minimal Job",
			Valid: &Job{},
		},
		{
			Name: "minimal User",
			Valid: &User{
				Username: "user",
			},
		},
		{
			Name: "minimal Environment",
			Valid: &Environment{
				Slug: "prod",
			},
		},
		{
			Name: "minimal AllowedAgentsForJob",
			Valid: &AllowedAgentsForJob{
				Job:      &Job{},
				Pipeline: &Pipeline{},
				Project:  &Project{},
				User: &User{
					Username: "user",
				},
			},
		},
	}
	testhelpers.AssertValid(t, tests)
}

func TestValidation_Invalid(t *testing.T) {
	tests := []testhelpers.InvalidTestcase{
		{
			ErrString: "invalid AllowedAgent.ConfigProject: value is required",
			Invalid:   &AllowedAgent{},
		},
		{
			ErrString: "invalid User.Username: value length must be at least 1 bytes",
			Invalid:   &User{},
		},
		{
			ErrString: "invalid Environment.Slug: value length must be at least 1 bytes",
			Invalid:   &Environment{},
		},
		{
			ErrString: "invalid AllowedAgentsForJob.Job: value is required",
			Invalid:   &AllowedAgentsForJob{},
		},
		{
			ErrString: "invalid AllowedAgentsForJob.Pipeline: value is required",
			Invalid: &AllowedAgentsForJob{
				Job: &Job{},
			},
		},
		{
			ErrString: "invalid AllowedAgentsForJob.Project: value is required",
			Invalid: &AllowedAgentsForJob{
				Job:      &Job{},
				Pipeline: &Pipeline{},
			},
		},
		{
			ErrString: "invalid AllowedAgentsForJob.User: value is required",
			Invalid: &AllowedAgentsForJob{
				Job:      &Job{},
				Pipeline: &Pipeline{},
				Project:  &Project{},
			},
		},
	}
	testhelpers.AssertInvalid(t, tests)
}
