// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: internal/module/reverse_tunnel/tracker/tracker.proto

package tracker

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
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
)

// Validate checks the field values on TunnelInfo with the rules defined in the
// proto definition for this message. If any rules are violated, an error is
// returned. When asked to return all errors, validation continues after first
// violation, and the result is a list of violation errors wrapped in
// TunnelInfoMultiError, or nil if none found. Otherwise, only the first error
// is returned, if any.
func (m *TunnelInfo) Validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetAgentDescriptor() == nil {
		err := TunnelInfoValidationError{
			field:  "AgentDescriptor",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if v, ok := interface{}(m.GetAgentDescriptor()).(interface{ Validate(bool) error }); ok {
		if err := v.Validate(all); err != nil {
			err = TunnelInfoValidationError{
				field:  "AgentDescriptor",
				reason: "embedded message failed validation",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}
	}

	// no validation rules for ConnectionId

	// no validation rules for AgentId

	if !_TunnelInfo_KasUrl_Pattern.MatchString(m.GetKasUrl()) {
		err := TunnelInfoValidationError{
			field:  "KasUrl",
			reason: "value does not match regex pattern \"(?:^$|^grpcs?://)\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return TunnelInfoMultiError(errors)
	}
	return nil
}

// TunnelInfoMultiError is an error wrapping multiple validation errors
// returned by TunnelInfo.Validate(true) if the designated constraints aren't met.
type TunnelInfoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TunnelInfoMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TunnelInfoMultiError) AllErrors() []error { return m }

// TunnelInfoValidationError is the validation error returned by
// TunnelInfo.Validate if the designated constraints aren't met.
type TunnelInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TunnelInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TunnelInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TunnelInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TunnelInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TunnelInfoValidationError) ErrorName() string { return "TunnelInfoValidationError" }

// Error satisfies the builtin error interface
func (e TunnelInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTunnelInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TunnelInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TunnelInfoValidationError{}

var _TunnelInfo_KasUrl_Pattern = regexp.MustCompile("(?:^$|^grpcs?://)")
