package agentkapp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/module/modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_modagent"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/pkg/agentcfg"
	"k8s.io/apimachinery/pkg/util/wait"
)

var (
	_ modagent.Module = (*leaderModuleWrapper)(nil)
)

func TestLMW_DefaultAndValidateConfiguration_IsDelegated(t *testing.T) {
	w, _, m := setupLMW(t)
	c := &agentcfg.AgentConfiguration{}
	m.EXPECT().
		DefaultAndValidateConfiguration(c).
		Return(errors.New("boom"))
	err := w.DefaultAndValidateConfiguration(c)
	assert.EqualError(t, err, "boom")
}

func TestLMW_Name_IsDelegated(t *testing.T) {
	w, _, m := setupLMW(t)
	m.EXPECT().
		Name().
		Return("name1")
	assert.Equal(t, "name1", w.Name())
}

func TestLMW_Run_NotRunnableConfig(t *testing.T) {
	w, _, m := setupLMW(t)
	c := &agentcfg.AgentConfiguration{}
	m.EXPECT().
		IsRunnableConfiguration(c).
		Return(false)
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c
	close(cfg)
	err := w.Run(context.Background(), cfg)
	require.NoError(t, err)
}

func TestLMW_Run_RunnableThenNotRunnableConfig(t *testing.T) {
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	c2 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		r.EXPECT().
			RunWhenLeader(gomock.Any()).
			DoAndReturn(func(f func(context.Context)) func() {
				var wg wait.Group
				ctx, cancel := context.WithCancel(context.Background())
				wg.StartWithContext(ctx, f)
				return func() {
					cancel()
					wg.Wait()
					stopCalled = true
				}
			}),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1m := <-cfgm
				assert.Same(t, c1, c1m)
				cfg <- c2
				close(cfg)
				_, ok := <-cfgm
				assert.False(t, ok)
				return nil
			}),
		m.EXPECT().
			IsRunnableConfiguration(c2).
			Return(false),
	)
	err := w.Run(context.Background(), cfg)
	require.NoError(t, err)
	assert.True(t, stopCalled)
}

func TestLMW_Run_RunnableThenNotRunnableThenRunnableConfig(t *testing.T) {
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	c2 := &agentcfg.AgentConfiguration{}
	c3 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	r.EXPECT().
		RunWhenLeader(gomock.Any()).
		DoAndReturn(func(f func(context.Context)) func() {
			var wg wait.Group
			ctx, cancel := context.WithCancel(context.Background())
			wg.StartWithContext(ctx, f)
			return func() {
				cancel()
				wg.Wait()
				stopCalled = true
			}
		}).Times(2)
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1m := <-cfgm
				assert.Same(t, c1, c1m)
				cfg <- c2
				_, ok := <-cfgm
				assert.False(t, ok)
				return nil
			}),
		m.EXPECT().
			IsRunnableConfiguration(c2).
			DoAndReturn(func(_ *agentcfg.AgentConfiguration) bool {
				cfg <- c3
				return false
			}),
		m.EXPECT().
			IsRunnableConfiguration(c3).
			Return(true),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c3m := <-cfgm
				assert.Same(t, c3, c3m)
				close(cfg)
				_, ok := <-cfgm
				assert.False(t, ok)
				return nil
			}),
	)
	err := w.Run(context.Background(), cfg)
	require.NoError(t, err)
	assert.True(t, stopCalled)
}

func TestLMW_Run_RunnableThenNotRunnableStopError(t *testing.T) {
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	c2 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		r.EXPECT().
			RunWhenLeader(gomock.Any()).
			DoAndReturn(func(f func(context.Context)) func() {
				var wg wait.Group
				ctx, cancel := context.WithCancel(context.Background())
				wg.StartWithContext(ctx, f)
				return func() {
					cancel()
					wg.Wait()
					stopCalled = true
				}
			}),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1m := <-cfgm
				assert.Same(t, c1, c1m)
				cfg <- c2
				_, ok := <-cfgm
				assert.False(t, ok)
				return errors.New("boom")
			}),
		m.EXPECT().
			IsRunnableConfiguration(c2).
			Return(false),
	)
	err := w.Run(context.Background(), cfg)
	assert.EqualError(t, err, "boom")
	assert.True(t, stopCalled)
}

func TestLMW_Run_EarlyReturnNoError(t *testing.T) {
	t.SkipNow() // TODO https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/-/issues/318
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		r.EXPECT().
			RunWhenLeader(gomock.Any()).
			DoAndReturn(func(f func(context.Context)) func() {
				var wg wait.Group
				ctx, cancel := context.WithCancel(context.Background())
				wg.StartWithContext(ctx, f)
				return func() {
					cancel()
					wg.Wait()
					stopCalled = true
				}
			}),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				close(cfg)
				return nil
			}),
	)
	err := w.Run(context.Background(), cfg)
	require.NoError(t, err)
	assert.True(t, stopCalled)
}

func TestLMW_Run_EarlyReturnNoErrorSameConfig(t *testing.T) {
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	r.EXPECT().
		RunWhenLeader(gomock.Any()).
		DoAndReturn(func(f func(context.Context)) func() {
			var wg wait.Group
			ctx, cancel := context.WithCancel(context.Background())
			wg.StartWithContext(ctx, f)
			return func() {
				cancel()
				wg.Wait()
				stopCalled = true
			}
		}).Times(2)
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				return nil
			}),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1m := <-cfgm
				assert.Same(t, c1, c1m)
				close(cfg)
				_, ok := <-cfgm
				assert.False(t, ok)
				return nil
			}),
	)
	err := w.Run(context.Background(), cfg)
	require.NoError(t, err)
	assert.True(t, stopCalled)
}

func TestLMW_Run_EarlyReturnNoErrorMoreConfig(t *testing.T) {
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	c2 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	r.EXPECT().
		RunWhenLeader(gomock.Any()).
		DoAndReturn(func(f func(context.Context)) func() {
			var wg wait.Group
			ctx, cancel := context.WithCancel(context.Background())
			wg.StartWithContext(ctx, f)
			return func() {
				cancel()
				wg.Wait()
				stopCalled = true
			}
		}).Times(2)
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c1m := <-cfgm
				assert.Same(t, c1, c1m)
				go func() {
					time.Sleep(time.Second)
					cfg <- c2
				}()
				return nil
			}),
		m.EXPECT().
			IsRunnableConfiguration(c2).
			Return(true),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				c2m := <-cfgm
				assert.Same(t, c2, c2m)
				close(cfg)
				_, ok := <-cfgm
				assert.False(t, ok)
				return nil
			}),
	)
	err := w.Run(context.Background(), cfg)
	require.NoError(t, err)
	assert.True(t, stopCalled)
}

func TestLMW_Run_EarlyReturnError(t *testing.T) {
	w, r, m := setupLMW(t)
	c1 := &agentcfg.AgentConfiguration{}
	cfg := make(chan *agentcfg.AgentConfiguration, 1)
	cfg <- c1
	stopCalled := false
	gomock.InOrder(
		m.EXPECT().
			IsRunnableConfiguration(c1).
			Return(true),
		r.EXPECT().
			RunWhenLeader(gomock.Any()).
			DoAndReturn(func(f func(context.Context)) func() {
				var wg wait.Group
				ctx, cancel := context.WithCancel(context.Background())
				wg.StartWithContext(ctx, f)
				return func() {
					cancel()
					wg.Wait()
					stopCalled = true
				}
			}),
		m.EXPECT().
			Run(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cfgm <-chan *agentcfg.AgentConfiguration) error {
				return errors.New("boom")
			}),
	)
	err := w.Run(context.Background(), cfg)
	assert.EqualError(t, err, "boom")
	assert.True(t, stopCalled)
}

func setupLMW(t *testing.T) (*leaderModuleWrapper, *MockRunner, *mock_modagent.MockLeaderModule) {
	ctrl := gomock.NewController(t)
	m := mock_modagent.NewMockLeaderModule(ctrl)
	r := NewMockRunner(ctrl)
	w := newLeaderModuleWrapper(m, r)
	return w, r, m
}
