package gitaly

import (
	"bytes"
	"context"
	"io"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/gitalypb"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitaly/vendored/stats"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/ioz"
)

const (
	DefaultBranch = "HEAD"
	// Same as Gitaly,
	//see https://gitlab.com/gitlab-org/gitaly/blob/2cb0d9f0604daabe63edc2c8271e65ef36ff6483/internal/git/repository.go#L16-16
	DefaultRef = "refs/heads/main"
	// Same as Gitaly,
	// see https://gitlab.com/gitlab-org/gitaly/blob/2cb0d9f0604daabe63edc2c8271e65ef36ff6483/internal/git/repository.go#L21-21
	LegacyDefaultRef = "refs/heads/master"
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

// Poll searched the given repository for the given fullRefName and returns a PollInfo containing a resolved Commit Object ID.
// Valid fullRefNames are:
// * `refs/heads/*` => for branches
// * `refs/tags/*` => for tags
// * `HEAD` => for the repository's default branch
func (p *Poller) Poll(ctx context.Context, repo *gitalypb.Repository, lastProcessedCommitId, fullRefName string) (*PollInfo, error) {
	noRefs := true
	var wanted, defaultRef, legacyDefaultRef *stats.Reference
	err := p.fetchRefs(ctx, repo, func(ref stats.Reference) bool {
		noRefs = false
		// We implement a similar logic here to what Gitaly does in their `GetDefaultBranch` logic
		// to find the default branch.
		// see https://gitlab.com/gitlab-org/gitaly/blob/2cb0d9f0604daabe63edc2c8271e65ef36ff6483/internal/git/localrepo/refs.go#L345-345
		switch string(ref.Name) {
		case fullRefName:
			wanted = cloneReference(ref)
			return true
		case DefaultRef:
			defaultRef = cloneReference(ref)
		case LegacyDefaultRef:
			legacyDefaultRef = cloneReference(ref)
		}
		return false
	})
	if err != nil {
		return nil, err // don't wrap
	}
	if noRefs {
		return &PollInfo{
			EmptyRepository: true,
		}, nil
	}

	if wanted == nil { // not found
		if fullRefName != DefaultBranch { // was looking for arbitrary branch, but didn't find it.
			return nil, NewNotFoundError("InfoRefsUploadPack", fullRefName)
		}

		// we have been searching for the default branch
		switch {
		case defaultRef != nil:
			wanted = defaultRef
		case legacyDefaultRef != nil:
			wanted = legacyDefaultRef
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
			if err == io.EOF { // nolint:errorlint
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
		Oid:  bytes.Clone(ref.Oid),
		Name: bytes.Clone(ref.Name),
	}
}
