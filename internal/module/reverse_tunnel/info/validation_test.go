package info

import (
	"testing"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
)

func TestValidation_Invalid(t *testing.T) {
	tests := []testhelpers.InvalidTestcase{
		{
			ErrString: "invalid Service.Name: value length must be at least 1 bytes",
			Invalid:   &Service{},
		},
		{
			ErrString: "invalid Method.Name: value length must be at least 1 bytes",
			Invalid:   &Method{},
		},
	}
	testhelpers.AssertInvalid(t, tests)
}
