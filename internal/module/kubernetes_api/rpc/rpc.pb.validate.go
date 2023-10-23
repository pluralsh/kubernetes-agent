// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: internal/module/kubernetes_api/rpc/rpc.proto

package rpc

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

// Validate checks the field values on HeaderExtra with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *HeaderExtra) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HeaderExtra with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in HeaderExtraMultiError, or
// nil if none found.
func (m *HeaderExtra) ValidateAll() error {
	return m.validate(true)
}

func (m *HeaderExtra) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetImpConfig()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, HeaderExtraValidationError{
					field:  "ImpConfig",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, HeaderExtraValidationError{
					field:  "ImpConfig",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetImpConfig()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HeaderExtraValidationError{
				field:  "ImpConfig",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return HeaderExtraMultiError(errors)
	}

	return nil
}

// HeaderExtraMultiError is an error wrapping multiple validation errors
// returned by HeaderExtra.ValidateAll() if the designated constraints aren't met.
type HeaderExtraMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HeaderExtraMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HeaderExtraMultiError) AllErrors() []error { return m }

// HeaderExtraValidationError is the validation error returned by
// HeaderExtra.Validate if the designated constraints aren't met.
type HeaderExtraValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HeaderExtraValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HeaderExtraValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HeaderExtraValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HeaderExtraValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HeaderExtraValidationError) ErrorName() string { return "HeaderExtraValidationError" }

// Error satisfies the builtin error interface
func (e HeaderExtraValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHeaderExtra.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HeaderExtraValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HeaderExtraValidationError{}

// Validate checks the field values on ImpersonationConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ImpersonationConfig) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ImpersonationConfig with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ImpersonationConfigMultiError, or nil if none found.
func (m *ImpersonationConfig) ValidateAll() error {
	return m.validate(true)
}

func (m *ImpersonationConfig) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Username

	// no validation rules for Uid

	for idx, item := range m.GetExtra() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ImpersonationConfigValidationError{
						field:  fmt.Sprintf("Extra[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ImpersonationConfigValidationError{
						field:  fmt.Sprintf("Extra[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ImpersonationConfigValidationError{
					field:  fmt.Sprintf("Extra[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ImpersonationConfigMultiError(errors)
	}

	return nil
}

// ImpersonationConfigMultiError is an error wrapping multiple validation
// errors returned by ImpersonationConfig.ValidateAll() if the designated
// constraints aren't met.
type ImpersonationConfigMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ImpersonationConfigMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ImpersonationConfigMultiError) AllErrors() []error { return m }

// ImpersonationConfigValidationError is the validation error returned by
// ImpersonationConfig.Validate if the designated constraints aren't met.
type ImpersonationConfigValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ImpersonationConfigValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ImpersonationConfigValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ImpersonationConfigValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ImpersonationConfigValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ImpersonationConfigValidationError) ErrorName() string {
	return "ImpersonationConfigValidationError"
}

// Error satisfies the builtin error interface
func (e ImpersonationConfigValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sImpersonationConfig.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ImpersonationConfigValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ImpersonationConfigValidationError{}

// Validate checks the field values on ExtraKeyVal with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ExtraKeyVal) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ExtraKeyVal with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ExtraKeyValMultiError, or
// nil if none found.
func (m *ExtraKeyVal) ValidateAll() error {
	return m.validate(true)
}

func (m *ExtraKeyVal) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Key

	if len(errors) > 0 {
		return ExtraKeyValMultiError(errors)
	}

	return nil
}

// ExtraKeyValMultiError is an error wrapping multiple validation errors
// returned by ExtraKeyVal.ValidateAll() if the designated constraints aren't met.
type ExtraKeyValMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ExtraKeyValMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ExtraKeyValMultiError) AllErrors() []error { return m }

// ExtraKeyValValidationError is the validation error returned by
// ExtraKeyVal.Validate if the designated constraints aren't met.
type ExtraKeyValValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ExtraKeyValValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ExtraKeyValValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ExtraKeyValValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ExtraKeyValValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ExtraKeyValValidationError) ErrorName() string { return "ExtraKeyValValidationError" }

// Error satisfies the builtin error interface
func (e ExtraKeyValValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sExtraKeyVal.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ExtraKeyValValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ExtraKeyValValidationError{}