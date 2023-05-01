package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_cache"
)

func TestGetItem_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher[int](ctrl)
	errCacher.EXPECT().GetError(gomock.Any(), key)
	c := NewWithError[int, int](time.Minute, time.Minute, errCacher, alwaysCache)
	item, err := c.GetItem(context.Background(), key, func() (int, error) {
		return itemVal, nil
	})
	require.NoError(t, err)
	assert.Equal(t, itemVal, item)

	item, err = c.GetItem(context.Background(), key, func() (int, error) {
		t.FailNow()
		return 0, nil
	})
	require.NoError(t, err)
	assert.Equal(t, itemVal, item)
}

func TestGetItem_CacheableError(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher[int](ctrl)
	errToCache := errors.New("boom")
	gomock.InOrder(
		errCacher.EXPECT().
			GetError(gomock.Any(), key),
		errCacher.EXPECT().
			CacheError(gomock.Any(), key, errToCache, time.Minute),
		errCacher.EXPECT().
			GetError(gomock.Any(), key).
			Return(errToCache),
	)
	c := NewWithError[int, int](time.Second, time.Minute, errCacher, alwaysCache)
	_, err := c.GetItem(context.Background(), key, func() (int, error) {
		return 0, errToCache
	})
	assert.EqualError(t, err, "boom")

	_, err = c.GetItem(context.Background(), key, func() (int, error) {
		t.FailNow()
		return 0, nil
	})
	assert.EqualError(t, err, "boom")
}

func TestGetItem_NonCacheableError(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher[int](ctrl)
	errCacher.EXPECT().
		GetError(gomock.Any(), key).
		Times(2)
	c := NewWithError[int, int](time.Minute, time.Minute, errCacher, func(err error) bool {
		return false
	})
	_, err := c.GetItem(context.Background(), key, func() (int, error) {
		return 0, errors.New("boom")
	})
	assert.EqualError(t, err, "boom")

	_, err = c.GetItem(context.Background(), key, func() (int, error) {
		return 0, errors.New("bAAm")
	})
	assert.EqualError(t, err, "bAAm")
}

func TestGetItem_Context(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher[int](ctrl)
	errCacher.EXPECT().GetError(gomock.Any(), key)
	c := NewWithError[int, int](time.Minute, time.Minute, errCacher, alwaysCache)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-start
		_, err := c.GetItem(ctx, key, func() (int, error) {
			return -itemVal, nil
		})
		assert.Equal(t, context.Canceled, err)
	}()
	item, err := c.GetItem(context.Background(), key, func() (int, error) {
		close(start)
		cancel()
		<-done
		return itemVal, nil
	})
	require.NoError(t, err)
	assert.Equal(t, itemVal, item)
}

func alwaysCache(err error) bool {
	return true
}
