package chartops

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	revision    = "rev12341234"
	releaseName = "rel1"
)

var (
	defaultNamespace        = defaultChartNamespace
	maxHistory       int32  = defaultChartMaxHistory
	projectPath             = "proj1"
	maxFileSize      uint32 = defaultUrlValueMaxFileSize
)

var (
	_ Helm                                      = (*HelmActions)(nil)
	_ modagent.Factory                          = (*Factory)(nil)
	_ modagent.Worker                           = (*worker)(nil)
	_ modagent.WorkerFactory[*agentcfg.ChartCF] = (*workerFactory)(nil)
	_ modagent.WorkSource[*agentcfg.ChartCF]    = (*manifestSource)(nil)
)

func TestRun_HappyPath_ChartInRoot(t *testing.T) {
	inline, err := structpb.NewStruct(map[string]any{
		"c": "d",
	})
	require.NoError(t, err)
	w, helm, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id: projectPath,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_File{
					File: &agentcfg.ChartValuesFileCF{
						ProjectId: &projectPath,
						File:      "prod.yaml",
					},
				},
			},
			{
				From: &agentcfg.ChartValuesCF_Inline{
					Inline: inline,
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "/**",
				},
			},
			{
				Path: &rpc.PathCF_File{
					File: "prod.yaml",
				},
			},
		},
	}
	gomock.InOrder(
		watcher.EXPECT().
			Watch(gomock.Any(), matcher.ProtoEq(nil, req), gomock.Any()).
			Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
				callback(ctx, rpc.ObjectsToSynchronizeData{
					CommitId:  revision,
					ProjectId: testhelpers.ProjectId,
					Sources: []rpc.ObjectSource{
						{
							Name: "Chart.yaml",
							Data: []byte(`apiVersion: v2
name: test1
version: 1`),
						},
						{
							Name: "values.yaml",
							Data: []byte(`x: z`),
						},
						{
							Name: "prod.yaml",
							Data: []byte(`a: b
c: x`),
						},
					},
				})
				<-ctx.Done()
			}),
		helm.EXPECT().
			History(releaseName).
			Return(nil, driver.ErrReleaseNotFound),
		helm.EXPECT().
			Install(gomock.Any(), gomock.Any(), gomock.Any(), InstallConfig{
				Namespace:   defaultNamespace,
				ReleaseName: releaseName,
			}).
			Do(func(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) {
				assert.Equal(t, ChartValues{
					"a": "b", // from prod.yaml
					"c": "d", // from inline, overrides prod.yaml
				}, vals)
				assert.Equal(t, map[string]any{
					"x": "z", // from values.yaml
				}, chart.Values)
				cancel()
			}),
	)
	w.Run(ctx)
}

func TestRun_HappyPath_ValuesNotWithChart(t *testing.T) {
	w, helm, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id:   projectPath,
					Path: "chart",
					//Ref:  nil,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_File{
					File: &agentcfg.ChartValuesFileCF{
						ProjectId: &projectPath,
						File:      "prod/prod.yaml",
					},
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "chart/**",
				},
			},
			{
				Path: &rpc.PathCF_File{
					File: "prod/prod.yaml",
				},
			},
		},
	}
	gomock.InOrder(
		watcher.EXPECT().
			Watch(gomock.Any(), matcher.ProtoEq(nil, req), gomock.Any()).
			Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
				callback(ctx, rpc.ObjectsToSynchronizeData{
					CommitId:  revision,
					ProjectId: testhelpers.ProjectId,
					Sources: []rpc.ObjectSource{
						{
							Name: "chart/Chart.yaml",
							Data: []byte(`apiVersion: v2
name: test1
version: 1`),
						},
						{
							Name: "chart/values.yaml",
							Data: []byte(`x: z`),
						},
						{
							Name: "prod/prod.yaml",
							Data: []byte(`a: b
c: x`),
						},
					},
				})
				<-ctx.Done()
			}),
		helm.EXPECT().
			History(releaseName).
			Return(nil, driver.ErrReleaseNotFound),
		helm.EXPECT().
			Install(gomock.Any(), gomock.Any(), gomock.Any(), InstallConfig{
				Namespace:   defaultNamespace,
				ReleaseName: releaseName,
			}).
			Do(func(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) {
				assert.Equal(t, ChartValues{
					"a": "b", // from prod.yaml
					"c": "x", // from prod.yaml
				}, vals)
				assert.Equal(t, map[string]any{
					"x": "z", // from values.yaml
				}, chart.Values)
				cancel()
			}),
	)
	w.Run(ctx)
}

func TestRun_HappyPath_NoChartOrValueChanges(t *testing.T) {
	w, helm, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id: projectPath,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_File{
					File: &agentcfg.ChartValuesFileCF{
						ProjectId: &projectPath,
						File:      "prod.yaml",
					},
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "/**",
				},
			},
			{
				Path: &rpc.PathCF_File{
					File: "prod.yaml",
				},
			},
		},
	}
	installed := make(chan struct{})
	gomock.InOrder(
		watcher.EXPECT().
			Watch(gomock.Any(), matcher.ProtoEq(nil, req), gomock.Any()).
			Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
				data := rpc.ObjectsToSynchronizeData{
					CommitId:  revision,
					ProjectId: testhelpers.ProjectId,
					Sources: []rpc.ObjectSource{
						{
							Name: "Chart.yaml",
							Data: []byte(`apiVersion: v2
name: test1
version: 1`),
						},
						{
							Name: "prod.yaml",
							Data: []byte(`a: b
c: x`),
						},
					},
				}
				callback(ctx, data)
				<-installed
				callback(ctx, data) // same things
				<-ctx.Done()
			}),
		helm.EXPECT().
			History(releaseName).
			Return(nil, driver.ErrReleaseNotFound),
		helm.EXPECT().
			Install(gomock.Any(), gomock.Any(), gomock.Any(), InstallConfig{
				Namespace:   defaultNamespace,
				ReleaseName: releaseName,
			}).
			Do(func(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) {
				assert.Equal(t, ChartValues{
					"a": "b", // from prod.yaml
					"c": "x", // from prod.yaml
				}, vals)
				close(installed)
			}),
	)
	w.Run(ctx)
}

func TestRun_HappyPath_ValuesFromAnotherProject(t *testing.T) {
	project2Path := "proj2"
	w, helm, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id: projectPath,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_File{
					File: &agentcfg.ChartValuesFileCF{
						ProjectId: &project2Path,
						Ref:       &agentcfg.GitRefCF{Ref: &agentcfg.GitRefCF_Branch{Branch: "prod"}},
						File:      "prod.yaml",
					},
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req1 := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "/**",
				},
			},
		},
	}
	req2 := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: project2Path,
		Ref:       &rpc.GitRefCF{Ref: &rpc.GitRefCF_Branch{Branch: "prod"}},
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_File{
					File: "prod.yaml",
				},
			},
		},
	}
	watcher.EXPECT().
		Watch(gomock.Any(), matcher.ProtoEq(nil, req1), gomock.Any()).
		Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
			callback(ctx, rpc.ObjectsToSynchronizeData{
				CommitId:  revision,
				ProjectId: testhelpers.ProjectId,
				Sources: []rpc.ObjectSource{
					{
						Name: "Chart.yaml",
						Data: []byte(`apiVersion: v2
name: test1
version: 1`),
					},
					{
						Name: "values.yaml",
						Data: []byte(`x: z`),
					},
				},
			})
			<-ctx.Done()
		})
	watcher.EXPECT().
		Watch(gomock.Any(), matcher.ProtoEq(nil, req2), gomock.Any()).
		Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
			callback(ctx, rpc.ObjectsToSynchronizeData{
				CommitId:  revision,
				ProjectId: testhelpers.ProjectId,
				Sources: []rpc.ObjectSource{
					{
						Name: "prod.yaml",
						Data: []byte(`a: b
c: x`),
					},
				},
			})
			<-ctx.Done()
		})

	gomock.InOrder(
		helm.EXPECT().
			History(releaseName).
			Return(nil, driver.ErrReleaseNotFound),
		helm.EXPECT().
			Install(gomock.Any(), gomock.Any(), gomock.Any(), InstallConfig{
				Namespace:   defaultNamespace,
				ReleaseName: releaseName,
			}).
			Do(func(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) {
				assert.Equal(t, ChartValues{
					"a": "b", // from prod.yaml
					"c": "x", // from prod.yaml
				}, vals)
				assert.Equal(t, map[string]any{
					"x": "z", // from values.yaml
				}, chart.Values)
				cancel()
			}),
	)
	w.Run(ctx)
}

func TestRun_HappyPath_ValuesFromUrl(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/yaml, application/json", r.Header.Get(httpz.AcceptHeader))
		switch r.URL.Path {
		case "/file.json":
			w.Header().Set(httpz.ContentTypeHeader, "application/json")
			_, err := w.Write([]byte(`{"json":"value"}`))
			assert.NoError(t, err)
		case "/file.yaml":
			w.Header().Set(httpz.ContentTypeHeader, "application/yaml")
			_, err := w.Write([]byte(`yaml: val`))
			assert.NoError(t, err)
		case "/empty":
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer s.Close()
	w, helm, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id: projectPath,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_Url{
					Url: &agentcfg.ChartValuesUrlCF{
						Url:         s.URL + "/file.yaml",
						PollPeriod:  durationpb.New(defaultUrlValuePollPeriod),
						MaxFileSize: &maxFileSize,
					},
				},
			},
			{
				From: &agentcfg.ChartValuesCF_Url{
					Url: &agentcfg.ChartValuesUrlCF{
						Url:         s.URL + "/file.json",
						PollPeriod:  durationpb.New(defaultUrlValuePollPeriod),
						MaxFileSize: &maxFileSize,
					},
				},
			},
			{
				From: &agentcfg.ChartValuesCF_Url{
					Url: &agentcfg.ChartValuesUrlCF{
						Url:         s.URL + "/empty",
						PollPeriod:  durationpb.New(defaultUrlValuePollPeriod),
						MaxFileSize: &maxFileSize,
					},
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "/**",
				},
			},
		},
	}
	gomock.InOrder(
		watcher.EXPECT().
			Watch(gomock.Any(), matcher.ProtoEq(nil, req), gomock.Any()).
			Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
				callback(ctx, rpc.ObjectsToSynchronizeData{
					CommitId:  revision,
					ProjectId: testhelpers.ProjectId,
					Sources: []rpc.ObjectSource{
						{
							Name: "Chart.yaml",
							Data: []byte(`apiVersion: v2
name: test1
version: 1`),
						},
					},
				})
				<-ctx.Done()
			}),
		helm.EXPECT().
			History(releaseName).
			Return(nil, driver.ErrReleaseNotFound),
		helm.EXPECT().
			Install(gomock.Any(), gomock.Any(), gomock.Any(), InstallConfig{
				Namespace:   defaultNamespace,
				ReleaseName: releaseName,
			}).
			Do(func(ctx context.Context, chart *chart.Chart, vals ChartValues, cfg InstallConfig) {
				assert.Equal(t, ChartValues{
					"yaml": "val",   // from file.yaml
					"json": "value", // from file.json
				}, vals)
				cancel()
			}),
	)
	w.Run(ctx)
}

func TestRun_MissingValueFile(t *testing.T) {
	w, _, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id: projectPath,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_File{
					File: &agentcfg.ChartValuesFileCF{
						ProjectId: &projectPath,
						File:      "prod.yaml",
					},
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "/**",
				},
			},
			{
				Path: &rpc.PathCF_File{
					File: "prod.yaml",
				},
			},
		},
	}
	watcher.EXPECT().
		Watch(gomock.Any(), matcher.ProtoEq(nil, req), gomock.Any()).
		Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
			callback(ctx, rpc.ObjectsToSynchronizeData{
				CommitId:  revision,
				ProjectId: testhelpers.ProjectId,
				Sources: []rpc.ObjectSource{
					{
						Name: "Chart.yaml",
						Data: []byte(`apiVersion: v2
name: test1
version: 1`),
					},
				},
			})
			<-ctx.Done()
		})
	w.Run(ctx)
}

func TestRun_MissingValuesFromAnotherProject(t *testing.T) {
	project2Path := "proj2"
	w, _, watcher := setupWorker(t, &agentcfg.ChartCF{
		ReleaseName: releaseName,
		Source: &agentcfg.ChartSourceCF{
			Source: &agentcfg.ChartSourceCF_Project{
				Project: &agentcfg.ChartProjectSourceCF{
					Id: projectPath,
				},
			},
		},
		Values: []*agentcfg.ChartValuesCF{
			{
				From: &agentcfg.ChartValuesCF_File{
					File: &agentcfg.ChartValuesFileCF{
						ProjectId: &project2Path,
						Ref:       &agentcfg.GitRefCF{Ref: &agentcfg.GitRefCF_Branch{Branch: "prod"}},
						File:      "prod.yaml",
					},
				},
			},
		},
		Namespace:  &defaultNamespace,
		MaxHistory: &maxHistory,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req1 := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: projectPath,
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_Glob{
					Glob: "/**",
				},
			},
		},
	}
	req2 := &rpc.ObjectsToSynchronizeRequest{
		ProjectId: project2Path,
		Ref:       &rpc.GitRefCF{Ref: &rpc.GitRefCF_Branch{Branch: "prod"}},
		Paths: []*rpc.PathCF{
			{
				Path: &rpc.PathCF_File{
					File: "prod.yaml",
				},
			},
		},
	}
	watcher.EXPECT().
		Watch(gomock.Any(), matcher.ProtoEq(nil, req1), gomock.Any()).
		Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
			callback(ctx, rpc.ObjectsToSynchronizeData{
				CommitId:  revision,
				ProjectId: testhelpers.ProjectId,
				Sources: []rpc.ObjectSource{
					{
						Name: "Chart.yaml",
						Data: []byte(`apiVersion: v2
name: test1
version: 1`),
					},
					{
						Name: "values.yaml",
						Data: []byte(`x: z`),
					},
				},
			})
			<-ctx.Done()
		})
	watcher.EXPECT().
		Watch(gomock.Any(), matcher.ProtoEq(nil, req2), gomock.Any()).
		Do(func(ctx context.Context, req *rpc.ObjectsToSynchronizeRequest, callback rpc.ObjectsToSynchronizeCallback) {
			<-ctx.Done()
		})
	w.Run(ctx)
}

func setupWorker(t *testing.T, chartCfg *agentcfg.ChartCF) (*worker, *MockHelm, *mock_rpc.MockObjectsToSynchronizeWatcherInterface) {
	ctrl := gomock.NewController(t)
	watcher := mock_rpc.NewMockObjectsToSynchronizeWatcherInterface(ctrl)
	helm := NewMockHelm(ctrl)

	w := &worker{
		log:               zaptest.NewLogger(t),
		chartCfg:          chartCfg,
		installPollConfig: testhelpers.NewPollConfig(time.Minute)(),
		helm:              helm,
		httpClient:        http.DefaultTransport,
		objWatcher:        watcher,
	}
	return w, helm, watcher
}
