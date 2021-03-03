// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: cmd/kas/kasapp/kasapp.proto

package kasapp

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
var _kasapp_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on StartStreaming with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *StartStreaming) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// StartStreamingValidationError is the validation error returned by
// StartStreaming.Validate if the designated constraints aren't met.
type StartStreamingValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StartStreamingValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StartStreamingValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StartStreamingValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StartStreamingValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StartStreamingValidationError) ErrorName() string { return "StartStreamingValidationError" }

// Error satisfies the builtin error interface
func (e StartStreamingValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStartStreaming.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StartStreamingValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StartStreamingValidationError{}

// Validate checks the field values on GatewayKasResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GatewayKasResponse) Validate() error {
	if m == nil {
		return nil
	}

	switch m.Msg.(type) {

	case *GatewayKasResponse_TunnelReady_:

		if v, ok := interface{}(m.GetTunnelReady()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponseValidationError{
					field:  "TunnelReady",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *GatewayKasResponse_Header_:

		if v, ok := interface{}(m.GetHeader()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponseValidationError{
					field:  "Header",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *GatewayKasResponse_Message_:

		if v, ok := interface{}(m.GetMessage()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponseValidationError{
					field:  "Message",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *GatewayKasResponse_Trailer_:

		if v, ok := interface{}(m.GetTrailer()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponseValidationError{
					field:  "Trailer",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *GatewayKasResponse_Error_:

		if v, ok := interface{}(m.GetError()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponseValidationError{
					field:  "Error",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		return GatewayKasResponseValidationError{
			field:  "Msg",
			reason: "value is required",
		}

	}

	return nil
}

// GatewayKasResponseValidationError is the validation error returned by
// GatewayKasResponse.Validate if the designated constraints aren't met.
type GatewayKasResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GatewayKasResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GatewayKasResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GatewayKasResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GatewayKasResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GatewayKasResponseValidationError) ErrorName() string {
	return "GatewayKasResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GatewayKasResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGatewayKasResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GatewayKasResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GatewayKasResponseValidationError{}

// Validate checks the field values on GatewayKasResponse_TunnelReady with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GatewayKasResponse_TunnelReady) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// GatewayKasResponse_TunnelReadyValidationError is the validation error
// returned by GatewayKasResponse_TunnelReady.Validate if the designated
// constraints aren't met.
type GatewayKasResponse_TunnelReadyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GatewayKasResponse_TunnelReadyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GatewayKasResponse_TunnelReadyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GatewayKasResponse_TunnelReadyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GatewayKasResponse_TunnelReadyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GatewayKasResponse_TunnelReadyValidationError) ErrorName() string {
	return "GatewayKasResponse_TunnelReadyValidationError"
}

// Error satisfies the builtin error interface
func (e GatewayKasResponse_TunnelReadyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGatewayKasResponse_TunnelReady.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GatewayKasResponse_TunnelReadyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GatewayKasResponse_TunnelReadyValidationError{}

// Validate checks the field values on GatewayKasResponse_Header with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GatewayKasResponse_Header) Validate() error {
	if m == nil {
		return nil
	}

	for key, val := range m.GetMeta() {
		_ = val

		// no validation rules for Meta[key]

		if v, ok := interface{}(val).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponse_HeaderValidationError{
					field:  fmt.Sprintf("Meta[%v]", key),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// GatewayKasResponse_HeaderValidationError is the validation error returned by
// GatewayKasResponse_Header.Validate if the designated constraints aren't met.
type GatewayKasResponse_HeaderValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GatewayKasResponse_HeaderValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GatewayKasResponse_HeaderValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GatewayKasResponse_HeaderValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GatewayKasResponse_HeaderValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GatewayKasResponse_HeaderValidationError) ErrorName() string {
	return "GatewayKasResponse_HeaderValidationError"
}

// Error satisfies the builtin error interface
func (e GatewayKasResponse_HeaderValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGatewayKasResponse_Header.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GatewayKasResponse_HeaderValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GatewayKasResponse_HeaderValidationError{}

// Validate checks the field values on GatewayKasResponse_Message with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GatewayKasResponse_Message) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Data

	return nil
}

// GatewayKasResponse_MessageValidationError is the validation error returned
// by GatewayKasResponse_Message.Validate if the designated constraints aren't met.
type GatewayKasResponse_MessageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GatewayKasResponse_MessageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GatewayKasResponse_MessageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GatewayKasResponse_MessageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GatewayKasResponse_MessageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GatewayKasResponse_MessageValidationError) ErrorName() string {
	return "GatewayKasResponse_MessageValidationError"
}

// Error satisfies the builtin error interface
func (e GatewayKasResponse_MessageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGatewayKasResponse_Message.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GatewayKasResponse_MessageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GatewayKasResponse_MessageValidationError{}

// Validate checks the field values on GatewayKasResponse_Trailer with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GatewayKasResponse_Trailer) Validate() error {
	if m == nil {
		return nil
	}

	for key, val := range m.GetMeta() {
		_ = val

		// no validation rules for Meta[key]

		if v, ok := interface{}(val).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GatewayKasResponse_TrailerValidationError{
					field:  fmt.Sprintf("Meta[%v]", key),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// GatewayKasResponse_TrailerValidationError is the validation error returned
// by GatewayKasResponse_Trailer.Validate if the designated constraints aren't met.
type GatewayKasResponse_TrailerValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GatewayKasResponse_TrailerValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GatewayKasResponse_TrailerValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GatewayKasResponse_TrailerValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GatewayKasResponse_TrailerValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GatewayKasResponse_TrailerValidationError) ErrorName() string {
	return "GatewayKasResponse_TrailerValidationError"
}

// Error satisfies the builtin error interface
func (e GatewayKasResponse_TrailerValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGatewayKasResponse_Trailer.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GatewayKasResponse_TrailerValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GatewayKasResponse_TrailerValidationError{}

// Validate checks the field values on GatewayKasResponse_Error with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GatewayKasResponse_Error) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetStatus() == nil {
		return GatewayKasResponse_ErrorValidationError{
			field:  "Status",
			reason: "value is required",
		}
	}

	if v, ok := interface{}(m.GetStatus()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GatewayKasResponse_ErrorValidationError{
				field:  "Status",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// GatewayKasResponse_ErrorValidationError is the validation error returned by
// GatewayKasResponse_Error.Validate if the designated constraints aren't met.
type GatewayKasResponse_ErrorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GatewayKasResponse_ErrorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GatewayKasResponse_ErrorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GatewayKasResponse_ErrorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GatewayKasResponse_ErrorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GatewayKasResponse_ErrorValidationError) ErrorName() string {
	return "GatewayKasResponse_ErrorValidationError"
}

// Error satisfies the builtin error interface
func (e GatewayKasResponse_ErrorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGatewayKasResponse_Error.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GatewayKasResponse_ErrorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GatewayKasResponse_ErrorValidationError{}
