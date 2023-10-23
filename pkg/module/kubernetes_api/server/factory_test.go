package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetAuthorizedProxyUserCacheKeyFunc_AllFieldsUsed(t *testing.T) {
	keyFunc := getAuthorizedProxyUserCacheKey("any-prefix")

	redisKeys := map[string]struct{}{}
	redisKeys[keyFunc(proxyUserCacheKey{
		agentId:    1,
		accessType: "any",
		accessKey:  "any",
		csrfToken:  "any",
	})] = struct{}{}
	redisKeys[keyFunc(proxyUserCacheKey{
		accessType: "any",
		accessKey:  "any",
	})] = struct{}{}
	redisKeys[keyFunc(proxyUserCacheKey{
		agentId:   1,
		accessKey: "any",
		csrfToken: "any",
	})] = struct{}{}
	redisKeys[keyFunc(proxyUserCacheKey{
		agentId:    1,
		accessType: "any",
		csrfToken:  "any",
	})] = struct{}{}

	assert.Equal(t, 4, len(redisKeys))
}
