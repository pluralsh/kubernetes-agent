package prototool

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	_ json.Marshaler   = (*JsonBox)(nil)
	_ json.Unmarshaler = (*JsonBox)(nil)
)

// Test round trip using stdlib json package.
func TestJsonBox_RoundTrip(t *testing.T) {
	source := &JsonBox{
		Message: &HttpRequest{
			Method: "POST",
			Header: map[string]*Values{
				"Req-Header": {
					Value: []string{"x1", "x2"},
				},
			},
			UrlPath: "adsad",
			Query: map[string]*Values{
				"asdsad": {
					Value: []string{"asdasds"},
				},
			},
		},
	}

	data, err := json.Marshal(source)
	require.NoError(t, err)

	actual := &JsonBox{
		Message: &HttpRequest{},
	}
	err = json.Unmarshal(data, actual)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(source.Message, actual.Message, protocmp.Transform()))
}
