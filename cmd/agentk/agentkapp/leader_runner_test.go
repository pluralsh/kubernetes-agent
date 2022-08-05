package agentkapp

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_modagent"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ Runner        = (*leaderRunner)(nil)
	_ LeaderElector = (*leaseLeaderElector)(nil)
)

func TestLR_MaybeWrapModule_Module(t *testing.T) {
	_, _, lr, _ := setupLR(t)
	ctrl := gomock.NewController(t)
	module := mock_modagent.NewMockModule(ctrl)
	wrapped := lr.MaybeWrapModule(module)
	assert.Same(t, module, wrapped)
}

func TestLR_MaybeWrapModule_LeaderModule(t *testing.T) {
	_, _, lr, _ := setupLR(t)
	ctrl := gomock.NewController(t)
	module := mock_modagent.NewMockLeaderModule(ctrl)
	wrapped := lr.MaybeWrapModule(module)
	assert.IsType(t, (*leaderModuleWrapper)(nil), wrapped)
}

func TestLR_RunNoLeaderStop(t *testing.T) {
	ctx, cancel, lr, elector := setupLR(t)
	elector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			<-ctx.Done()
			onStoppedLeading()
		})
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.StartWithContext(ctx, lr.Run)
	stop := lr.RunWhenLeader(func(ctx context.Context) {
		assert.Fail(t, "unexpected call")
	})
	stop()
}

func TestLR_RunNoLeaderRunStop(t *testing.T) {
	ctx, cancel, lr, elector := setupLR(t)
	elector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			<-ctx.Done()
			onStoppedLeading()
		})
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.StartWithContext(ctx, lr.Run)
	stop1 := lr.RunWhenLeader(func(ctx context.Context) {
		assert.Fail(t, "unexpected call")
	})
	stop2 := lr.RunWhenLeader(func(ctx context.Context) {
		assert.Fail(t, "unexpected call")
	})
	stop1()
	stop2()
}

func TestLR_RunLeaderStop(t *testing.T) {
	ctx, cancel, lr, elector := setupLR(t)
	elector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			onStartedLeading()
			<-ctx.Done()
			onStoppedLeading()
		})
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.StartWithContext(ctx, lr.Run)
	started := make(chan struct{})
	stopped := make(chan struct{})
	stop := lr.RunWhenLeader(func(ctx context.Context) {
		close(started)
		<-ctx.Done()
		close(stopped)
	})
	<-started
	stop()
	<-stopped
}

func TestLR_RunLeaderRunStopStop(t *testing.T) {
	ctx, cancel, lr, elector := setupLR(t)
	elector.EXPECT().
		Run(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
			onStartedLeading()
			<-ctx.Done()
			onStoppedLeading()
		})
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.StartWithContext(ctx, lr.Run)
	started1 := make(chan struct{})
	stopped1 := make(chan struct{})
	stop1 := lr.RunWhenLeader(func(ctx context.Context) {
		close(started1)
		<-ctx.Done()
		close(stopped1)
	})
	<-started1
	started2 := make(chan struct{})
	stopped2 := make(chan struct{})
	stop2 := lr.RunWhenLeader(func(ctx context.Context) {
		close(started2)
		<-ctx.Done()
		close(stopped2)
	})
	<-started2
	stop1()
	stop2()
	<-stopped1
	<-stopped2
}

func TestLR_RunLeaderNotLeaderLeaderStop(t *testing.T) {
	ctx, cancel, lr, elector := setupLR(t)

	var wgStart sync.WaitGroup
	wgStart.Add(1)
	var wgStop sync.WaitGroup
	wgStop.Add(1)
	callStop := make(chan struct{})

	gomock.InOrder(
		elector.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
				onStartedLeading()
				wgStart.Wait()
				onStoppedLeading()
			}),
		elector.EXPECT().
			Run(gomock.Any(), gomock.Any(), gomock.Any()).
			Do(func(ctx context.Context, onStartedLeading, onStoppedLeading func()) {
				wgStop.Wait()
				wgStart.Add(1)
				wgStop.Add(1)
				onStartedLeading()
				wgStart.Wait()
				close(callStop)
				<-ctx.Done()
				onStoppedLeading()
			}),
	)
	var wg wait.Group
	defer wg.Wait()
	defer cancel()
	wg.StartWithContext(ctx, lr.Run)

	stop := lr.RunWhenLeader(func(ctx context.Context) {
		wgStart.Done()
		<-ctx.Done()
		wgStop.Done()
	})
	<-callStop
	stop()
}

func setupLR(t *testing.T) (context.Context, context.CancelFunc, *leaderRunner, *MockLeaderElector) {
	ctrl := gomock.NewController(t)
	elector := NewMockLeaderElector(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel, newLeaderRunner(elector), elector
}
