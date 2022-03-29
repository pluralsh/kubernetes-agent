package agentkapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/module/modagent"
	"go.uber.org/zap/zaptest"
)

const (
	testFeature modagent.Feature = -42
)

func TestFeatureTracker_FirstEnableCallsCallback(t *testing.T) {
	ft := newFeatureTracker(zaptest.NewLogger(t))
	called := 0
	ft.Subscribe(testFeature, func(enabled bool) {
		assert.True(t, enabled)
		called++
	})
	ft.ToggleFeature(testFeature, "bla", true)
	assert.EqualValues(t, 1, called)
}

func TestFeatureTracker_MultipleEnableCallCallbackOnce(t *testing.T) {
	t.Run("same consumer", func(t *testing.T) {
		ft := newFeatureTracker(zaptest.NewLogger(t))
		called := 0
		ft.Subscribe(testFeature, func(enabled bool) {
			assert.True(t, enabled)
			called++
		})
		ft.ToggleFeature(testFeature, "bla", true)
		ft.ToggleFeature(testFeature, "bla", true)
		assert.EqualValues(t, 1, called)
	})
	t.Run("different consumers", func(t *testing.T) {
		ft := newFeatureTracker(zaptest.NewLogger(t))
		called := 0
		ft.Subscribe(testFeature, func(enabled bool) {
			assert.True(t, enabled)
			called++
		})
		ft.ToggleFeature(testFeature, "bla1", true)
		ft.ToggleFeature(testFeature, "bla2", true)
		assert.EqualValues(t, 1, called)
	})
}

func TestFeatureTracker_DisableIsCalledOnce(t *testing.T) {
	t.Run("same consumer", func(t *testing.T) {
		ft := newFeatureTracker(zaptest.NewLogger(t))
		called := 0
		ft.ToggleFeature(testFeature, "bla", true)
		ft.Subscribe(testFeature, func(enabled bool) {
			assert.False(t, enabled)
			called++
		})
		ft.ToggleFeature(testFeature, "bla", false)
		assert.EqualValues(t, 1, called)
		ft.ToggleFeature(testFeature, "bla", false)
		assert.EqualValues(t, 1, called) // still one
	})
	t.Run("different consumers", func(t *testing.T) {
		ft := newFeatureTracker(zaptest.NewLogger(t))
		called := 0
		ft.ToggleFeature(testFeature, "bla1", true)
		ft.ToggleFeature(testFeature, "bla2", true)
		ft.Subscribe(testFeature, func(enabled bool) {
			assert.False(t, enabled)
			called++
		})
		ft.ToggleFeature(testFeature, "bla1", false)
		assert.Zero(t, called)
		ft.ToggleFeature(testFeature, "bla2", false)
		assert.EqualValues(t, 1, called)
	})
}
