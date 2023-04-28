package gitlab_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/httpz"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/mock_gitlab"
	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v16/internal/tool/testing/testhelpers"
)

var (
	_ gitlab.ClientInterface = &gitlab.Client{}
	_ gitlab.ResponseHandler = gitlab.ResponseHandlerStruct{}
)

func TestRequestOptions(t *testing.T) {
	ctx, traceId := testhelpers.CtxWithSpanContext(t)
	c := mock_gitlab.SetupClient(t, "/ok", func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertRequestMethod(t, r, "CUSTOM_METHOD")
		testhelpers.AssertRequestAccept(t, r, "Bla")
		testhelpers.AssertAgentToken(t, r, testhelpers.AgentkToken)
		assert.Empty(t, r.Header[httpz.ContentTypeHeader])
		testhelpers.AssertCommonRequestParams(t, r, traceId)
		testhelpers.AssertJWTSignature(t, r)
		assert.Equal(t, "1", r.URL.Query().Get("a"))
		assert.Equal(t, "val1", r.URL.Query().Get("key"))
		assert.Equal(t, "val2", r.Header.Get("h1"))
	})
	c.Backend.RawQuery = "a=1"

	err := c.Do(ctx,
		gitlab.WithMethod("CUSTOM_METHOD"),
		gitlab.WithPath("/ok"),
		gitlab.WithQuery(url.Values{
			"key": []string{"val1"},
		}),
		gitlab.WithHeader(http.Header{
			"h1": []string{"val2"},
		}),
		gitlab.WithAgentToken(testhelpers.AgentkToken),
		gitlab.WithJWT(true),
		gitlab.WithResponseHandler(gitlab.ResponseHandlerStruct{
			AcceptHeader: "Bla",
			HandleFunc: func(resp *http.Response, err error) error {
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				// Do nothing
				return nil
			},
		}),
	)
	require.NoError(t, err)
}

func TestDoWithPath(t *testing.T) {
	tests := []struct {
		backendPath  string
		requestPath  string
		expectedPath string
	}{
		{
			backendPath:  "/ok",
			requestPath:  "",
			expectedPath: "/ok",
		},
		{
			backendPath:  "/ok/",
			requestPath:  "",
			expectedPath: "/ok/",
		},
		{
			backendPath:  "/ok/",
			requestPath:  "/",
			expectedPath: "/ok/",
		},
		{
			backendPath:  "",
			requestPath:  "",
			expectedPath: "/",
		},
		{
			backendPath:  "/",
			requestPath:  "",
			expectedPath: "/",
		},
		{
			backendPath:  "/",
			requestPath:  "/",
			expectedPath: "/",
		},
		{
			backendPath:  "",
			requestPath:  "/",
			expectedPath: "/",
		},
		{
			backendPath:  "/ok",
			requestPath:  "NONE",
			expectedPath: "/ok",
		},
		{
			backendPath:  "/a",
			requestPath:  "/b",
			expectedPath: "/a/b",
		},
		{
			backendPath:  "/a",
			requestPath:  "/b/",
			expectedPath: "/a/b/",
		},
		{
			backendPath:  "/a/",
			requestPath:  "/b",
			expectedPath: "/a/b",
		},
		{
			backendPath:  "/a/",
			requestPath:  "/b/",
			expectedPath: "/a/b/",
		},
	}
	for _, test := range tests {
		t.Run(test.backendPath+"-"+test.requestPath, func(t *testing.T) {
			c := mock_gitlab.SetupClient(t, "/", func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, test.expectedPath, r.URL.Path)
			})
			c.Backend.Path = test.backendPath

			opts := []gitlab.DoOption{
				gitlab.WithResponseHandler(gitlab.ResponseHandlerStruct{
					HandleFunc: func(resp *http.Response, err error) error {
						if err != nil {
							return err
						}
						defer resp.Body.Close()
						assert.EqualValues(t, http.StatusOK, resp.StatusCode)
						// Do nothing
						return nil
					},
				}),
			}
			if test.requestPath != "NONE" {
				opts = append(opts, gitlab.WithPath(test.requestPath))
			}
			err := c.Do(context.Background(), opts...)
			require.NoError(t, err)
		})
	}
}

func TestDoWithSlashAndSlashBackendPath(t *testing.T) {
	c := mock_gitlab.SetupClient(t, "/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ok/", r.URL.Path)
	})
	c.Backend.Path = "/ok/"

	err := c.Do(context.Background(),
		gitlab.WithPath("/"),
		gitlab.WithResponseHandler(gitlab.ResponseHandlerStruct{
			HandleFunc: func(resp *http.Response, err error) error {
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				assert.EqualValues(t, http.StatusOK, resp.StatusCode)
				// Do nothing
				return nil
			},
		}),
	)
	require.NoError(t, err)
}

func TestDoMasksQueryParametersInReturnedError(t *testing.T) {
	u, err := url.Parse("http://example.com:0")
	require.NoError(t, err)
	c := gitlab.NewClient(u, []byte(testhelpers.AuthSecretKey))
	err = c.Do(context.Background(),
		gitlab.WithoutRetries(),
		gitlab.WithPath("/abc"),
		gitlab.WithQuery(url.Values{"a": []string{"1", "2"}, "b": []string{"3"}}),
		gitlab.WithResponseHandler(gitlab.ResponseHandlerStruct{
			HandleFunc: func(resp *http.Response, err error) error {
				if err != nil {
					return err
				}
				t.Fail()
				return nil
			},
		}),
	)
	require.Error(t, err)
	ue, ok := err.(*url.Error) // nolint: errorlint
	require.True(t, ok)
	assert.Equal(t, "http://example.com:0/abc?a=x&b=x", ue.URL)
}

func TestJsonResponseHandler_Errors(t *testing.T) {
	tests := map[int]func(error) bool{
		http.StatusForbidden:    gitlab.IsForbidden,
		http.StatusUnauthorized: gitlab.IsUnauthorized,
		http.StatusNotFound:     gitlab.IsNotFound,
	}
	for statusCode, f := range tests {
		t.Run(strconv.Itoa(statusCode), func(t *testing.T) {
			c := mock_gitlab.SetupClient(t, "/bla", func(w http.ResponseWriter, r *http.Request) {
				testhelpers.AssertGetJsonRequest(t, r)
				w.WriteHeader(statusCode)
			})

			var resp interface{}

			err := c.Do(context.Background(),
				gitlab.WithPath("/bla"),
				gitlab.WithResponseHandler(gitlab.JsonResponseHandler(&resp)),
			)
			require.Error(t, err)
			assert.True(t, f(err))
			assert.True(t, errHasPath(err, "/bla"))
		})
	}
}

func TestJsonResponseHandler_HappyPath(t *testing.T) {
	c := mock_gitlab.SetupClient(t, "/ok", func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertGetJsonRequest(t, r)
		testhelpers.RespondWithJSON(t, w, 42)
	})

	var resp interface{}

	err := c.Do(context.Background(),
		gitlab.WithPath("/ok"),
		gitlab.WithResponseHandler(gitlab.JsonResponseHandler(&resp)),
	)
	require.NoError(t, err)
	assert.EqualValues(t, 42, resp)
}

func TestJsonResponseHandler_Cancellation(t *testing.T) {
	ctxClient, cancelClient := context.WithCancel(context.Background())
	defer cancelClient()
	cancelServer := make(chan struct{})
	c := mock_gitlab.SetupClient(t, "/cancel", func(w http.ResponseWriter, r *http.Request) {
		testhelpers.AssertGetJsonRequest(t, r)
		cancelClient() // unblock client
		<-cancelServer // wait for client to get the error and unblock server
	})

	var resp interface{}

	err := c.Do(ctxClient,
		gitlab.WithPath("/cancel"),
		gitlab.WithResponseHandler(gitlab.JsonResponseHandler(&resp)),
	)
	close(cancelServer) // unblock server
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}

func TestNoContentResponseHandler_Errors(t *testing.T) {
	tests := map[int]func(error) bool{
		http.StatusForbidden:    gitlab.IsForbidden,
		http.StatusUnauthorized: gitlab.IsUnauthorized,
		http.StatusNotFound:     gitlab.IsNotFound,
	}
	for statusCode, f := range tests {
		t.Run(strconv.Itoa(statusCode), func(t *testing.T) {
			c := mock_gitlab.SetupClient(t, "/bla", func(w http.ResponseWriter, r *http.Request) {
				assertNoContentRequest(t, r)
				w.WriteHeader(statusCode)
			})

			err := c.Do(context.Background(),
				gitlab.WithPath("/bla"),
				gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
			)
			require.Error(t, err)
			assert.True(t, f(err))
			assert.True(t, errHasPath(err, "/bla"))
		})
	}
}

func TestNoContentResponseHandler_Unauthorized(t *testing.T) {
	c := mock_gitlab.SetupClient(t, "/unauthorized", func(w http.ResponseWriter, r *http.Request) {
		assertNoContentRequest(t, r)
		w.WriteHeader(http.StatusUnauthorized)
	})

	err := c.Do(context.Background(),
		gitlab.WithPath("/unauthorized"),
		gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
	)
	require.Error(t, err)
	assert.True(t, gitlab.IsUnauthorized(err))
	assert.True(t, errHasPath(err, "/unauthorized"))
}

func TestNoContentResponseHandler_HappyPath(t *testing.T) {
	c := mock_gitlab.SetupClient(t, "/ok", func(w http.ResponseWriter, r *http.Request) {
		assertNoContentRequest(t, r)
		testhelpers.RespondWithJSON(t, w, 42)
	})

	err := c.Do(context.Background(),
		gitlab.WithPath("/ok"),
		gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
	)
	require.NoError(t, err)
}

func TestNoContentResponseHandler_Cancellation(t *testing.T) {
	ctxClient, cancelClient := context.WithCancel(context.Background())
	defer cancelClient()
	cancelServer := make(chan struct{})
	c := mock_gitlab.SetupClient(t, "/cancel", func(w http.ResponseWriter, r *http.Request) {
		assertNoContentRequest(t, r)
		cancelClient() // unblock client
		<-cancelServer // wait for client to get the error and unblock server
	})

	err := c.Do(ctxClient,
		gitlab.WithPath("/cancel"),
		gitlab.WithResponseHandler(gitlab.NoContentResponseHandler()),
	)
	close(cancelServer) // unblock server
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
}

func assertNoContentRequest(t *testing.T, r *http.Request) {
	testhelpers.AssertRequestMethod(t, r, http.MethodGet)
	assert.Empty(t, r.Header.Values(httpz.AcceptHeader))
}

func errHasPath(err error, path string) bool {
	var e *gitlab.ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.Path == path
}
