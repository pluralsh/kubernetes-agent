package server

import (
	"context"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v2"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/api/apiutil"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/gitops"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/errz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/logz"
	"gitlab.com/gitlab-org/gitaly/proto/go/gitalypb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	gitOpsManifestMaxChunkSize = 128 * 1024
)

var (
	// globPrefix captures glob prefix that does not contain any special characters, recognized by doublestar.Match.
	// See https://github.com/bmatcuk/doublestar#about and
	// https://pkg.go.dev/github.com/bmatcuk/doublestar/v2#Match for globbing rules.
	globPrefix = regexp.MustCompile(`^/?([^\\*?[\]{}]+)/(.*)$`)
)

type module struct {
	log                            *zap.Logger
	api                            modserver.API
	gitalyPool                     gitaly.PoolInterface
	projectInfoClient              *projectInfoClient
	gitopsSyncCount                usage_metrics.Counter
	gitopsPollPeriod               time.Duration
	connectionMaxAge               time.Duration
	maxGitopsManifestFileSize      int64
	maxGitopsTotalManifestFileSize int64
	maxGitopsNumberOfPaths         uint32
	maxGitopsNumberOfFiles         uint32
}

func (m *module) Run(ctx context.Context) error {
	return nil
}

func (m *module) GetObjectsToSynchronize(req *rpc.ObjectsToSynchronizeRequest, server rpc.Gitops_GetObjectsToSynchronizeServer) error {
	ctx := server.Context()
	agentMeta := apiutil.AgentMetaFromContext(ctx)
	l := m.log.With(logz.CorrelationIdFromContext(ctx))
	agentInfo, err, retErr := m.api.GetAgentInfo(ctx, l, agentMeta, false)
	if retErr {
		return err
	}
	err = m.validateGetObjectsToSynchronizeRequest(req)
	if err != nil {
		return err // no wrap
	}
	return m.api.PollImmediateUntil(ctx, m.gitopsPollPeriod, m.connectionMaxAge, m.sendObjectsToSynchronize(agentInfo, req, server))
}

func (m *module) Name() string {
	return gitops.ModuleName
}

func (m *module) validateGetObjectsToSynchronizeRequest(req *rpc.ObjectsToSynchronizeRequest) error {
	numberOfPaths := uint32(len(req.Paths))
	if numberOfPaths > m.maxGitopsNumberOfPaths {
		// TODO validate config in GetConfiguration too and send it somewhere the user can see it https://gitlab.com/gitlab-org/gitlab/-/issues/277323
		// This check must be here, but there too.
		return status.Errorf(codes.InvalidArgument, "maximum number of GitOps paths per manifest project is %d, but %d was requested", m.maxGitopsNumberOfPaths, numberOfPaths)
	}
	return nil
}

func (m *module) sendObjectsToSynchronize(agentInfo *api.AgentInfo, req *rpc.ObjectsToSynchronizeRequest, server rpc.Gitops_GetObjectsToSynchronizeServer) modserver.ConditionFunc {
	ctx := server.Context()
	l := m.log.With(logz.AgentId(agentInfo.Id), logz.ProjectId(req.ProjectId), logz.CorrelationIdFromContext(ctx))
	return func() (bool /*done*/, error) {
		// This call is made on each poll because:
		// - it checks that the agent's token is still valid
		// - repository location in Gitaly might have changed
		projectInfo, err, retErr := m.getProjectInfo(ctx, l, &agentInfo.Meta, req.ProjectId)
		if retErr {
			return false, err
		}
		revision := gitaly.DefaultBranch // TODO support user-specified branches/tags
		p, err := m.gitalyPool.Poller(ctx, &projectInfo.GitalyInfo)
		if err != nil {
			m.api.HandleProcessingError(ctx, l, "GitOps: Poller", err)
			return false, nil // don't want to close the response stream, so report no error
		}
		info, err := p.Poll(ctx, &projectInfo.Repository, req.CommitId, revision)
		if err != nil {
			m.api.HandleProcessingError(ctx, l, "GitOps: repository poll failed", err)
			return false, nil // don't want to close the response stream, so report no error
		}
		if !info.UpdateAvailable {
			l.Debug("GitOps: no updates", logz.CommitId(req.CommitId))
			return false, nil
		}
		// Create a new l variable, don't want to mutate the one from the outer scope
		l := l.With(logz.CommitId(info.CommitId)) // nolint:govet
		l.Info("GitOps: new commit")
		err = m.sendObjectsToSynchronizeHeaders(server, l, info.CommitId)
		if err != nil {
			return false, err // no wrap
		}
		numberOfFiles, err := m.sendObjectsToSynchronizeBody(req, server, l, &projectInfo.Repository, &projectInfo.GitalyInfo, info.CommitId)
		if err != nil {
			return false, err // no wrap
		}
		err = m.sendObjectsToSynchronizeTrailers(server, l)
		if err != nil {
			return false, err // no wrap
		}
		l.Info("GitOps: fetched files", logz.NumberOfFiles(numberOfFiles))
		m.gitopsSyncCount.Inc()
		return true, nil
	}
}

func (m *module) sendObjectsToSynchronizeHeaders(server rpc.Gitops_GetObjectsToSynchronizeServer, log *zap.Logger, commitId string) error {
	err := server.Send(&rpc.ObjectsToSynchronizeResponse{
		Message: &rpc.ObjectsToSynchronizeResponse_Headers_{
			Headers: &rpc.ObjectsToSynchronizeResponse_Headers{
				CommitId: commitId,
			},
		},
	})
	if err != nil {
		return m.api.HandleSendError(log, "GitOps: failed to send headers for objects to synchronize", err)
	}
	return nil
}

func (m *module) sendObjectsToSynchronizeBody(req *rpc.ObjectsToSynchronizeRequest, server rpc.Gitops_GetObjectsToSynchronizeServer, log *zap.Logger, repo *gitalypb.Repository, gitalyInfo *api.GitalyInfo, commitId string) (uint32, error) {
	ctx := server.Context()
	pf, err := m.gitalyPool.PathFetcher(ctx, gitalyInfo)
	if err != nil {
		m.api.HandleProcessingError(ctx, log, "GitOps: PathFetcher", err)
		return 0, status.Error(codes.Unavailable, "GitOps: PathFetcher")
	}
	v := &objectsToSynchronizeVisitor{
		server:                 server,
		remainingTotalFileSize: m.maxGitopsTotalManifestFileSize,
		fileSizeLimit:          m.maxGitopsManifestFileSize,
		maxNumberOfFiles:       m.maxGitopsNumberOfFiles,
	}
	vChunk := gitaly.ChunkingFetchVisitor{
		MaxChunkSize: gitOpsManifestMaxChunkSize,
		Delegate:     v,
	}
	for _, p := range req.Paths {
		repoPath, recursive, glob := globToGitaly(p.Glob)
		v.glob = glob // set new glob for each path
		err = pf.Visit(ctx, repo, []byte(commitId), repoPath, recursive, vChunk)
		if err != nil {
			if v.sendFailed {
				return 0, m.api.HandleSendError(log, "GitOps: failed to send objects to synchronize", err)
			} else {
				m.api.HandleProcessingError(ctx, log, "GitOps: failed to get objects to synchronize", err)
				return 0, status.Error(codes.Unavailable, "GitOps: failed to get objects to synchronize")
			}
		}
	}
	return v.numberOfFiles, nil
}

func (m *module) sendObjectsToSynchronizeTrailers(server rpc.Gitops_GetObjectsToSynchronizeServer, log *zap.Logger) error {
	err := server.Send(&rpc.ObjectsToSynchronizeResponse{
		Message: &rpc.ObjectsToSynchronizeResponse_Trailers_{
			Trailers: &rpc.ObjectsToSynchronizeResponse_Trailers{},
		},
	})
	if err != nil {
		return m.api.HandleSendError(log, "GitOps: failed to send trailers for objects to synchronize", err)
	}
	return nil
}

func (m *module) getProjectInfo(ctx context.Context, log *zap.Logger, agentMeta *api.AgentMeta, projectId string) (*api.ProjectInfo, error, bool /* return the error? */) {
	projectInfo, err := m.projectInfoClient.GetProjectInfo(ctx, agentMeta, projectId)
	switch {
	case err == nil:
		return projectInfo, nil, false
	case errz.ContextDone(err):
		err = status.Error(codes.Unavailable, "unavailable")
	case gitlab.IsForbidden(err):
		err = status.Error(codes.PermissionDenied, "forbidden")
	case gitlab.IsUnauthorized(err):
		err = status.Error(codes.Unauthenticated, "unauthenticated")
	default:
		m.api.LogAndCapture(ctx, log, "GetProjectInfo()", err)
		err = nil // don't want to close the response stream, so report no error
	}
	return nil, err, true
}

type objectsToSynchronizeVisitor struct {
	server                 rpc.Gitops_GetObjectsToSynchronizeServer
	glob                   string
	remainingTotalFileSize int64
	fileSizeLimit          int64
	maxNumberOfFiles       uint32
	numberOfFiles          uint32
	sendFailed             bool
}

func (v *objectsToSynchronizeVisitor) Entry(entry *gitalypb.TreeEntry) (bool /* download? */, int64 /* max size */, error) {
	if v.numberOfFiles == v.maxNumberOfFiles {
		return false, 0, errz.NewUserErrorf("maximum number of manifest files limit reached: %d", v.maxNumberOfFiles)
	}
	v.numberOfFiles++
	filename := string(entry.Path)
	if isHiddenDir(filename) {
		return false, 0, nil
	}
	shouldDownload, err := doublestar.Match(v.glob, filename)
	if err != nil {
		return false, 0, errz.NewUserErrorWithCausef(err, "glob %s match failed", v.glob)
	}
	return shouldDownload, minInt64(v.remainingTotalFileSize, v.fileSizeLimit), nil
}

func (v *objectsToSynchronizeVisitor) StreamChunk(path []byte, data []byte) (bool /* done? */, error) {
	v.remainingTotalFileSize -= int64(len(data))
	if v.remainingTotalFileSize < 0 {
		// This should never happen because we told Gitaly the maximum file size that we'd like to get.
		// i.e. we should have gotten an error from Gitaly if file is bigger than the limit.
		return false, status.Error(codes.Internal, "unexpected negative remaining total file size")
	}
	err := v.server.Send(&rpc.ObjectsToSynchronizeResponse{
		Message: &rpc.ObjectsToSynchronizeResponse_Object_{
			Object: &rpc.ObjectsToSynchronizeResponse_Object{
				Source: string(path),
				Data:   data,
			},
		},
	})
	if err != nil {
		v.sendFailed = true
	}
	return false, err
}

// isHiddenDir checks if a file is in a directory, which name starts with a dot.
func isHiddenDir(filename string) bool {
	dir := path.Dir(filename)
	if dir == "." { // root directory special case
		return false
	}
	parts := strings.Split(dir, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

func globToGitaly(glob string) ([]byte /* repoPath */, bool /* recursive */, string /* glob */) {
	var repoPath []byte
	matches := globPrefix.FindStringSubmatch(glob)
	if matches == nil {
		repoPath = []byte{'.'}
		glob = strings.TrimPrefix(glob, "/") // remove at most one slash to match regex
	} else {
		repoPath = []byte(matches[1])
		glob = matches[2]
	}
	recursive := strings.ContainsAny(glob, "[/") || // cannot determine if recursive or not because character class may contain ranges, etc
		strings.Contains(glob, "**") // contains directory match
	return repoPath, recursive, glob
}