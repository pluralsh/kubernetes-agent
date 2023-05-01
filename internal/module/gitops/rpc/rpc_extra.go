package rpc

import (
	"fmt"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
)

func NewRpcRef(configRef *agentcfg.GitRefCF) *GitRefCF {
	switch ref := configRef.GetRef().(type) {
	case *agentcfg.GitRefCF_Tag:
		return &GitRefCF{Ref: &GitRefCF_Tag{Tag: ref.Tag}}
	case *agentcfg.GitRefCF_Branch:
		return &GitRefCF{Ref: &GitRefCF_Branch{Branch: ref.Branch}}
	case *agentcfg.GitRefCF_Commit:
		return &GitRefCF{Ref: &GitRefCF_Commit{Commit: ref.Commit}}
	case nil:
		return nil
	default:
		// Nah, this doesn't happen - UNLESS you forgot to add a `case` when changing the `agentcfg.GitRefCF` proto message ;)
		panic(fmt.Sprintf("unexpected ref to resolve: %T", ref))
	}
}

// GetResolvedRef resolved the `Ref` into a full unambiguous Git reference.
func (x *GitRefCF) GetResolvedRef() string {
	switch ref := x.GetRef().(type) {
	case *GitRefCF_Tag:
		return "refs/tags/" + ref.Tag
	case *GitRefCF_Branch:
		return "refs/heads/" + ref.Branch
	case *GitRefCF_Commit:
		return ref.Commit
	case nil:
		// as a default and for backward-compatibility reasons we assume that if no ref is specified the default project branch is used.
		return gitaly.DefaultBranch
	default:
		// Nah, this doesn't happen - UNLESS you forgot to add a `case` when changing the `agentcfg.GitRefCF` proto message ;)
		panic(fmt.Sprintf("unexpected ref to resolve: %T", ref))
	}
}
