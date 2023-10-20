// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: internal/gitlab/gitlab.proto

package gitlab

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on ClientError with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ClientError) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClientError with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ClientErrorMultiError, or
// nil if none found.
func (m *ClientError) ValidateAll() error {
	return m.validate(true)
}

func (m *ClientError) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for StatusCode

	// no validation rules for Path

	// no validation rules for Reason

	if len(errors) > 0 {
		return ClientErrorMultiError(errors)
	}

	return nil
}

// ClientErrorMultiError is an error wrapping multiple validation errors
// returned by ClientError.ValidateAll() if the designated constraints aren't met.
type ClientErrorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClientErrorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClientErrorMultiError) AllErrors() []error { return m }

// ClientErrorValidationError is the validation error returned by
// ClientError.Validate if the designated constraints aren't met.
type ClientErrorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClientErrorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClientErrorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClientErrorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClientErrorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClientErrorValidationError) ErrorName() string { return "ClientErrorValidationError" }

// Error satisfies the builtin error interface
func (e ClientErrorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClientError.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClientErrorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClientErrorValidationError{}

// Validate checks the field values on DefaultApiError with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DefaultApiError) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DefaultApiError with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DefaultApiErrorMultiError, or nil if none found.
func (m *DefaultApiError) ValidateAll() error {
	return m.validate(true)
}

func (m *DefaultApiError) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Message

	if len(errors) > 0 {
		return DefaultApiErrorMultiError(errors)
	}

	return nil
}

// DefaultApiErrorMultiError is an error wrapping multiple validation errors
// returned by DefaultApiError.ValidateAll() if the designated constraints
// aren't met.
type DefaultApiErrorMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DefaultApiErrorMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DefaultApiErrorMultiError) AllErrors() []error { return m }

// DefaultApiErrorValidationError is the validation error returned by
// DefaultApiError.Validate if the designated constraints aren't met.
type DefaultApiErrorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DefaultApiErrorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DefaultApiErrorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DefaultApiErrorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DefaultApiErrorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DefaultApiErrorValidationError) ErrorName() string { return "DefaultApiErrorValidationError" }

// Error satisfies the builtin error interface
func (e DefaultApiErrorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDefaultApiError.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DefaultApiErrorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DefaultApiErrorValidationError{}
