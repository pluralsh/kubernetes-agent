package redistool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	errKey = "test1"
)

func TestErrCacher_GetError_ReturnsNilOnClientError(t *testing.T) {
	ec, mock := setupNormal(t)
	mock.ExpectGet(errKey).SetErr(errors.New("boom"))
	err := ec.GetError(context.Background(), errKey)
	require.NoError(t, err)
}

func TestErrCacher_GetError_ReturnsNilOnClientNil(t *testing.T) {
	ec, mock := setupNormal(t)
	mock.ExpectGet(errKey).SetVal("")
	err := ec.GetError(context.Background(), errKey)
	require.NoError(t, err)
}

func TestErrCacher_GetError_ReturnsNilOnUnmarshalFail(t *testing.T) {
	ec, mock := setupError(t)
	mock.ExpectGet(errKey).SetVal("boom")
	err := ec.GetError(context.Background(), errKey)
	require.NoError(t, err)
}

func TestErrCacher_GetError_ReturnsCachedError(t *testing.T) {
	ec, mock := setupNormal(t)
	mock.ExpectGet(errKey).SetVal("boom")
	err := ec.GetError(context.Background(), errKey)
	require.EqualError(t, err, "boom")
}

func TestErrCacher_CacheError_HappyPath(t *testing.T) {
	ec, mock := setupNormal(t)
	mock.ExpectSet(errKey, []byte("boom"), time.Minute).SetVal("")
	ec.CacheError(context.Background(), errKey, errors.New("boom"), time.Minute)
}

func TestErrCacher_CacheError_MarshalError(t *testing.T) {
	ec, _ := setupError(t)
	ec.CacheError(context.Background(), errKey, errors.New("boom"), time.Minute)
}

func setupNormal(t *testing.T) (*ErrCacher, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()
	ec := &ErrCacher{
		Log:    zaptest.NewLogger(t),
		Client: client,
		ErrMarshaler: testErrMarshaler{
			marshal: func(err error) ([]byte, error) {
				return []byte(err.Error()), nil
			},
			unmarshal: func(data []byte) (error, error) {
				return errors.New(string(data)), nil
			},
		},
		KeyToRedisKey: func(key interface{}) string {
			return key.(string)
		},
	}
	return ec, mock
}

func setupError(t *testing.T) (*ErrCacher, redismock.ClientMock) {
	client, mock := redismock.NewClientMock()
	ec := &ErrCacher{
		Log:    zaptest.NewLogger(t),
		Client: client,
		ErrMarshaler: testErrMarshaler{
			marshal: func(err error) ([]byte, error) {
				return nil, errors.New("marshal error")
			},
			unmarshal: func(data []byte) (error, error) {
				return nil, errors.New("unmarshal error")
			},
		},
		KeyToRedisKey: func(key interface{}) string {
			return key.(string)
		},
	}
	return ec, mock
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
