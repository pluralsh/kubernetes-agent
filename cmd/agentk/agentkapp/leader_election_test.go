package agentkapp

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"k8s.io/apimachinery/pkg/util/wait"
)

func TestLeaderElection_NotLeader_NoModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, _, mockElector, _ := setup(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var wg wait.Group
	defer wg.Wait()

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			<-ctx.Done()
			onStoppedLeading()
		})

	// WHEN
	wg.StartWithContext(ctx, lr.Run)

	// THEN
	cancel()
}

func TestLeaderElection_NotLeader_OneModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, lmw, mockElector, mockModule := setup(t)

	// contexts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	moduleCtx, moduleCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer moduleCancel()

	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	c := &agentcfg.AgentConfiguration{}

	var wg wait.Group
	defer wg.Wait()

	// setup mock expectations
	mockElector.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			<-ctx.Done()
			onStoppedLeading()
		})
	mockModule.EXPECT().IsRunnableConfiguration(gomock.Any()).
		DoAndReturn(func(_ *agentcfg.AgentConfiguration) bool {
			// shut down the module
			close(cfg)
			return true
		})

	// WHEN
	wg.StartWithContext(ctx, lr.Run)
	wg.Start(func() {
		lmw.Run(moduleCtx, cfg)
	})
	cfg <- c

	// THEN
	// we need to give the leader module wrapper some time to register the module etc.
	// and we don't currently have a way to properly wait for that, so we are just going to sleep a little bit ...
	time.Sleep(500 * time.Millisecond)
	cancel()
}

func TestLeaderElection_Leader_NoModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, _, mockElector, _ := setup(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var wg wait.Group
	defer wg.Wait()

	startedLeading := make(chan struct{})

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			onStartedLeading()
			close(startedLeading)
			<-ctx.Done()
			onStoppedLeading()
		})

	// WHEN
	wg.StartWithContext(ctx, lr.Run)

	// THEN
	select {
	case <-ctx.Done():
		require.FailNow(t, ctx.Err().Error())
	case <-startedLeading:
	}
	cancel()
}

func TestLeaderElection_Leader_OneModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, lmw, mockElector, mockModule := setup(t)

	// contexts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	moduleCtx, moduleCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer moduleCancel()

	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	c := &agentcfg.AgentConfiguration{}

	var wg wait.Group
	defer wg.Wait()

	complete := make(chan struct{})

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			onStartedLeading()
			<-ctx.Done()
			onStoppedLeading()
		})
	mockModule.EXPECT().IsRunnableConfiguration(gomock.Any()).
		Return(true)
	mockModule.EXPECT().Run(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ <-chan *agentcfg.AgentConfiguration) error {
			// shut down the module
			close(cfg)
			close(complete)
			return nil
		})

	// WHEN
	wg.StartWithContext(ctx, lr.Run)
	wg.Start(func() {
		lmw.Run(moduleCtx, cfg)
	})
	cfg <- c

	// THEN
	select {
	case <-ctx.Done():
		require.FailNow(t, ctx.Err().Error())
	case <-complete:
		// we need to give the leader runner some time to unregister the module,
		// and we don't currently have a way to properly wait for that, so we are just going to sleep a little bit ...
		time.Sleep(500 * time.Millisecond)
	}
	cancel()
}

func TestLeaderElection_LeaderNotLeaderLeader_NoModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, _, mockElector, _ := setup(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var wg wait.Group
	defer wg.Wait()

	secondStartedLeading := make(chan struct{})

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			// first leading
			onStartedLeading()
			onStoppedLeading()

			// second leading
			onStartedLeading()
			close(secondStartedLeading)
			<-ctx.Done()
			onStoppedLeading()
		})

	// WHEN
	wg.StartWithContext(ctx, lr.Run)

	// THEN
	select {
	case <-ctx.Done():
		require.FailNow(t, ctx.Err().Error())
	case <-secondStartedLeading:
	}
	cancel()
}

func TestLeaderElection_LeaderNotLeaderLeader_OneModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, lmw, mockElector, mockModule := setup(t)

	// contexts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	moduleCtx, moduleCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer moduleCancel()

	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	c := &agentcfg.AgentConfiguration{}

	var wg wait.Group
	defer wg.Wait()

	firstModuleStarted := make(chan struct{})
	complete := make(chan struct{})

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			// first leading
			onStartedLeading()
			<-firstModuleStarted
			onStoppedLeading()

			// second leading
			onStartedLeading()
			<-ctx.Done()
			onStoppedLeading()
		})

	gomock.InOrder(
		mockModule.EXPECT().IsRunnableConfiguration(gomock.Any()).
			Return(true),
		mockModule.EXPECT().Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1 := <-cfgm
				assert.Same(t, c, c1)
				close(firstModuleStarted)

				<-ctx.Done()
				return nil
			}),
		mockModule.EXPECT().Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1 := <-cfgm
				assert.Same(t, c, c1)

				// shut down the module
				close(cfg)
				close(complete)
				return nil
			}),
	)

	// WHEN
	wg.StartWithContext(ctx, lr.Run)
	wg.Start(func() {
		lmw.Run(moduleCtx, cfg)
	})
	cfg <- c

	// THEN
	select {
	case <-ctx.Done():
		require.FailNow(t, ctx.Err().Error())
	case <-complete:
		// we need to give the leader runner some time to unregister the module,
		// and we don't currently have a way to properly wait for that, so we are just going to sleep a little bit ...
		time.Sleep(500 * time.Millisecond)
	}
	cancel()
}

func TestLeaderElection_Leader_OneModule_RunnableNotRunnableRunnable_Shutdown(t *testing.T) {
	// GIVEN
	lr, lmw, mockElector, mockModule := setup(t)

	// contexts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	moduleCtx, moduleCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer moduleCancel()

	cfg := make(chan *agentcfg.AgentConfiguration)
	firstRunnableCfg := &agentcfg.AgentConfiguration{}
	secondRunnableCfg := &agentcfg.AgentConfiguration{}
	notRunnableCfg := &agentcfg.AgentConfiguration{}

	var wg wait.Group
	defer wg.Wait()

	consumedFirstCfg := make(chan struct{})
	complete := make(chan struct{})

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			onStartedLeading()
			<-ctx.Done()
			onStoppedLeading()
		})
	gomock.InOrder(
		mockModule.EXPECT().IsRunnableConfiguration(firstRunnableCfg).
			Return(true),
		mockModule.EXPECT().Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c := <-cfgm
				assert.Same(t, firstRunnableCfg, c)
				close(consumedFirstCfg)
				<-cfgm
				return nil
			}),
		mockModule.EXPECT().IsRunnableConfiguration(notRunnableCfg).
			Return(false),
		mockModule.EXPECT().IsRunnableConfiguration(secondRunnableCfg).
			Return(true),
		mockModule.EXPECT().Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c := <-cfgm
				assert.Same(t, secondRunnableCfg, c)

				// shut down the module
				close(cfg)
				close(complete)
				return nil
			}),
	)

	// WHEN
	wg.StartWithContext(ctx, lr.Run)
	wg.Start(func() {
		lmw.Run(moduleCtx, cfg)
	})

	cfg <- firstRunnableCfg
	<-consumedFirstCfg
	cfg <- notRunnableCfg
	cfg <- secondRunnableCfg

	// THEN
	select {
	case <-ctx.Done():
		require.FailNow(t, ctx.Err().Error())
	case <-complete:
		// we need to give the leader runner some time to unregister the module,
		// and we don't currently have a way to properly wait for that, so we are just going to sleep a little bit ...
		time.Sleep(500 * time.Millisecond)
	}
	cancel()
}

func TestLeaderElection_NotLeaderRunnableRunnableLeader_OneModule_Shutdown(t *testing.T) {
	// GIVEN
	lr, lmw, mockElector, mockModule := setup(t)

	// contexts
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	moduleCtx, moduleCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer moduleCancel()

	cfg := make(chan *agentcfg.AgentConfiguration)
	firstRunnableCfg := &agentcfg.AgentConfiguration{}
	secondRunnableCfg := &agentcfg.AgentConfiguration{}

	var wg wait.Group
	defer wg.Wait()

	validatedSecondCfg := make(chan struct{})
	complete := make(chan struct{})

	// setup mock expectations
	mockElector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			<-validatedSecondCfg
			// Give the leader module wrapper some time to register the module with the second runnable config
			// before getting the leadership
			time.Sleep(500 * time.Millisecond)

			onStartedLeading()
			<-ctx.Done()
			onStoppedLeading()
		})
	gomock.InOrder(
		mockModule.EXPECT().IsRunnableConfiguration(firstRunnableCfg).
			Return(true),
		mockModule.EXPECT().IsRunnableConfiguration(secondRunnableCfg).
			DoAndReturn(func(_ *agentcfg.AgentConfiguration) bool {
				close(validatedSecondCfg)
				return true
			}),
		mockModule.EXPECT().Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c := <-cfgm
				assert.Same(t, secondRunnableCfg, c)

				// shut down the module
				close(cfg)
				close(complete)
				return nil
			}),
	)

	// WHEN
	wg.StartWithContext(ctx, lr.Run)
	wg.Start(func() {
		lmw.Run(moduleCtx, cfg)
	})

	cfg <- firstRunnableCfg
	cfg <- secondRunnableCfg

	// THEN
	select {
	case <-ctx.Done():
		require.FailNow(t, ctx.Err().Error())
	case <-complete:
		// we need to give the leader runner some time to unregister the module,
		// and we don't currently have a way to properly wait for that, so we are just going to sleep a little bit ...
		time.Sleep(500 * time.Millisecond)
	}
	cancel()
}

func setup(t *testing.T) (*leaderRunner, *leaderModuleWrapper, *MockLeaderElector, *mock_modagent.MockLeaderModule) {
	ctrl := gomock.NewController(t)

	mockModule := mock_modagent.NewMockLeaderModule(ctrl)
	mockElector := NewMockLeaderElector(ctrl)
	lr := newLeaderRunner(mockElector)
	lmw := newLeaderModuleWrapper(mockModule, lr)
	return lr, lmw, mockElector, mockModule
}
