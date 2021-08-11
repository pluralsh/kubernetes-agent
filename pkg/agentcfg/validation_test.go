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
			Name: "empty CiAccessGroup.DefaultNamespace",
			Valid: &CiAccessGroup{
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
			ErrString: "invalid CiAccessGroup.Id: value length must be at least 1 bytes",
			Invalid: &CiAccessGroup{
				Id: "", // empty id is not ok
			},
		},
	}
	testhelpers.AssertInvalid(t, tests)
}
