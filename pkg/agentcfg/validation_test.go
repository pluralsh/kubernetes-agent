package agentcfg

import (
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
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
			ErrString: "invalid ExtraKeyValCF.Key: value length must be at least 1 bytes",
			Invalid:   &ExtraKeyValCF{},
		},
		{
			ErrString: "invalid ExtraKeyValCF.Val: value must contain at least 1 item(s)",
			Invalid: &ExtraKeyValCF{
				Key: "1",
			},
		},
	}
	testhelpers.AssertInvalid(t, tests)
}
