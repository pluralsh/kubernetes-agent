package agentcfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveRef(t *testing.T) {
	tcs := []struct {
		name                string
		cf                  *GitRefCF
		expectedResolvedRef string
	}{
		{
			name:                "no ref should default to HEAD",
			cf:                  nil,
			expectedResolvedRef: "HEAD",
		},
		{
			name:                "empty ref should default to HEAD",
			cf:                  &GitRefCF{},
			expectedResolvedRef: "HEAD",
		},
		{
			name: "resolve arbitrary branch",
			cf: &GitRefCF{
				Ref: &GitRefCF_Branch{Branch: "any-branch-name"},
			},
			expectedResolvedRef: "refs/heads/any-branch-name",
		},
		{
			name: "resolve arbitrary tag",
			cf: &GitRefCF{
				Ref: &GitRefCF_Tag{Tag: "any-tag-name"},
			},
			expectedResolvedRef: "refs/tags/any-tag-name",
		},
		{
			name: "resolve arbitrary commit",
			cf: &GitRefCF{
				Ref: &GitRefCF_Commit{Commit: "any-commit-sha"},
			},
			expectedResolvedRef: "any-commit-sha",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actualResolvedRef := tc.cf.GetResolvedRef()
			assert.Equal(t, tc.expectedResolvedRef, actualResolvedRef)
		})
	}
}
