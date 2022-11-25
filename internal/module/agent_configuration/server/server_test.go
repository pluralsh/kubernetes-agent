package server

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitaly"
	gapi "gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/gitlab/api"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_configuration/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/module/modshared"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_agent_tracker"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_internalgitaly"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modserver"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/pkg/agentcfg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
	"sigs.k8s.io/yaml"
)

const (
	projectId    = "some/project"
	revision     = "507ebc6de9bcac25628aa7afd52802a91a0685d8"
	branchPrefix = "refs/heads/"

	maxConfigurationFileSize = 128 * 1024
)

var (
	_ modserver.Module             = (*module)(nil)
	_ modserver.Factory            = (*Factory)(nil)
	_ modserver.ApplyDefaults      = ApplyDefaults
	_ rpc.AgentConfigurationServer = (*server)(nil)
)

func TestEmptyConfig(t *testing.T) {
	t.Run("comments", func(t *testing.T) {
		data := []byte(`
#gitops:
#  manifest_projects:
#  - id: "root/gitops-manifests"
#    paths:
#      - glob: "/bla/**"
`)
		assertEmpty(t, data)
	})
	t.Run("empty", func(t *testing.T) {
		data := []byte("")
		assertEmpty(t, data)
	})
	t.Run("newline", func(t *testing.T) {
		data := []byte("\n")
		assertEmpty(t, data)
	})
	t.Run("missing", func(t *testing.T) {
		var data []byte
		assertEmpty(t, data)
	})
}

func assertEmpty(t *testing.T, data []byte) {
	config, err := parseYAMLToConfiguration(data)
	require.NoError(t, err)
	diff := cmp.Diff(config, &agentcfg.ConfigurationFile{}, protocmp.Transform())
	assert.Empty(t, diff)
}

func TestYAMLToConfigurationAndBack(t *testing.T) {
	testCases := []struct {
		given, expected string
	}{
		{
			given: `{}
`, // empty config
			expected: `{}
`,
		},
		{
			given: `gitops: {}
`,
			expected: `gitops: {}
`,
		},
		{
			given: `gitops:
  manifest_projects: []
`,
			expected: `gitops: {}
`, // empty slice is omitted
		},
		{
			expected: `gitops:
  manifest_projects:
  - id: gitlab-org/cluster-integration/gitlab-agent
`,
			given: `gitops:
  manifest_projects:
  - id: gitlab-org/cluster-integration/gitlab-agent
`,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			config, err := parseYAMLToConfiguration([]byte(tc.given))
			require.NoError(t, err)
			configJson, err := protojson.Marshal(config)
			require.NoError(t, err)
			configYaml, err := yaml.JSONToYAML(configJson)
			require.NoError(t, err)
			diff := cmp.Diff(tc.expected, string(configYaml))
			assert.Empty(t, diff)
		})
	}
}

func TestGetConfiguration_HappyPath(t *testing.T) {
	s, agentInfo, ctrl, gitalyPool, resp, _ := setupServer(t)
	configFile := sampleConfig()
	resp.EXPECT().
		Send(matcher.ProtoEq(t, &rpc.ConfigurationResponse{
			Configuration: &agentcfg.AgentConfiguration{
				Gitops: &agentcfg.GitopsCF{
					ManifestProjects: []*agentcfg.ManifestProjectCF{
						{
							Id: projectId,
						},
					},
				},
				AgentId:   agentInfo.Id,
				ProjectId: agentInfo.ProjectId,
			},
			CommitId: revision,
		}))
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	pf := mock_internalgitaly.NewMockPathFetcherInterface(ctrl)
	configFileName := agent_configuration.Directory + "/" + agentInfo.Name + "/" + agent_configuration.FileName
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &agentInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), "", branchPrefix+agentInfo.DefaultBranch).
			Return(&gitaly.PollInfo{
				CommitId:        revision,
				UpdateAvailable: true,
			}, nil),
		gitalyPool.EXPECT().
			PathFetcher(gomock.Any(), &agentInfo.GitalyInfo).
			Return(pf, nil),
		pf.EXPECT().
			FetchFile(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), []byte(revision), []byte(configFileName), int64(maxConfigurationFileSize)).
			Return(configToBytes(t, configFile), nil),
	)
	err := s.GetConfiguration(&rpc.ConfigurationRequest{
		AgentMeta: agentMeta(),
	}, resp)
	require.NoError(t, err)
}

func TestGetConfiguration_ResumeConnection(t *testing.T) {
	s, agentInfo, ctrl, gitalyPool, resp, _ := setupServer(t)
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &agentInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), revision, branchPrefix+agentInfo.DefaultBranch).
			Return(&gitaly.PollInfo{
				CommitId:        revision,
				UpdateAvailable: false,
			}, nil),
	)
	err := s.GetConfiguration(&rpc.ConfigurationRequest{
		CommitId:  revision, // same commit id
		AgentMeta: agentMeta(),
	}, resp)
	require.NoError(t, err)
}

func TestGetConfiguration_ConfigNotFound(t *testing.T) {
	s, agentInfo, ctrl, gitalyPool, resp, _ := setupServer(t)
	resp.EXPECT().
		Send(matcher.ProtoEq(t, &rpc.ConfigurationResponse{
			Configuration: &agentcfg.AgentConfiguration{
				AgentId:   agentInfo.Id,
				ProjectId: agentInfo.ProjectId,
			},
			CommitId: revision,
		}))
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	pf := mock_internalgitaly.NewMockPathFetcherInterface(ctrl)
	configFileName := agent_configuration.Directory + "/" + agentInfo.Name + "/" + agent_configuration.FileName
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &agentInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), "", branchPrefix+agentInfo.DefaultBranch).
			Return(&gitaly.PollInfo{
				CommitId:        revision,
				UpdateAvailable: true,
			}, nil),
		gitalyPool.EXPECT().
			PathFetcher(gomock.Any(), &agentInfo.GitalyInfo).
			Return(pf, nil),
		pf.EXPECT().
			FetchFile(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), []byte(revision), []byte(configFileName), int64(maxConfigurationFileSize)).
			Return(nil, gitaly.NewNotFoundError("Bla", "some/file")),
	)
	err := s.GetConfiguration(&rpc.ConfigurationRequest{
		AgentMeta: agentMeta(),
	}, resp)
	require.NoError(t, err)
}

func TestGetConfiguration_EmptyRepository(t *testing.T) {
	s, agentInfo, ctrl, gitalyPool, resp, _ := setupServer(t)
	p := mock_internalgitaly.NewMockPollerInterface(ctrl)
	gomock.InOrder(
		gitalyPool.EXPECT().
			Poller(gomock.Any(), &agentInfo.GitalyInfo).
			Return(p, nil),
		p.EXPECT().
			Poll(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), "", branchPrefix+agentInfo.DefaultBranch).
			Return(&gitaly.PollInfo{
				EmptyRepository: true,
			}, nil),
	)
	err := s.GetConfiguration(&rpc.ConfigurationRequest{
		AgentMeta: agentMeta(),
	}, resp)
	require.NoError(t, err)
}

func TestGetConfiguration_UserErrors(t *testing.T) {
	gitalyErrs := []error{
		gitaly.NewFileTooBigError(nil, "Bla", "some/file"),
		gitaly.NewUnexpectedTreeEntryTypeError("Bla", "some/file"),
	}
	for _, gitalyErr := range gitalyErrs {
		t.Run(gitalyErr.(*gitaly.Error).Code.String(), func(t *testing.T) { // nolint: errorlint
			s, agentInfo, ctrl, gitalyPool, resp, mockRpcApi := setupServer(t)
			p := mock_internalgitaly.NewMockPollerInterface(ctrl)
			pf := mock_internalgitaly.NewMockPathFetcherInterface(ctrl)
			configFileName := agent_configuration.Directory + "/" + agentInfo.Name + "/" + agent_configuration.FileName
			gomock.InOrder(
				gitalyPool.EXPECT().
					Poller(gomock.Any(), &agentInfo.GitalyInfo).
					Return(p, nil),
				p.EXPECT().
					Poll(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), "", branchPrefix+agentInfo.DefaultBranch).
					Return(&gitaly.PollInfo{
						CommitId:        revision,
						UpdateAvailable: true,
					}, nil),
				gitalyPool.EXPECT().
					PathFetcher(gomock.Any(), &agentInfo.GitalyInfo).
					Return(pf, nil),
				pf.EXPECT().
					FetchFile(gomock.Any(), matcher.ProtoEq(nil, agentInfo.Repository), []byte(revision), []byte(configFileName), int64(maxConfigurationFileSize)).
					Return(nil, gitalyErr),
				mockRpcApi.EXPECT().
					HandleProcessingError(gomock.Any(), testhelpers.AgentId, "Config: failed to fetch",
						matcher.ErrorEq(fmt.Sprintf("agent configuration file: %v", gitalyErr))),
			)
			err := s.GetConfiguration(&rpc.ConfigurationRequest{
				AgentMeta: agentMeta(),
			}, resp)
			assert.EqualError(t, err, fmt.Sprintf("rpc error: code = FailedPrecondition desc = Config: agent configuration file: %v", gitalyErr))
		})
	}
}

func TestGetConfiguration_GetAgentInfo_Error(t *testing.T) {
	s, _, _, resp, mockRpcApi, _ := setupServerBare(t, 1)
	mockRpcApi.EXPECT().
		AgentInfo(gomock.Any(), gomock.Any()).
		Return(nil, status.Error(codes.PermissionDenied, "expected err")) // code doesn't matter, we test that we return on error
	err := s.GetConfiguration(&rpc.ConfigurationRequest{
		AgentMeta: agentMeta(),
	}, resp)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = expected err")
}

func TestGetConfiguration_GetAgentInfo_RetriableError(t *testing.T) {
	s, _, _, resp, mockRpcApi, _ := setupServerBare(t, 2)
	gomock.InOrder(
		mockRpcApi.EXPECT().
			AgentInfo(gomock.Any(), gomock.Any()).
			Return(nil, status.Error(codes.Unavailable, "unavailable")),
		mockRpcApi.EXPECT().
			AgentInfo(gomock.Any(), gomock.Any()).
			Return(nil, status.Error(codes.PermissionDenied, "expected err")), // code doesn't matter, we test that we return on error
	)
	err := s.GetConfiguration(&rpc.ConfigurationRequest{
		AgentMeta: agentMeta(),
	}, resp)
	assert.EqualError(t, err, "rpc error: code = PermissionDenied desc = expected err")
}

func TestFetchContainerScanningConfiguration(t *testing.T) {
	s := &server{}
	starboardConfig := &agentcfg.StarboardCF{Cadence: "0 * * * *"}
	containerScanningConfig := &agentcfg.StarboardCF{Cadence: "30 * * * *"}
	t.Run("ContainerScanning config missing, Starboard config present", func(t *testing.T) {
		result := s.fetchContainerScanningConfiguration(&agentcfg.ConfigurationFile{Starboard: starboardConfig})
		assert.Equal(t, starboardConfig, result)
	})
	t.Run("ContainerScanning config present, Starboard config missing", func(t *testing.T) {
		result := s.fetchContainerScanningConfiguration(&agentcfg.ConfigurationFile{ContainerScanning: containerScanningConfig})
		assert.Equal(t, containerScanningConfig, result)
	})
	t.Run("ContainerScanning config present, Starboard config present", func(t *testing.T) {
		result := s.fetchContainerScanningConfiguration(&agentcfg.ConfigurationFile{ContainerScanning: containerScanningConfig, Starboard: starboardConfig})
		assert.Equal(t, containerScanningConfig, result)
	})
}

func setupServerBare(t *testing.T, pollTimes int) (*server, *gomock.Controller, *mock_internalgitaly.MockPoolInterface, *mock_rpc.MockAgentConfiguration_GetConfigurationServer, *mock_modserver.MockAgentRpcApi, *mock_agent_tracker.MockTracker) {
	ctrl := gomock.NewController(t)
	mockRpcApi := mock_modserver.NewMockAgentRpcApiWithMockPoller(ctrl, pollTimes)
	gitalyPool := mock_internalgitaly.NewMockPoolInterface(ctrl)
	agentTracker := mock_agent_tracker.NewMockTracker(ctrl)
	gitLabClient := mock_gitlab.SetupClient(t, gapi.AgentConfigurationApiPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	s := &server{
		agentRegisterer:            agentTracker,
		gitaly:                     gitalyPool,
		gitLabClient:               gitLabClient,
		getConfigurationPollConfig: testhelpers.NewPollConfig(10 * time.Minute),
		maxConfigurationFileSize:   maxConfigurationFileSize,
	}
	resp := mock_rpc.NewMockAgentConfiguration_GetConfigurationServer(ctrl)
	resp.EXPECT().
		Context().
		Return(mock_modserver.IncomingAgentCtx(t, mockRpcApi)).
		MinTimes(1)
	return s, ctrl, gitalyPool, resp, mockRpcApi, agentTracker
}

func setupServer(t *testing.T) (*server, *api.AgentInfo, *gomock.Controller, *mock_internalgitaly.MockPoolInterface, *mock_rpc.MockAgentConfiguration_GetConfigurationServer, *mock_modserver.MockAgentRpcApi) {
	s, ctrl, gitalyPool, resp, mockRpcApi, agentTracker := setupServerBare(t, 1)
	agentInfo := testhelpers.AgentInfoObj()
	connMatcher := matcher.ProtoEq(t, &agent_tracker.ConnectedAgentInfo{
		AgentMeta: agentMeta(),
		AgentId:   agentInfo.Id,
		ProjectId: agentInfo.ProjectId,
	}, protocmp.IgnoreFields(&agent_tracker.ConnectedAgentInfo{}, "connected_at", "connection_id"))
	gomock.InOrder(
		mockRpcApi.EXPECT().
			AgentInfo(gomock.Any(), gomock.Any()).
			Return(agentInfo, nil),
		agentTracker.EXPECT().
			RegisterConnection(gomock.Any(), connMatcher),
	)
	agentTracker.EXPECT().
		UnregisterConnection(gomock.Any(), connMatcher)
	return s, agentInfo, ctrl, gitalyPool, resp, mockRpcApi
}

func configToBytes(t *testing.T, configFile *agentcfg.ConfigurationFile) []byte {
	configJson, err := protojson.Marshal(configFile)
	require.NoError(t, err)
	configYaml, err := yaml.JSONToYAML(configJson)
	require.NoError(t, err)
	return configYaml
}

func sampleConfig() *agentcfg.ConfigurationFile {
	return &agentcfg.ConfigurationFile{
		Gitops: &agentcfg.GitopsCF{
			ManifestProjects: []*agentcfg.ManifestProjectCF{
				{
					Id: projectId,
				},
			},
		},
	}
}

func agentMeta() *modshared.AgentMeta {
	return &modshared.AgentMeta{
		Version:      "v1.2.3",
		CommitId:     "32452345",
		PodNamespace: "ns1",
		PodName:      "n1",
	}
}
