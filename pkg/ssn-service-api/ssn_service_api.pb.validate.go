// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: ozonmp/ssn_service_api/v1/ssn_service_api.proto

package ssn_service_api

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

// Validate checks the field values on Service with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *Service) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Id

	// no validation rules for Name

	// no validation rules for Description

	if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
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

// Validate checks the field values on CreateServiceV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateServiceV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 1 || l > 100 {
		return CreateServiceV1RequestValidationError{
			field:  "Name",
			reason: "value length must be between 1 and 100 runes, inclusive",
		}
	}

	if l := utf8.RuneCountInString(m.GetDescription()); l < 1 || l > 200 {
		return CreateServiceV1RequestValidationError{
			field:  "Description",
			reason: "value length must be between 1 and 200 runes, inclusive",
		}
	}

	return nil
}

// CreateServiceV1RequestValidationError is the validation error returned by
// CreateServiceV1Request.Validate if the designated constraints aren't met.
type CreateServiceV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateServiceV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateServiceV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateServiceV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateServiceV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateServiceV1RequestValidationError) ErrorName() string {
	return "CreateServiceV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateServiceV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateServiceV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateServiceV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateServiceV1RequestValidationError{}

// Validate checks the field values on CreateServiceV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *CreateServiceV1Response) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ServiceId

	return nil
}

// CreateServiceV1ResponseValidationError is the validation error returned by
// CreateServiceV1Response.Validate if the designated constraints aren't met.
type CreateServiceV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateServiceV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateServiceV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateServiceV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateServiceV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateServiceV1ResponseValidationError) ErrorName() string {
	return "CreateServiceV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateServiceV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateServiceV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateServiceV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateServiceV1ResponseValidationError{}

// Validate checks the field values on DescribeServiceV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribeServiceV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetServiceId() <= 0 {
		return DescribeServiceV1RequestValidationError{
			field:  "ServiceId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// DescribeServiceV1RequestValidationError is the validation error returned by
// DescribeServiceV1Request.Validate if the designated constraints aren't met.
type DescribeServiceV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribeServiceV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribeServiceV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribeServiceV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribeServiceV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribeServiceV1RequestValidationError) ErrorName() string {
	return "DescribeServiceV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e DescribeServiceV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribeServiceV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribeServiceV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribeServiceV1RequestValidationError{}

// Validate checks the field values on DescribeServiceV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *DescribeServiceV1Response) Validate() error {
	if m == nil {
		return nil
	}

	if v, ok := interface{}(m.GetService()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DescribeServiceV1ResponseValidationError{
				field:  "Service",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// DescribeServiceV1ResponseValidationError is the validation error returned by
// DescribeServiceV1Response.Validate if the designated constraints aren't met.
type DescribeServiceV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DescribeServiceV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DescribeServiceV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DescribeServiceV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DescribeServiceV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DescribeServiceV1ResponseValidationError) ErrorName() string {
	return "DescribeServiceV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e DescribeServiceV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDescribeServiceV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DescribeServiceV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DescribeServiceV1ResponseValidationError{}

// Validate checks the field values on UpdateServiceV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UpdateServiceV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetServiceId() <= 0 {
		return UpdateServiceV1RequestValidationError{
			field:  "ServiceId",
			reason: "value must be greater than 0",
		}
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 1 || l > 100 {
		return UpdateServiceV1RequestValidationError{
			field:  "Name",
			reason: "value length must be between 1 and 100 runes, inclusive",
		}
	}

	if l := utf8.RuneCountInString(m.GetDescription()); l < 1 || l > 200 {
		return UpdateServiceV1RequestValidationError{
			field:  "Description",
			reason: "value length must be between 1 and 200 runes, inclusive",
		}
	}

	return nil
}

// UpdateServiceV1RequestValidationError is the validation error returned by
// UpdateServiceV1Request.Validate if the designated constraints aren't met.
type UpdateServiceV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateServiceV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateServiceV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateServiceV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateServiceV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateServiceV1RequestValidationError) ErrorName() string {
	return "UpdateServiceV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateServiceV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateServiceV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateServiceV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateServiceV1RequestValidationError{}

// Validate checks the field values on UpdateServiceV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *UpdateServiceV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// UpdateServiceV1ResponseValidationError is the validation error returned by
// UpdateServiceV1Response.Validate if the designated constraints aren't met.
type UpdateServiceV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateServiceV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateServiceV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateServiceV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateServiceV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateServiceV1ResponseValidationError) ErrorName() string {
	return "UpdateServiceV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateServiceV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateServiceV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateServiceV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateServiceV1ResponseValidationError{}

// Validate checks the field values on ListServicesV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListServicesV1Request) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Offset

	if val := m.GetLimit(); val <= 0 || val > 500 {
		return ListServicesV1RequestValidationError{
			field:  "Limit",
			reason: "value must be inside range (0, 500]",
		}
	}

	return nil
}

// ListServicesV1RequestValidationError is the validation error returned by
// ListServicesV1Request.Validate if the designated constraints aren't met.
type ListServicesV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListServicesV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListServicesV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListServicesV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListServicesV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListServicesV1RequestValidationError) ErrorName() string {
	return "ListServicesV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListServicesV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListServicesV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListServicesV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListServicesV1RequestValidationError{}

// Validate checks the field values on ListServicesV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ListServicesV1Response) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetServices() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListServicesV1ResponseValidationError{
					field:  fmt.Sprintf("Services[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListServicesV1ResponseValidationError is the validation error returned by
// ListServicesV1Response.Validate if the designated constraints aren't met.
type ListServicesV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListServicesV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListServicesV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListServicesV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListServicesV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListServicesV1ResponseValidationError) ErrorName() string {
	return "ListServicesV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListServicesV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListServicesV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListServicesV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListServicesV1ResponseValidationError{}

// Validate checks the field values on RemoveServiceV1Request with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemoveServiceV1Request) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetServiceId() <= 0 {
		return RemoveServiceV1RequestValidationError{
			field:  "ServiceId",
			reason: "value must be greater than 0",
		}
	}

	return nil
}

// RemoveServiceV1RequestValidationError is the validation error returned by
// RemoveServiceV1Request.Validate if the designated constraints aren't met.
type RemoveServiceV1RequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveServiceV1RequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveServiceV1RequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveServiceV1RequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveServiceV1RequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveServiceV1RequestValidationError) ErrorName() string {
	return "RemoveServiceV1RequestValidationError"
}

// Error satisfies the builtin error interface
func (e RemoveServiceV1RequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveServiceV1Request.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveServiceV1RequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveServiceV1RequestValidationError{}

// Validate checks the field values on RemoveServiceV1Response with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *RemoveServiceV1Response) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// RemoveServiceV1ResponseValidationError is the validation error returned by
// RemoveServiceV1Response.Validate if the designated constraints aren't met.
type RemoveServiceV1ResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveServiceV1ResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveServiceV1ResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveServiceV1ResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveServiceV1ResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveServiceV1ResponseValidationError) ErrorName() string {
	return "RemoveServiceV1ResponseValidationError"
}

// Error satisfies the builtin error interface
func (e RemoveServiceV1ResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveServiceV1Response.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveServiceV1ResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveServiceV1ResponseValidationError{}

// Validate checks the field values on ServiceEventPayload with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ServiceEventPayload) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetServiceId() <= 0 {
		return ServiceEventPayloadValidationError{
			field:  "ServiceId",
			reason: "value must be greater than 0",
		}
	}

	if l := utf8.RuneCountInString(m.GetName()); l < 1 || l > 100 {
		return ServiceEventPayloadValidationError{
			field:  "Name",
			reason: "value length must be between 1 and 100 runes, inclusive",
		}
	}

	if l := utf8.RuneCountInString(m.GetDescription()); l < 1 || l > 200 {
		return ServiceEventPayloadValidationError{
			field:  "Description",
			reason: "value length must be between 1 and 200 runes, inclusive",
		}
	}

	return nil
}

// ServiceEventPayloadValidationError is the validation error returned by
// ServiceEventPayload.Validate if the designated constraints aren't met.
type ServiceEventPayloadValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ServiceEventPayloadValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ServiceEventPayloadValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ServiceEventPayloadValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ServiceEventPayloadValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ServiceEventPayloadValidationError) ErrorName() string {
	return "ServiceEventPayloadValidationError"
}

// Error satisfies the builtin error interface
func (e ServiceEventPayloadValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sServiceEventPayload.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ServiceEventPayloadValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ServiceEventPayloadValidationError{}

// Validate checks the field values on ServiceEvent with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ServiceEvent) Validate() error {
	if m == nil {
		return nil
	}

	if m.GetId() <= 0 {
		return ServiceEventValidationError{
			field:  "Id",
			reason: "value must be greater than 0",
		}
	}

	if m.GetServiceId() <= 0 {
		return ServiceEventValidationError{
			field:  "ServiceId",
			reason: "value must be greater than 0",
		}
	}

	if utf8.RuneCountInString(m.GetType()) < 1 {
		return ServiceEventValidationError{
			field:  "Type",
			reason: "value length must be at least 1 runes",
		}
	}

	if utf8.RuneCountInString(m.GetSubtype()) < 1 {
		return ServiceEventValidationError{
			field:  "Subtype",
			reason: "value length must be at least 1 runes",
		}
	}

	if v, ok := interface{}(m.GetPayload()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ServiceEventValidationError{
				field:  "Payload",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	return nil
}

// ServiceEventValidationError is the validation error returned by
// ServiceEvent.Validate if the designated constraints aren't met.
type ServiceEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ServiceEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ServiceEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ServiceEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ServiceEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ServiceEventValidationError) ErrorName() string { return "ServiceEventValidationError" }

// Error satisfies the builtin error interface
func (e ServiceEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sServiceEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ServiceEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ServiceEventValidationError{}
