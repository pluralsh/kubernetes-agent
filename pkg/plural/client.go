package plural

import (
	"context"
	"net/http"

	console "github.com/pluralsh/console-client-go"
)

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Token "+t.token)
	return t.wrapped.RoundTrip(req)
}

type Client struct {
	ctx           context.Context
	consoleClient *console.Client
}

func New(url, token string) *Client {
	httpClient := http.Client{
		Transport: &authedTransport{
			token:   token,
			wrapped: http.DefaultTransport,
		},
	}

	return &Client{
		consoleClient: console.NewClient(&httpClient, url),
		ctx:           context.Background(),
	}
}
