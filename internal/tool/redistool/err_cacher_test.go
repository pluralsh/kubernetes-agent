package redistool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/matcher"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_tool"
	"go.uber.org/zap/zaptest"
)

const (
	errKey = "test1"
)

func TestErrCacher_GetError_ReturnsNilOnClientError(t *testing.T) {
	ec, mock, rep := setupNormal(t)
	mock.ExpectGet(errKey).SetErr(errors.New("boom"))
	rep.EXPECT().
		HandleProcessingError(gomock.Any(), gomock.Any(), "Failed to get cached error from Redis", matcher.ErrorEq("boom"))
	err := ec.GetError(context.Background(), errKey)
	require.NoError(t, err)
}

func TestErrCacher_GetError_ReturnsNilOnClientNil(t *testing.T) {
	ec, mock, _ := setupNormal(t)
	mock.ExpectGet(errKey).SetVal("")
	err := ec.GetError(context.Background(), errKey)
	require.NoError(t, err)
}

func TestErrCacher_GetError_ReturnsNilOnUnmarshalFail(t *testing.T) {
	ec, mock, rep := setupError(t)
	mock.ExpectGet(errKey).SetVal("boom")
	rep.EXPECT().
		HandleProcessingError(gomock.Any(), gomock.Any(), "Failed to unmarshal cached error", matcher.ErrorEq("unmarshal error"))
	err := ec.GetError(context.Background(), errKey)
	require.NoError(t, err)
}

func TestErrCacher_GetError_ReturnsCachedError(t *testing.T) {
	ec, mock, _ := setupNormal(t)
	mock.ExpectGet(errKey).SetVal("boom")
	err := ec.GetError(context.Background(), errKey)
	require.EqualError(t, err, "boom")
}

func TestErrCacher_CacheError_HappyPath(t *testing.T) {
	ec, mock, _ := setupNormal(t)
	mock.ExpectSet(errKey, []byte("boom"), time.Minute).SetVal("")
	ec.CacheError(context.Background(), errKey, errors.New("boom"), time.Minute)
}

func TestErrCacher_CacheError_MarshalError(t *testing.T) {
	ec, _, rep := setupError(t)
	rep.EXPECT().
		HandleProcessingError(gomock.Any(), gomock.Any(), "Failed to marshal error for caching", matcher.ErrorEq("marshal error"))
	ec.CacheError(context.Background(), errKey, errors.New("boom"), time.Minute)
}

func setupNormal(t *testing.T) (*ErrCacher[string], redismock.ClientMock, *mock_tool.MockErrReporter) {
	client, mock := redismock.NewClientMock()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	ec := &ErrCacher[string]{
		Log:    zaptest.NewLogger(t),
		ErrRep: rep,
		Client: client,
		ErrMarshaler: testErrMarshaler{
			marshal: func(err error) ([]byte, error) {
				return []byte(err.Error()), nil
			},
			unmarshal: func(data []byte) (error, error) {
				return errors.New(string(data)), nil
			},
		},
		KeyToRedisKey: func(key string) string {
			return key
		},
	}
	return ec, mock, rep
}

func setupError(t *testing.T) (*ErrCacher[string], redismock.ClientMock, *mock_tool.MockErrReporter) {
	client, mock := redismock.NewClientMock()
	ctrl := gomock.NewController(t)
	rep := mock_tool.NewMockErrReporter(ctrl)
	ec := &ErrCacher[string]{
		Log:    zaptest.NewLogger(t),
		ErrRep: rep,
		Client: client,
		ErrMarshaler: testErrMarshaler{
			marshal: func(err error) ([]byte, error) {
				return nil, errors.New("marshal error")
			},
			unmarshal: func(data []byte) (error, error) {
				return nil, errors.New("unmarshal error")
			},
		},
		KeyToRedisKey: func(key string) string {
			return key
		},
	}
	return ec, mock, rep
}

type testErrMarshaler struct {
	marshal   func(err error) ([]byte, error)
	unmarshal func(data []byte) (error, error)
}

func (m testErrMarshaler) Marshal(err error) ([]byte, error) {
	return m.marshal(err)
}

func (m testErrMarshaler) Unmarshal(data []byte) (error, error) {
	return m.unmarshal(data)
}
