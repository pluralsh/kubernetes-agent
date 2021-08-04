package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/kube_testing"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_internalgitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/mock_usage_metrics"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/agentcfg"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/pkg/kascfg"
	"gitlab.com/gitlab-org/gitaly/v14/proto/go/gitalypb"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	defaultGitOpsManifestPathGlob = "**/*.{yaml,yml,json}"

	projectId        = "some/project"
	revision         = "507ebc6de9bcac25628aa7afd52802a91a0685d8"
	manifestRevision = "7afd52802a91a0685d8507ebc6de9bcac25628aa"
)

var (
	_ rpc.GitopsServer = &server{}
)

func TestGlobToGitaly(t *testing.T) {
	tests := []struct {
		name              string
		glob              string
		expectedRepoPath  []byte
		expectedRecursive bool
	}{
		{
			name:              "full file name",
			glob:              "simple-path/manifest.yaml",
			expectedRepoPath:  []byte("simple-path"),
			expectedRecursive: false,
		},
		{
			name:              "empty",
			glob:              "",
			expectedRepoPath:  []byte{'.'},
			expectedRecursive: false,
		},
		{
			name:              "simple file1",
			glob:              "*.yaml",
			expectedRepoPath:  []byte{'.'},
			expectedRecursive: false,
		},
		{
			name:              "files in directory1",
			glob:              "bla/*.yaml",
			expectedRepoPath:  []byte("bla"),
			expectedRecursive: false,
		},
		{
			name:              "recursive files in directory1",
			glob:              "bla/**/*.yaml",
			expectedRepoPath:  []byte("bla"),
			expectedRecursive: true,
		},
		{
			name:              "all files1",
			glob:              "**/*.yaml",
			expectedRepoPath:  []byte{'.'},
			expectedRecursive: true,
		},
		{
			name:              "group1",
			glob:              "[a-z]*/*.yaml",
			expectedRepoPath:  []byte{'.'},
			expectedRecursive: true,
		},
		{
			name:              "group2",
			glob:              "?bla/*.yaml",
			expectedRepoPath:  []byte{'.'},
			expectedRecursive: true,
		},
		{
			name:              "group3",
			glob:              "bla/?aaa/*.yaml",
			expectedRepoPath:  []byte("bla"),
			expectedRecursive: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotRepoPath, gotRecursive := globToGitaly(tc.glob)
			assert.Equal(t, tc.expectedRepoPath, gotRepoPath)
			assert.Equal(t, tc.expectedRecursive, gotRecursive)
		})
	}
}

func TestGetObjectsToSynchronize_GetProjectInfo_Forbidden(t *testing.T) {
	s, ctrl, _, _ := setupServerBare(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
	server.EXPECT().
		Context().
		Return(mock_modserver.IncomingCtx(context.Background(), t, testhelpers.AgentkToken)).
		MinTimes(1)
	err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{ProjectId: projectId}, server)
	require.Error(t, err)
	assert.Equal(t, codes.PermissionDenied, status.Code(err))
}

func TestGetObjectsToSynchronize_GetProjectInfo_Unauthorized(t *testing.T) {
	s, ctrl, _, _ := setupServerBare(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
	server.EXPECT().
		Context().
		Return(mock_modserver.IncomingCtx(context.Background(), t, testhelpers.AgentkToken)).
		MinTimes(1)
	err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{ProjectId: projectId}, server)
	require.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestGetObjectsToSynchronize_GetProjectInfo_InternalServerError(t *testing.T) {
	s, ctrl, mockApi, _ := setupServerBare(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})
	mockApi.EXPECT().
		HandleProcessingError(gomock.Any(), gomock.Any(), testhelpers.AgentId, "GetProjectInfo()", matcher.ErrorEq("error kind: 0; status: 500"))
	server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
	server.EXPECT().
		Context().
		Return(mock_modserver.IncomingCtx(context.Background(), t, testhelpers.AgentkToken)).
		MinTimes(1)
	err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{ProjectId: projectId}, server)
	require.NoError(t, err) // no error here, it keep trying for any other errors
}

func TestGetObjectsToSynchronize_HappyPath(t *testing.T) {
	ctx, cancel, s, ctrl, gitalyPool, _ := setupServer(t)
	s.syncCount.(*mock_usage_metrics.MockCounter).EXPECT().Inc()

	objs := objectsYAML(t)
	projInfo := projectInfo()
	server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
	server.EXPECT().
		Context().
		Return(mock_modserver.IncomingCtx(ctx, t, testhelpers.AgentkToken)).
		MinTimes(1)
	gomock.InOrder(
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Header_{
					Header: &rpc.ObjectsToSynchronizeResponse_Header{
						CommitId:  revision,
						ProjectId: projInfo.ProjectId,
					},
				},
			})),
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Object_{
					Object: &rpc.ObjectsToSynchronizeResponse_Object{
						Source: "manifest.yaml",
						Data:   objs[:1],
					},
				},
			})),
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Object_{
					Object: &rpc.ObjectsToSynchronizeResponse_Object{
						Source: "manifest.yaml",
						Data:   objs[1:],
					},
				},
			})),
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Trailer_{
					Trailer: &rpc.ObjectsToSynchronizeResponse_Trailer{},
				},
			})).
			Do(func(resp *rpc.ObjectsToSynchronizeResponse) {
				cancel() // stop streaming call after the first response has been sent
			}),
	)
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	pf := mock_internalgitaly.NewMockPathFetcherInterface(ctrl)
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &projInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), "", gitaly.DefaultBranch).
			Return(&gitaly.PollInfo{
				UpdateAvailable: true,
				CommitId:        revision,
			}, nil),
		gitalyPool.EXPECT().
			PathFetcher(gomock.Any(), &projInfo.GitalyInfo).
			Return(pf, nil),
		pf.EXPECT().
			Visit(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), []byte(revision), []byte("."), true, gomock.Any()).
			Do(func(ctx context.Context, repo *gitalypb.Repository, revision, repoPath []byte, recursive bool, visitor gitaly.FetchVisitor) {
				download, maxSize, err := visitor.Entry(&gitalypb.TreeEntry{
					Path:      []byte("manifest.yaml"),
					Type:      gitalypb.TreeEntry_BLOB,
					CommitOid: manifestRevision,
				})
				require.NoError(t, err)
				assert.EqualValues(t, defaultGitopsMaxManifestFileSize, maxSize)
				assert.True(t, download)

				done, err := visitor.StreamChunk([]byte("manifest.yaml"), objs[:1])
				require.NoError(t, err)
				assert.False(t, done)
				done, err = visitor.StreamChunk([]byte("manifest.yaml"), objs[1:])
				require.NoError(t, err)
				assert.False(t, done)
			}),
	)
	err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectId,
		Paths: []*agentcfg.PathCF{
			{
				Glob: defaultGitOpsManifestPathGlob,
			},
		},
	}, server)
	require.NoError(t, err)
}

func TestGetObjectsToSynchronize_HappyPath_Glob(t *testing.T) {
	ctx, cancel, s, ctrl, gitalyPool, _ := setupServer(t)
	s.syncCount.(*mock_usage_metrics.MockCounter).EXPECT().Inc()

	objs := objectsYAML(t)
	projInfo := projectInfo()
	server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
	server.EXPECT().
		Context().
		Return(mock_modserver.IncomingCtx(ctx, t, testhelpers.AgentkToken)).
		MinTimes(1)
	gomock.InOrder(
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Header_{
					Header: &rpc.ObjectsToSynchronizeResponse_Header{
						CommitId:  revision,
						ProjectId: projInfo.ProjectId,
					},
				},
			})),
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Object_{
					Object: &rpc.ObjectsToSynchronizeResponse_Object{
						Source: "path/manifest.yaml",
						Data:   objs,
					},
				},
			})),
		server.EXPECT().
			Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
				Message: &rpc.ObjectsToSynchronizeResponse_Trailer_{
					Trailer: &rpc.ObjectsToSynchronizeResponse_Trailer{},
				},
			})).
			Do(func(resp *rpc.ObjectsToSynchronizeResponse) {
				cancel() // stop streaming call after the first response has been sent
			}),
	)
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	pf := mock_internalgitaly.NewMockPathFetcherInterface(ctrl)
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &projInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), "", gitaly.DefaultBranch).
			Return(&gitaly.PollInfo{
				UpdateAvailable: true,
				CommitId:        revision,
			}, nil),
		gitalyPool.EXPECT().
			PathFetcher(gomock.Any(), &projInfo.GitalyInfo).
			Return(pf, nil),
		pf.EXPECT().
			Visit(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), []byte(revision), []byte("path"), false, gomock.Any()).
			Do(func(ctx context.Context, repo *gitalypb.Repository, revision, repoPath []byte, recursive bool, visitor gitaly.FetchVisitor) {
				download, maxSize, err := visitor.Entry(&gitalypb.TreeEntry{
					Path:      []byte("path/manifest.yaml"),
					Type:      gitalypb.TreeEntry_BLOB,
					CommitOid: manifestRevision,
				})
				require.NoError(t, err)
				assert.EqualValues(t, defaultGitopsMaxManifestFileSize, maxSize)
				assert.True(t, download)

				done, err := visitor.StreamChunk([]byte("path/manifest.yaml"), objs)
				require.NoError(t, err)
				assert.False(t, done)
			}),
	)
	err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectId,
		Paths: []*agentcfg.PathCF{
			{
				Glob: "/path/*.yaml",
			},
		},
	}, server)
	require.NoError(t, err)
}

func TestGetObjectsToSynchronize_ResumeConnection(t *testing.T) {
	ctx, cancel, s, ctrl, gitalyPool, _ := setupServer(t)
	projInfo := projectInfo()
	server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
	server.EXPECT().
		Context().
		Return(mock_modserver.IncomingCtx(ctx, t, testhelpers.AgentkToken)).
		MinTimes(1)
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &projInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), revision, gitaly.DefaultBranch).
			DoAndReturn(func(ctx context.Context, repo *gitalypb.Repository, lastProcessedCommitId, refName string) (*gitaly.PollInfo, error) {
				cancel() // stop the test
				return &gitaly.PollInfo{
					UpdateAvailable: false,
					CommitId:        revision,
				}, nil
			}),
	)
	err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectId,
		CommitId:  revision,
	}, server)
	require.NoError(t, err)
}

func TestGetObjectsToSynchronize_UserErrors(t *testing.T) {
	pathFetcherErrs := []struct {
		errMsg string
		err    error
	}{
		{
			errMsg: "manifest file: FileNotFound: Bla: file/directory/ref not found: some/file",
			err:    gitaly.NewNotFoundError("Bla", "some/file"),
		},
		{
			errMsg: "manifest file: FileTooBig: Bla: file is too big: some/file",
			err:    gitaly.NewFileTooBigError(nil, "Bla", "some/file"),
		},
		{
			errMsg: "manifest file: UnexpectedTreeEntryType: Bla: file is not a usual file: some/file",
			err:    gitaly.NewUnexpectedTreeEntryTypeError("Bla", "some/file"),
		},
		{
			errMsg: "manifest file: path visited more than once: bla",
			err:    &gitaly.DuplicatePathFoundError{Path: "bla"},
		},
		{
			errMsg: "manifest file: glob *.yaml match failed: bad glob",
			err:    &gitaly.GlobMatchFailedError{Cause: errors.New("bad glob"), Glob: "*.yaml"},
		},
		{
			errMsg: "manifest file: maximum number of files limit reached: 10",
			err:    &gitaly.MaxNumberOfFilesError{MaxNumberOfFiles: 10},
		},
	}
	for _, tc := range pathFetcherErrs {
		t.Run(tc.errMsg, func(t *testing.T) {
			ctx, _, s, ctrl, gitalyPool, mockApi := setupServer(t)

			projInfo := projectInfo()
			server := mock_rpc.NewMockGitops_GetObjectsToSynchronizeServer(ctrl)
			server.EXPECT().
				Context().
				Return(mock_modserver.IncomingCtx(ctx, t, testhelpers.AgentkToken)).
				MinTimes(1)
			server.EXPECT().
				Send(matcher.ProtoEq(t, &rpc.ObjectsToSynchronizeResponse{
					Message: &rpc.ObjectsToSynchronizeResponse_Header_{
						Header: &rpc.ObjectsToSynchronizeResponse_Header{
							CommitId:  revision,
							ProjectId: projInfo.ProjectId,
						},
					},
				}))
			mockApi.EXPECT().
				HandleProcessingError(gomock.Any(), gomock.Any(), testhelpers.AgentId, "GitOps: failed to get objects to synchronize",
					matcher.ErrorEq(tc.errMsg),
				)
			p := mock_internalgitaly.NewMockPollerInterface(ctrl)
			pf := mock_internalgitaly.NewMockPathFetcherInterface(ctrl)
			gomock.InOrder(
				gitalyPool.EXPECT().
					Poller(gomock.Any(), &projInfo.GitalyInfo).
					Return(p, nil),
				p.EXPECT().
					Poll(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), "", gitaly.DefaultBranch).
					Return(&gitaly.PollInfo{
						UpdateAvailable: true,
						CommitId:        revision,
					}, nil),
				gitalyPool.EXPECT().
					PathFetcher(gomock.Any(), &projInfo.GitalyInfo).
					Return(pf, nil),
				pf.EXPECT().
					Visit(gomock.Any(), matcher.ProtoEq(nil, projInfo.Repository), []byte(revision), []byte("."), true, gomock.Any()).
					Return(tc.err),
			)
			err := s.GetObjectsToSynchronize(&rpc.ObjectsToSynchronizeRequest{
				ProjectId: projectId,
				Paths: []*agentcfg.PathCF{
					{
						Glob: defaultGitOpsManifestPathGlob,
					},
				},
			}, server)
			assert.EqualError(t, err, fmt.Sprintf("rpc error: code = FailedPrecondition desc = GitOps: failed to get objects to synchronize: %s", tc.errMsg))
		})
	}
}

func projectInfoRest() *gapi.ProjectInfoResponse {
	return &gapi.ProjectInfoResponse{
		ProjectId: 234,
		GitalyInfo: gitlab.GitalyInfo{
			Address: "127.0.0.1:321321",
			Token:   "cba",
			Features: map[string]string{
				"bla": "false",
			},
		},
		GitalyRepository: gitlab.GitalyRepository{
			StorageName:   "StorageName1",
			RelativePath:  "RelativePath1",
			GlRepository:  "GlRepository1",
			GlProjectPath: "GlProjectPath1",
		},
	}
}

func projectInfo() *api.ProjectInfo {
	rest := projectInfoRest()
	return &api.ProjectInfo{
		ProjectId:  rest.ProjectId,
		GitalyInfo: rest.GitalyInfo.ToGitalyInfo(),
		Repository: rest.GitalyRepository.ToProtoRepository(),
	}
}

func setupServer(t *testing.T) (context.Context, context.CancelFunc, *server, *gomock.Controller, *mock_internalgitaly.MockPoolInterface, *mock_modserver.MockAPI) {
	ctx, correlationId := testhelpers.CtxWithCorrelation(t)
	ctx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	s, ctrl, mockApi, gitalyPool := setupServerBare(t, func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertGetJsonRequestIsCorrect(t, r, correlationId)
		assert.Equal(t, projectId, r.URL.Query().Get(gapi.ProjectIdQueryParam))
		testhelpers.RespondWithJSON(t, w, projectInfoRest())
	})

	return ctx, cancel, s, ctrl, gitalyPool, mockApi
}

func setupServerBare(t *testing.T, handler func(http.ResponseWriter, *http.Request)) (*server, *gomock.Controller, *mock_modserver.MockAPI, *mock_internalgitaly.MockPoolInterface) {
	ctrl := gomock.NewController(t)
	gitalyPool := mock_internalgitaly.NewMockPoolInterface(ctrl)
	mockApi := mock_modserver.NewMockAPIWithMockPoller(ctrl, 1)
	agentInfo := testhelpers.AgentInfoObj()
	mockApi.EXPECT().
		GetAgentInfo(gomock.Any(), gomock.Any(), testhelpers.AgentkToken).
		Return(agentInfo, nil)
	usageTracker := mock_usage_metrics.NewMockUsageTrackerInterface(ctrl)
	usageTracker.EXPECT().
		RegisterCounter(gitopsSyncCountKnownMetric).
		Return(mock_usage_metrics.NewMockCounter(ctrl))

	config := &kascfg.ConfigurationFile{}
	ApplyDefaults(config)
	config.Agent.Gitops.ProjectInfoCacheTtl = durationpb.New(0)
	config.Agent.Gitops.ProjectInfoCacheErrorTtl = durationpb.New(0)
	s, err := newServerFromConfig(&modserver.Config{
		Log:          zaptest.NewLogger(t),
		Api:          mockApi,
		Config:       config,
		GitLabClient: mock_gitlab.SetupClient(t, gapi.ProjectInfoApiPath, handler),
		Registerer:   prometheus.NewPedanticRegistry(),
		UsageTracker: usageTracker,
		Gitaly:       gitalyPool,
	})
	require.NoError(t, err)
	return s, ctrl, mockApi, gitalyPool
}

func objectsYAML(t *testing.T) []byte {
	objects := []runtime.Object{
		&corev1.ConfigMap{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ConfigMap",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "map1",
			},
			Data: map[string]string{
				"key": "value",
			},
		},
		&corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "ns1",
			},
		},
	}
	return kube_testing.ObjsToYAML(t, objects...)
}