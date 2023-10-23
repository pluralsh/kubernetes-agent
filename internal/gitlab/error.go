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
	p := x.GetPath()
	if p == "" {
		p = "<unknown>"
	}
	r := x.GetReason()
	if r != "" {
		return fmt.Sprintf("HTTP status code: %d for path %s with reason %s", x.GetStatusCode(), p, r)
	}
	return fmt.Sprintf("HTTP status code: %d for path %s", x.GetStatusCode(), p)
}

func IsForbidden(err error) bool {
	var e *ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.GetStatusCode() == http.StatusForbidden
}

func IsUnauthorized(err error) bool {
	var e *ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.GetStatusCode() == http.StatusUnauthorized
}

func IsNotFound(err error) bool {
	var e *ClientError
	if !errors.As(err, &e) {
		return false
	}
	return e.GetStatusCode() == http.StatusNotFound
}
