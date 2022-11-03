package agentcfg

import (
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
)

func TestValidation_Valid(t *testing.T) {
	tests := []testhelpers.ValidTestcase{
		{
			Name:  "empty",
			Valid: &AgentConfiguration{},
		},
		{
			Name: "empty CiAccessGroupCF.DefaultNamespace",
			Valid: &CiAccessGroupCF{
				Id:               "abc",
				DefaultNamespace: "", // empty is ok
			},
		},
		{
			Name:  "empty CiAccessAsAgentCF",
			Valid: &CiAccessAsAgentCF{},
		},
		{
			Name:  "empty CiAccessAsCiJobCF",
			Valid: &CiAccessAsCiJobCF{},
		},
		{
			Name: "minimal CiAccessAsImpersonateCF",
			Valid: &CiAccessAsImpersonateCF{
				Username: "abc",
			},
		},
		{
			Name: "one group CiAccessAsImpersonateCF",
			Valid: &CiAccessAsImpersonateCF{
				Username: "abc",
				Groups:   []string{"g"},
			},
		},
	}
	testhelpers.AssertValid(t, tests)
}

func TestValidation_Invalid(t *testing.T) {
	tests := []testhelpers.InvalidTestcase{
		{
			ErrString: "invalid ManifestProjectCF.Id: value length must be at least 1 bytes",
			Invalid: &ManifestProjectCF{
				Id: "", // empty id is not ok
			},
		},
		{
			ErrString: "invalid PathCF.Glob: value length must be at least 1 bytes",
			Invalid: &PathCF{
				Glob: "",
			},
		},
		{
			ErrString: "invalid CiAccessGroupCF.Id: value length must be at least 1 bytes",
			Invalid: &CiAccessGroupCF{
				Id: "", // empty id is not ok
			},
		},
		{
			ErrString: "invalid CiAccessAsCF.As: value is required",
			Invalid:   &CiAccessAsCF{},
		},
		{
			ErrString: "invalid CiAccessAsCF.Agent: value is required",
			Invalid: &CiAccessAsCF{
				As: &CiAccessAsCF_Agent{},
			},
		},
		{
			ErrString: "invalid CiAccessAsCF.Impersonate: value is required",
			Invalid: &CiAccessAsCF{
				As: &CiAccessAsCF_Impersonate{},
			},
		},
		{
			ErrString: "invalid CiAccessAsCF.CiJob: value is required",
			Invalid: &CiAccessAsCF{
				As: &CiAccessAsCF_CiJob{},
			},
		},
		{
			ErrString: "invalid CiAccessAsImpersonateCF.Username: value length must be at least 1 bytes",
			Invalid:   &CiAccessAsImpersonateCF{},
		},
		{
			ErrString: "invalid CiAccessAsImpersonateCF.Groups[0]: value length must be at least 1 bytes",
			Invalid: &CiAccessAsImpersonateCF{
				Username: "a",
				Groups:   []string{""},
			},
		},
		{
			ErrString: "invalid ExtraKeyValCF.Key: value length must be at least 1 bytes; invalid ExtraKeyValCF.Val: value must contain at least 1 item(s)",
			Invalid:   &ExtraKeyValCF{},
		},
		{
			ErrString: "invalid ChartValuesCF.As: value is required",
			Invalid:   &ChartValuesCF{},
		},
	}
	testhelpers.AssertInvalid(t, tests)
}
