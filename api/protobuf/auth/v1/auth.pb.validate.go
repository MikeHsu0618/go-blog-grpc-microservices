// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: auth.proto

package v1

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

// Validate checks the field values on GenerateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GenerateTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GenerateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GenerateTokenRequestMultiError, or nil if none found.
func (m *GenerateTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GenerateTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	if len(errors) > 0 {
		return GenerateTokenRequestMultiError(errors)
	}

	return nil
}

// GenerateTokenRequestMultiError is an error wrapping multiple validation
// errors returned by GenerateTokenRequest.ValidateAll() if the designated
// constraints aren't met.
type GenerateTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GenerateTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GenerateTokenRequestMultiError) AllErrors() []error { return m }

// GenerateTokenRequestValidationError is the validation error returned by
// GenerateTokenRequest.Validate if the designated constraints aren't met.
type GenerateTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateTokenRequestValidationError) ErrorName() string {
	return "GenerateTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateTokenRequestValidationError{}

// Validate checks the field values on GenerateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GenerateTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GenerateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GenerateTokenResponseMultiError, or nil if none found.
func (m *GenerateTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GenerateTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if len(errors) > 0 {
		return GenerateTokenResponseMultiError(errors)
	}

	return nil
}

// GenerateTokenResponseMultiError is an error wrapping multiple validation
// errors returned by GenerateTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type GenerateTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GenerateTokenResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GenerateTokenResponseMultiError) AllErrors() []error { return m }

// GenerateTokenResponseValidationError is the validation error returned by
// GenerateTokenResponse.Validate if the designated constraints aren't met.
type GenerateTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GenerateTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GenerateTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GenerateTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GenerateTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GenerateTokenResponseValidationError) ErrorName() string {
	return "GenerateTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GenerateTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGenerateTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GenerateTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GenerateTokenResponseValidationError{}

// Validate checks the field values on ValidateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ValidateTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ValidateTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ValidateTokenRequestMultiError, or nil if none found.
func (m *ValidateTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ValidateTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if len(errors) > 0 {
		return ValidateTokenRequestMultiError(errors)
	}

	return nil
}

// ValidateTokenRequestMultiError is an error wrapping multiple validation
// errors returned by ValidateTokenRequest.ValidateAll() if the designated
// constraints aren't met.
type ValidateTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValidateTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValidateTokenRequestMultiError) AllErrors() []error { return m }

// ValidateTokenRequestValidationError is the validation error returned by
// ValidateTokenRequest.Validate if the designated constraints aren't met.
type ValidateTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValidateTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValidateTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValidateTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValidateTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValidateTokenRequestValidationError) ErrorName() string {
	return "ValidateTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ValidateTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValidateTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValidateTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValidateTokenRequestValidationError{}

// Validate checks the field values on ValidateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ValidateTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ValidateTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ValidateTokenResponseMultiError, or nil if none found.
func (m *ValidateTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ValidateTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Valid

	if len(errors) > 0 {
		return ValidateTokenResponseMultiError(errors)
	}

	return nil
}

// ValidateTokenResponseMultiError is an error wrapping multiple validation
// errors returned by ValidateTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type ValidateTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ValidateTokenResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ValidateTokenResponseMultiError) AllErrors() []error { return m }

// ValidateTokenResponseValidationError is the validation error returned by
// ValidateTokenResponse.Validate if the designated constraints aren't met.
type ValidateTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ValidateTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ValidateTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ValidateTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ValidateTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ValidateTokenResponseValidationError) ErrorName() string {
	return "ValidateTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ValidateTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sValidateTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ValidateTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ValidateTokenResponseValidationError{}

// Validate checks the field values on RefreshTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RefreshTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RefreshTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RefreshTokenResponseMultiError, or nil if none found.
func (m *RefreshTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *RefreshTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if len(errors) > 0 {
		return RefreshTokenResponseMultiError(errors)
	}

	return nil
}

// RefreshTokenResponseMultiError is an error wrapping multiple validation
// errors returned by RefreshTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type RefreshTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RefreshTokenResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RefreshTokenResponseMultiError) AllErrors() []error { return m }

// RefreshTokenResponseValidationError is the validation error returned by
// RefreshTokenResponse.Validate if the designated constraints aren't met.
type RefreshTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RefreshTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RefreshTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RefreshTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RefreshTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RefreshTokenResponseValidationError) ErrorName() string {
	return "RefreshTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RefreshTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRefreshTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RefreshTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RefreshTokenResponseValidationError{}

// Validate checks the field values on RefreshTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *RefreshTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RefreshTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RefreshTokenRequestMultiError, or nil if none found.
func (m *RefreshTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RefreshTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Token

	if len(errors) > 0 {
		return RefreshTokenRequestMultiError(errors)
	}

	return nil
}

// RefreshTokenRequestMultiError is an error wrapping multiple validation
// errors returned by RefreshTokenRequest.ValidateAll() if the designated
// constraints aren't met.
type RefreshTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RefreshTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RefreshTokenRequestMultiError) AllErrors() []error { return m }

// RefreshTokenRequestValidationError is the validation error returned by
// RefreshTokenRequest.Validate if the designated constraints aren't met.
type RefreshTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RefreshTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RefreshTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RefreshTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RefreshTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RefreshTokenRequestValidationError) ErrorName() string {
	return "RefreshTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e RefreshTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRefreshTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RefreshTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RefreshTokenRequestValidationError{}
