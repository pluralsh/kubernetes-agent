package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15/internal/tool/testing/mock_cache"
)

func TestGetItem_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher(ctrl)
	errCacher.EXPECT().GetError(gomock.Any(), key)
	c := NewWithError(time.Minute, time.Minute, errCacher, alwaysCache)
	item, err := c.GetItem(context.Background(), key, func() (interface{}, error) {
		return itemVal, nil
	})
	require.NoError(t, err)
	assert.Equal(t, itemVal, item)

	item, err = c.GetItem(context.Background(), key, func() (interface{}, error) {
		t.FailNow()
		return nil, nil
	})
	require.NoError(t, err)
	assert.Equal(t, itemVal, item)
}

func TestGetItem_CacheableError(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher(ctrl)
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
	c := NewWithError(time.Second, time.Minute, errCacher, alwaysCache)
	_, err := c.GetItem(context.Background(), key, func() (interface{}, error) {
		return nil, errToCache
	})
	assert.EqualError(t, err, "boom")

	_, err = c.GetItem(context.Background(), key, func() (interface{}, error) {
		t.FailNow()
		return nil, nil
	})
	assert.EqualError(t, err, "boom")
}

func TestGetItem_NonCacheableError(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher(ctrl)
	errCacher.EXPECT().
		GetError(gomock.Any(), key).
		Times(2)
	c := NewWithError(time.Minute, time.Minute, errCacher, func(err error) bool {
		return false
	})
	_, err := c.GetItem(context.Background(), key, func() (interface{}, error) {
		return nil, errors.New("boom")
	})
	assert.EqualError(t, err, "boom")

	_, err = c.GetItem(context.Background(), key, func() (interface{}, error) {
		return nil, errors.New("bAAm")
	})
	assert.EqualError(t, err, "bAAm")
}

func TestGetItem_Context(t *testing.T) {
	ctrl := gomock.NewController(t)
	errCacher := mock_cache.NewMockErrCacher(ctrl)
	errCacher.EXPECT().GetError(gomock.Any(), key)
	c := NewWithError(time.Minute, time.Minute, errCacher, alwaysCache)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-start
		_, err := c.GetItem(ctx, key, func() (interface{}, error) {
			return "Stalemate. No-oh, too late, too late", nil
		})
		assert.Equal(t, context.Canceled, err)
	}()
	item, err := c.GetItem(context.Background(), key, func() (interface{}, error) {
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
