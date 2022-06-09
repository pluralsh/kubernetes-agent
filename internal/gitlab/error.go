package gitlab

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	_ error = (*ClientError)(nil)
)

func (x *ClientError) Error() string {
	return fmt.Sprintf("HTTP status code: %d", x.StatusCode)
}

func IsForbidden(err error) bool {
	var e *ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.StatusCode == http.StatusForbidden
}

func IsUnauthorized(err error) bool {
	var e *ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.StatusCode == http.StatusUnauthorized
}

func IsNotFound(err error) bool {
	var e *ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.StatusCode == http.StatusNotFound
}
