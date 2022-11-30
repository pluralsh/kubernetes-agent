package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
)

func TestRpc_NewRpcRef(t *testing.T) {
	tcs := []struct {
		name        string
		ref         *agentcfg.GitRefCF
		expectedRef *GitRefCF
	}{
		{
			name:        "no ref should default to HEAD",
			ref:         nil,
			expectedRef: nil,
		},
		{
			name:        "empty ref should default to HEAD",
			ref:         &agentcfg.GitRefCF{},
			expectedRef: nil,
		},
		{
			name: "resolve arbitrary branch",
			ref: &agentcfg.GitRefCF{
				Ref: &agentcfg.GitRefCF_Branch{Branch: "any-branch-name"},
			},
			expectedRef: &GitRefCF{
				Ref: &GitRefCF_Branch{Branch: "any-branch-name"},
			},
		},
		{
			name: "resolve arbitrary tag",
			ref: &agentcfg.GitRefCF{
				Ref: &agentcfg.GitRefCF_Tag{Tag: "any-tag-name"},
			},
			expectedRef: &GitRefCF{
				Ref: &GitRefCF_Tag{Tag: "any-tag-name"},
			},
		},
		{
			name: "resolve arbitrary commit",
			ref: &agentcfg.GitRefCF{
				Ref: &agentcfg.GitRefCF_Commit{Commit: "any-commit-sha"},
			},
			expectedRef: &GitRefCF{
				Ref: &GitRefCF_Commit{Commit: "any-commit-sha"},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actualRef := NewRpcRef(tc.ref)
			assert.Equal(t, tc.expectedRef, actualRef)
		})
	}
}

func TestRpc_ResolveRef(t *testing.T) {
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
