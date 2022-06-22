package gitaly

import (
	"context"
	"errors"
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitaly/copied/stats"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/ioz"
	"gitlab.com/gitlab-org/gitaly/v15/proto/go/gitalypb"
)

const (
	DefaultBranch = ""
)

// PollerInterface does the following:
// - polls ref advertisement for updates to the repository
// - detects which is the main branch, if branch or tag name is not specified
// - compares the commit id the branch or tag is referring to with the last processed one
// - returns the information about the change
type PollerInterface interface {
	// Poll performs a poll on the repository.
	// revision can be a branch name or a tag.
	// Poll returns a wrapped context.Canceled, context.DeadlineExceeded or gRPC error if ctx signals done and interrupts a running gRPC call.
	// Poll returns *Error when a error occurs.
	Poll(ctx context.Context, repo *gitalypb.Repository, lastProcessedCommitId, refName string) (*PollInfo, error)
}

type Poller struct {
	Client   gitalypb.SmartHTTPServiceClient
	Features map[string]string
}

type PollInfo struct {
	CommitId        string
	UpdateAvailable bool
	// EmptyRepository is true when polling the default branch but no refs were found.
	// When polling non-default branch and no refs were found the NotFound error is returned.
	EmptyRepository bool
}

func (p *Poller) Poll(ctx context.Context, repo *gitalypb.Repository, lastProcessedCommitId, refName string) (*PollInfo, error) {
	refNameTag := "refs/tags/" + refName
	refNameBranch := "refs/heads/" + refName
	isEmpty := true
	var head, master, wanted *stats.Reference
	err := p.fetchRefs(ctx, repo, func(ref stats.Reference) bool {
		isEmpty = false
		switch string(ref.Name) {
		case refNameTag, refNameBranch:
			wanted = cloneReference(ref)
			return true
		case "HEAD":
			head = cloneReference(ref)
		case "refs/heads/master":
			master = cloneReference(ref)
		}
		return false
	})
	if err != nil {
		return nil, err // don't wrap
	}
	if wanted == nil { // not found
		if refName != DefaultBranch { // was looking for something specific, but didn't find it
			return nil, NewNotFoundError("InfoRefsUploadPack", refName)
		}
		// was looking for the default branch
		switch {
		case head != nil:
			wanted = head
		case master != nil:
			wanted = master
		case isEmpty:
			return &PollInfo{
				EmptyRepository: true,
			}, nil
		default:
			return nil, NewNotFoundError("InfoRefsUploadPack", "default branch")
		}
	}
	oid := string(wanted.Oid)
	return &PollInfo{
		CommitId:        oid,
		UpdateAvailable: oid != lastProcessedCommitId,
	}, nil
}

// fetchRefs returns a wrapped context.Canceled, context.DeadlineExceeded or gRPC error if ctx signals done and interrupts a running gRPC call.
// fetchRefs returns *Error when a error occurs.
func (p *Poller) fetchRefs(ctx context.Context, repo *gitalypb.Repository, cb stats.ReferenceCb) error {
	ctx, cancel := context.WithCancel(appendFeatureFlagsToContext(ctx, p.Features))
	defer cancel() // ensure streaming call is canceled
	uploadPackReq := &gitalypb.InfoRefsRequest{
		Repository: repo,
		// Do not set GitConfigOptions or GitProtocol because that would disable cache in Gitaly.
		// See https://gitlab.com/gitlab-org/gitaly/-/blob/bea500b301bbec8535fbcae58c1da2d29377c666/internal/gitaly/service/smarthttp/cache.go#L53-56
	}
	uploadPackResp, err := p.Client.InfoRefsUploadPack(ctx, uploadPackReq)
	if err != nil {
		return NewRpcError(err, "InfoRefsUploadPack", "")
	}
	err = stats.ParseReferenceDiscovery(ioz.NewReceiveReader(func() ([]byte, error) {
		entry, err := uploadPackResp.Recv() // nolint: govet
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, io.EOF
			}
			return nil, NewRpcError(err, "InfoRefsUploadPack.Recv", "")
		}
		return entry.Data, nil
	}), cb)
	if err != nil {
		if _, ok := err.(*Error); ok { // nolint: errorlint
			return err // A wrapped error already
		}
		return NewProtocolError(err, "failed to parse reference discovery", "", "")
	}
	return nil
}

func cloneReference(ref stats.Reference) *stats.Reference {
	return &stats.Reference{
		Oid:  cloneSlice(ref.Oid),
		Name: cloneSlice(ref.Name),
	}
}

func cloneSlice(in []byte) []byte {
	out := make([]byte, len(in))
	copy(out, in)
	return out
}
