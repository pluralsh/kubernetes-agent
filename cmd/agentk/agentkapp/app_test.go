package agentkapp

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrpcHostWithPort(t *testing.T) {
	tests := []struct {
		inUrl               string
		expectedOutHostPort string
	}{
		{
			inUrl:               "grpc://test.test",
			expectedOutHostPort: "test.test:80",
		},
		{
			inUrl:               "grpcs://test.test",
			expectedOutHostPort: "test.test:443",
		},
		{
			inUrl:               "grpc://test.test:123",
			expectedOutHostPort: "test.test:123",
		},
		{
			inUrl:               "grpcs://test.test:123",
			expectedOutHostPort: "test.test:123",
		},
		{
			inUrl:               "grpc://1.2.3.4",
			expectedOutHostPort: "1.2.3.4:80",
		},
		{
			inUrl:               "grpcs://1.2.3.4",
			expectedOutHostPort: "1.2.3.4:443",
		},
		{
			inUrl:               "grpc://[123::123]:123",
			expectedOutHostPort: "[123::123]:123",
		},
		{
			inUrl:               "grpcs://[123::123]:123",
			expectedOutHostPort: "[123::123]:123",
		},
	}
	for _, test := range tests {
		t.Run(test.inUrl, func(t *testing.T) {
			u, err := url.Parse(test.inUrl)
			require.NoError(t, err)
			hostAndPort := grpcHostWithPort(u)
			assert.Equal(t, test.expectedOutHostPort, hostAndPort)
		})
	}
}
