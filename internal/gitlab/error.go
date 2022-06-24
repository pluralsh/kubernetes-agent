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
	p := x.Path
	if p == "" {
		p = "<unknown>"
	}
	return fmt.Sprintf("HTTP status code: %d for path %s", x.StatusCode, p)
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
