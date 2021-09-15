package prototool

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	_ json.Marshaler   = JsonBox{}
	_ json.Unmarshaler = (*JsonBox)(nil)
)

// Test round trip using stdlib json package.
func TestJsonBox_RoundTrip(t *testing.T) {
	val := &HttpRequest{
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
	}
	tests := []struct {
		input    interface{}
		output   interface{}
		expected interface{} // input is used if not set
	}{
		{
			input: JsonBox{
				Message: val,
			},
			output: &JsonBox{
				Message: &HttpRequest{},
			},
			expected: &JsonBox{
				Message: val,
			},
		},
		{
			input: &JsonBox{
				Message: val,
			},
			output: &JsonBox{
				Message: &HttpRequest{},
			},
		},
		{
			input: &embeddedBox{ // pointer
				A: JsonBox{
					Message: val,
				},
				B: &JsonBox{
					Message: val,
				},
			},
			output: &embeddedBox{
				A: JsonBox{
					Message: &HttpRequest{},
				},
				B: &JsonBox{
					Message: &HttpRequest{},
				},
			},
		},
		{
			input: embeddedBox{ // not a pointer
				A: JsonBox{
					Message: val,
				},
				B: &JsonBox{
					Message: val,
				},
			},
			output: &embeddedBox{
				A: JsonBox{
					Message: &HttpRequest{},
				},
				B: &JsonBox{
					Message: &HttpRequest{},
				},
			},
			expected: &embeddedBox{
				A: JsonBox{
					Message: val,
				},
				B: &JsonBox{
					Message: val,
				},
			},
		},
		{
			input: &embeddedBox{},
			output: &embeddedBox{
				A: JsonBox{
					Message: &HttpRequest{},
				},
				B: &JsonBox{
					Message: &HttpRequest{},
				},
			},
			expected: &embeddedBox{
				A: JsonBox{
					Message: &HttpRequest{},
				},
			},
		},
	}
	for i, tc := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			data, err := json.Marshal(tc.input)
			require.NoError(t, err)

			err = json.Unmarshal(data, tc.output)
			require.NoError(t, err)

			var expected interface{}
			if tc.expected != nil {
				expected = tc.expected
			} else {
				expected = tc.input
			}
			assert.Empty(t, cmp.Diff(expected, tc.output, protocmp.Transform()))
		})
	}
}

func TestJsonBox_IgnoreUnknownFields(t *testing.T) {
	var (
		in  = []byte(`{"method":"POST", "some_field":"Ha-ha! I'll break your code!"}`)
		out = JsonBox{
			Message: &HttpRequest{},
		}
		expected = &HttpRequest{
			Method: "POST",
		}
	)
	err := json.Unmarshal(in, &out)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(expected, out.Message, protocmp.Transform()))
}

type embeddedBox struct {
	A JsonBox
	B *JsonBox
}
