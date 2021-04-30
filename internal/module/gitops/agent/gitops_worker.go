package agent

import (
	"context"

	"github.com/argoproj/gitops-engine/pkg/cache"
	"github.com/argoproj/gitops-engine/pkg/engine"
	"github.com/ash2k/stager"
	"github.com/go-logr/zapr"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/module/gitops/rpc"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/logz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/internal/tool/retry"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/pkg/agentcfg"
	"go.uber.org/zap"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"
)

type GitopsEngineFactory interface {
	New(engineOpts []engine.Option, cacheOpts []cache.UpdateSettingsFunc) engine.GitOpsEngine
}

type GitopsWorkerFactory interface {
	New(project *agentcfg.ManifestProjectCF) GitopsWorker
}

type GitopsWorker interface {
	Run(ctx context.Context)
}

type gitopsWorker struct {
	objWatcher           rpc.ObjectsToSynchronizeWatcherInterface
	engineFactory        GitopsEngineFactory
	engineBackoffFactory retry.BackoffManagerFactory
	synchronizerConfig
}

func (d *gitopsWorker) Run(ctx context.Context) {
	l := zapr.NewLogger(d.log)
	eng := d.engineFactory.New(
		[]engine.Option{
			engine.WithLogr(l),
		},
		[]cache.UpdateSettingsFunc{
			cache.SetPopulateResourceInfoHandler(populateResourceInfoHandler),
			cache.SetSettings(cache.Settings{
				ResourcesFilter: resourcesFilter{
					resourceInclusions: d.project.ResourceInclusions,
					resourceExclusions: d.project.ResourceExclusions,
				},
			}),
			cache.SetLogr(l),
		},
	)
	var stopEngine engine.StopFunc
	err := retry.PollWithBackoff(ctx, d.engineBackoffFactory(), true, 0, func() (error, retry.AttemptResult) {
		var err error
		stopEngine, err = eng.Run()
		if err != nil {
			d.log.Warn("engine.Run() failed", zap.Error(err))
			return nil, retry.Backoff
		}
		return nil, retry.Done
	})
	if err != nil {
		// context is done
		return
	}
	defer stopEngine()
	s := newSynchronizer(d.synchronizerConfig, eng)
	st := stager.New()
	stage := st.NextStage()
	stage.Go(func(ctx context.Context) error {
		s.run(ctx)
		return nil
	})
	stage = st.NextStage()
	stage.Go(func(ctx context.Context) error {
		req := &rpc.ObjectsToSynchronizeRequest{
			ProjectId: d.project.Id,
			Paths:     d.project.Paths,
		}
		d.objWatcher.Watch(ctx, req, func(ctx context.Context, data rpc.ObjectsToSynchronizeData) {
			s.setDesiredState(ctx, data)
		})
		return nil
	})
	_ = st.Run(ctx) // no errors possible
}

type defaultGitopsEngineFactory struct {
	kubeClientConfig *rest.Config
}

func (f *defaultGitopsEngineFactory) New(engineOpts []engine.Option, cacheOpts []cache.UpdateSettingsFunc) engine.GitOpsEngine {
	return engine.NewEngine(
		f.kubeClientConfig,
		cache.NewClusterCache(f.kubeClientConfig, cacheOpts...),
		engineOpts...,
	)
}

type defaultGitopsWorkerFactory struct {
	log                  *zap.Logger
	engineFactory        GitopsEngineFactory
	k8sClientGetter      resource.RESTClientGetter
	gitopsClient         rpc.GitopsClient
	watchBackoffFactory  retry.BackoffManagerFactory
	engineBackoffFactory retry.BackoffManagerFactory
}

func (m *defaultGitopsWorkerFactory) New(project *agentcfg.ManifestProjectCF) GitopsWorker {
	l := m.log.With(logz.ProjectId(project.Id))
	return &gitopsWorker{
		objWatcher: &rpc.ObjectsToSynchronizeWatcher{
			Log:          l,
			GitopsClient: m.gitopsClient,
			Backoff:      m.watchBackoffFactory,
		},
		engineFactory:        m.engineFactory,
		engineBackoffFactory: m.engineBackoffFactory,
		synchronizerConfig: synchronizerConfig{
			log:             l,
			project:         project,
			k8sClientGetter: m.k8sClientGetter,
		},
	}
}
