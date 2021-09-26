package api

import (
	"errors"
	"net/http"

	"gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v14/internal/gitlab"
)

// IsCacheableError checks if an error is cacheable.
func IsCacheableError(err error) bool {
	var e gitlab.ClientError
	if !errors.As(err, &e) {
		return false // not a client error, probably a network error
	}
	switch e.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		return true
	default:
		return false
	}
}
