// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: internal/module/reverse_tunnel/info/info.proto

package info

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

	"github.com/golang/protobuf/ptypes"
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
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _info_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Method with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Method) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		return MethodValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// MethodValidationError is the validation error returned by Method.Validate if
// the designated constraints aren't met.
type MethodValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MethodValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MethodValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MethodValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MethodValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MethodValidationError) ErrorName() string { return "MethodValidationError" }

// Error satisfies the builtin error interface
func (e MethodValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMethod.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MethodValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MethodValidationError{}

// Validate checks the field values on Service with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Service) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		return ServiceValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
	}

	for idx, item := range m.GetMethods() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ServiceValidationError{
					field:  fmt.Sprintf("Methods[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ServiceValidationError is the validation error returned by Service.Validate
// if the designated constraints aren't met.
type ServiceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ServiceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ServiceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ServiceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ServiceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ServiceValidationError) ErrorName() string { return "ServiceValidationError" }

// Error satisfies the builtin error interface
func (e ServiceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sService.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ServiceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ServiceValidationError{}

// Validate checks the field values on AgentDescriptor with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *AgentDescriptor) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetServices() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AgentDescriptorValidationError{
					field:  fmt.Sprintf("Services[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// AgentDescriptorValidationError is the validation error returned by
// AgentDescriptor.Validate if the designated constraints aren't met.
type AgentDescriptorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AgentDescriptorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AgentDescriptorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AgentDescriptorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AgentDescriptorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AgentDescriptorValidationError) ErrorName() string { return "AgentDescriptorValidationError" }

// Error satisfies the builtin error interface
func (e AgentDescriptorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAgentDescriptor.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AgentDescriptorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AgentDescriptorValidationError{}