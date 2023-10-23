package gitlab

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/pluralsh/kuberentes-agent/internal/tool/httpz"
	"github.com/stretchr/testify/require"
)

func TestDefaultErrorHandlerWithReason_DiscardUnknownFields(t *testing.T) {
	// GIVEN
	url, _ := url.Parse("https://gitlab.example.com/api/v4/anything")
	resp := &http.Response{
		Request:    &http.Request{URL: url},
		StatusCode: http.StatusUnauthorized,
		Header:     map[string][]string{httpz.ContentTypeHeader: {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"message": "anything", "ignored": "ignored"}`)),
	}

	// WHEN
	err := defaultErrorHandlerWithReason(resp)

	// THEN
	require.EqualError(t, err, "HTTP status code: 401 for path /api/v4/anything with reason anything")
}
