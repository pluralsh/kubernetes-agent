package gitaly

import (
	"context"
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/gitalypb"
)

type PathEntryVisitor interface {
	Entry(*gitalypb.TreeEntry) (bool /* done? */, error)
}

type PathVisitor struct {
	Client   gitalypb.CommitServiceClient
	Features map[string]string
}

func (v *PathVisitor) Visit(ctx context.Context, repo *gitalypb.Repository, revision, repoPath []byte, recursive bool, visitor PathEntryVisitor) error {
	ctx, cancel := context.WithCancel(appendFeatureFlagsToContext(ctx, v.Features))
	defer cancel() // ensure streaming call is canceled
	entries, err := v.Client.GetTreeEntries(ctx, &gitalypb.GetTreeEntriesRequest{
		Repository: repo,
		Revision:   revision,
		Path:       repoPath,
		Recursive:  recursive,
	})
	if err != nil {
		return NewRpcError(err, "GetTreeEntries", string(repoPath))
	}
entriesLoop:
	for {
		resp, err := entries.Recv()
		if err != nil {
			if err == io.EOF { // nolint:errorlint
				break
			}
			return NewRpcError(err, "GetTreeEntries.Recv", string(repoPath))
		}
		for _, entry := range resp.Entries {
			done, err := visitor.Entry(entry)
			if err != nil {
				return err // don't wrap
			}
			if done {
				break entriesLoop
			}
		}
	}
	return nil
}
