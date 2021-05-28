// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pkg/agentcfg/agentcfg.proto

package agentcfg

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

// Validate checks the field values on PathCF with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *PathCF) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetGlob()) < 1 {
		return PathCFValidationError{
			field:  "Glob",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// PathCFValidationError is the validation error returned by PathCF.Validate if
// the designated constraints aren't met.
type PathCFValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PathCFValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PathCFValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PathCFValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PathCFValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PathCFValidationError) ErrorName() string { return "PathCFValidationError" }

// Error satisfies the builtin error interface
func (e PathCFValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPathCF.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PathCFValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PathCFValidationError{}

// Validate checks the field values on ManifestProjectCF with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ManifestProjectCF) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetId()) < 1 {
		return ManifestProjectCFValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
	}

	// no validation rules for DefaultNamespace

	for idx, item := range m.GetPaths() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ManifestProjectCFValidationError{
					field:  fmt.Sprintf("Paths[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if v, ok := interface{}(m.GetReconcileTimeout()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ManifestProjectCFValidationError{
				field:  "ReconcileTimeout",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if _, ok := _ManifestProjectCF_DryRunStrategy_InLookup[m.GetDryRunStrategy()]; !ok {
		return ManifestProjectCFValidationError{
			field:  "DryRunStrategy",
			reason: "value must be in list [ none client server]",
		}
	}

	// no validation rules for NoPrune

	if v, ok := interface{}(m.GetPruneTimeout()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ManifestProjectCFValidationError{
				field:  "PruneTimeout",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if _, ok := _ManifestProjectCF_PrunePropagationPolicy_InLookup[m.GetPrunePropagationPolicy()]; !ok {
		return ManifestProjectCFValidationError{
			field:  "PrunePropagationPolicy",
			reason: "value must be in list [ orphan background foreground]",
		}
	}

	if _, ok := _ManifestProjectCF_InventoryPolicy_InLookup[m.GetInventoryPolicy()]; !ok {
		return ManifestProjectCFValidationError{
			field:  "InventoryPolicy",
			reason: "value must be in list [ must_match adopt_if_no_inventory adopt_all]",
		}
	}

	return nil
}

// ManifestProjectCFValidationError is the validation error returned by
// ManifestProjectCF.Validate if the designated constraints aren't met.
type ManifestProjectCFValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ManifestProjectCFValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ManifestProjectCFValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ManifestProjectCFValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ManifestProjectCFValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ManifestProjectCFValidationError) ErrorName() string {
	return "ManifestProjectCFValidationError"
}

// Error satisfies the builtin error interface
func (e ManifestProjectCFValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sManifestProjectCF.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ManifestProjectCFValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ManifestProjectCFValidationError{}

var _ManifestProjectCF_DryRunStrategy_InLookup = map[string]struct{}{
	"":       {},
	"none":   {},
	"client": {},
	"server": {},
}

var _ManifestProjectCF_PrunePropagationPolicy_InLookup = map[string]struct{}{
	"":           {},
	"orphan":     {},
	"background": {},
	"foreground": {},
}

var _ManifestProjectCF_InventoryPolicy_InLookup = map[string]struct{}{
	"":                      {},
	"must_match":            {},
	"adopt_if_no_inventory": {},
	"adopt_all":             {},
}

// Validate checks the field values on GitopsCF with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *GitopsCF) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetManifestProjects() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GitopsCFValidationError{
					field:  fmt.Sprintf("ManifestProjects[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// GitopsCFValidationError is the validation error returned by
// GitopsCF.Validate if the designated constraints aren't met.
type GitopsCFValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GitopsCFValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GitopsCFValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GitopsCFValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GitopsCFValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GitopsCFValidationError) ErrorName() string { return "GitopsCFValidationError" }

// Error satisfies the builtin error interface
func (e GitopsCFValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGitopsCF.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GitopsCFValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GitopsCFValidationError{}

// Validate checks the field values on ObservabilityCF with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ObservabilityCF) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetLogging()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ObservabilityCFValidationError{
				field:  "Logging",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// ObservabilityCFValidationError is the validation error returned by
// ObservabilityCF.Validate if the designated constraints aren't met.
type ObservabilityCFValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ObservabilityCFValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ObservabilityCFValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ObservabilityCFValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ObservabilityCFValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ObservabilityCFValidationError) ErrorName() string { return "ObservabilityCFValidationError" }

// Error satisfies the builtin error interface
func (e ObservabilityCFValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sObservabilityCF.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ObservabilityCFValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ObservabilityCFValidationError{}

// Validate checks the field values on LoggingCF with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *LoggingCF) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Level

	return nil
}

// LoggingCFValidationError is the validation error returned by
// LoggingCF.Validate if the designated constraints aren't met.
type LoggingCFValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoggingCFValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoggingCFValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoggingCFValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoggingCFValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoggingCFValidationError) ErrorName() string { return "LoggingCFValidationError" }

// Error satisfies the builtin error interface
func (e LoggingCFValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoggingCF.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoggingCFValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoggingCFValidationError{}

// Validate checks the field values on CiliumCF with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *CiliumCF) Validate() error {
	if m == nil {
		return nil
	}

	if utf8.RuneCountInString(m.GetHubbleRelayAddress()) < 1 {
		return CiliumCFValidationError{
			field:  "HubbleRelayAddress",
			reason: "value length must be at least 1 runes",
		}
	}

	return nil
}

// CiliumCFValidationError is the validation error returned by
// CiliumCF.Validate if the designated constraints aren't met.
type CiliumCFValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CiliumCFValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CiliumCFValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CiliumCFValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CiliumCFValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CiliumCFValidationError) ErrorName() string { return "CiliumCFValidationError" }

// Error satisfies the builtin error interface
func (e CiliumCFValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCiliumCF.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CiliumCFValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CiliumCFValidationError{}

// Validate checks the field values on ConfigurationFile with the rules defined
// in the proto definition for this message. If any rules are violated, an
// error is returned.
func (m *ConfigurationFile) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetGitops()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ConfigurationFileValidationError{
				field:  "Gitops",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetObservability()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ConfigurationFileValidationError{
				field:  "Observability",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetCilium()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ConfigurationFileValidationError{
				field:  "Cilium",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// ConfigurationFileValidationError is the validation error returned by
// ConfigurationFile.Validate if the designated constraints aren't met.
type ConfigurationFileValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ConfigurationFileValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ConfigurationFileValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ConfigurationFileValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ConfigurationFileValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ConfigurationFileValidationError) ErrorName() string {
	return "ConfigurationFileValidationError"
}

// Error satisfies the builtin error interface
func (e ConfigurationFileValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sConfigurationFile.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ConfigurationFileValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ConfigurationFileValidationError{}

// Validate checks the field values on AgentConfiguration with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *AgentConfiguration) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetGitops()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AgentConfigurationValidationError{
				field:  "Gitops",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetObservability()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AgentConfigurationValidationError{
				field:  "Observability",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetCilium()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AgentConfigurationValidationError{
				field:  "Cilium",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for AgentId

	// no validation rules for ProjectId

	return nil
}

// AgentConfigurationValidationError is the validation error returned by
// AgentConfiguration.Validate if the designated constraints aren't met.
type AgentConfigurationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AgentConfigurationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AgentConfigurationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AgentConfigurationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AgentConfigurationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AgentConfigurationValidationError) ErrorName() string {
	return "AgentConfigurationValidationError"
}

// Error satisfies the builtin error interface
func (e AgentConfigurationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAgentConfiguration.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AgentConfigurationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AgentConfigurationValidationError{}
