// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: queue/v1/queue.proto

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

// Validate checks the field values on QueueSendRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *QueueSendRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueSendRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueSendRequestMultiError, or nil if none found.
func (m *QueueSendRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueSendRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetQueue()) > 256 {
		err := QueueSendRequestValidationError{
			field:  "Queue",
			reason: "value length must be at most 256 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_QueueSendRequest_Queue_Pattern.MatchString(m.GetQueue()) {
		err := QueueSendRequestValidationError{
			field:  "Queue",
			reason: "value does not match regex pattern \"^\\\\w+([.\\\\-]\\\\w+)*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTask() == nil {
		err := QueueSendRequestValidationError{
			field:  "Task",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetTask()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, QueueSendRequestValidationError{
					field:  "Task",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, QueueSendRequestValidationError{
					field:  "Task",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTask()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return QueueSendRequestValidationError{
				field:  "Task",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return QueueSendRequestMultiError(errors)
	}

	return nil
}

// QueueSendRequestMultiError is an error wrapping multiple validation errors
// returned by QueueSendRequest.ValidateAll() if the designated constraints
// aren't met.
type QueueSendRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueSendRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueSendRequestMultiError) AllErrors() []error { return m }

// QueueSendRequestValidationError is the validation error returned by
// QueueSendRequest.Validate if the designated constraints aren't met.
type QueueSendRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueSendRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueSendRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueSendRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueSendRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueSendRequestValidationError) ErrorName() string { return "QueueSendRequestValidationError" }

// Error satisfies the builtin error interface
func (e QueueSendRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueSendRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueSendRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueSendRequestValidationError{}

var _QueueSendRequest_Queue_Pattern = regexp.MustCompile("^\\w+([.\\-]\\w+)*$")

// Validate checks the field values on QueueSendResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *QueueSendResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueSendResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueSendResponseMultiError, or nil if none found.
func (m *QueueSendResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueSendResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return QueueSendResponseMultiError(errors)
	}

	return nil
}

// QueueSendResponseMultiError is an error wrapping multiple validation errors
// returned by QueueSendResponse.ValidateAll() if the designated constraints
// aren't met.
type QueueSendResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueSendResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueSendResponseMultiError) AllErrors() []error { return m }

// QueueSendResponseValidationError is the validation error returned by
// QueueSendResponse.Validate if the designated constraints aren't met.
type QueueSendResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueSendResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueSendResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueSendResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueSendResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueSendResponseValidationError) ErrorName() string {
	return "QueueSendResponseValidationError"
}

// Error satisfies the builtin error interface
func (e QueueSendResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueSendResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueSendResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueSendResponseValidationError{}

// Validate checks the field values on QueueSendBatchRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *QueueSendBatchRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueSendBatchRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueSendBatchRequestMultiError, or nil if none found.
func (m *QueueSendBatchRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueSendBatchRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetQueue()) > 256 {
		err := QueueSendBatchRequestValidationError{
			field:  "Queue",
			reason: "value length must be at most 256 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_QueueSendBatchRequest_Queue_Pattern.MatchString(m.GetQueue()) {
		err := QueueSendBatchRequestValidationError{
			field:  "Queue",
			reason: "value does not match regex pattern \"^\\\\w+([.\\\\-]\\\\w+)*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetTasks()) < 1 {
		err := QueueSendBatchRequestValidationError{
			field:  "Tasks",
			reason: "value must contain at least 1 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetTasks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, QueueSendBatchRequestValidationError{
						field:  fmt.Sprintf("Tasks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, QueueSendBatchRequestValidationError{
						field:  fmt.Sprintf("Tasks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return QueueSendBatchRequestValidationError{
					field:  fmt.Sprintf("Tasks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return QueueSendBatchRequestMultiError(errors)
	}

	return nil
}

// QueueSendBatchRequestMultiError is an error wrapping multiple validation
// errors returned by QueueSendBatchRequest.ValidateAll() if the designated
// constraints aren't met.
type QueueSendBatchRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueSendBatchRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueSendBatchRequestMultiError) AllErrors() []error { return m }

// QueueSendBatchRequestValidationError is the validation error returned by
// QueueSendBatchRequest.Validate if the designated constraints aren't met.
type QueueSendBatchRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueSendBatchRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueSendBatchRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueSendBatchRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueSendBatchRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueSendBatchRequestValidationError) ErrorName() string {
	return "QueueSendBatchRequestValidationError"
}

// Error satisfies the builtin error interface
func (e QueueSendBatchRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueSendBatchRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueSendBatchRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueSendBatchRequestValidationError{}

var _QueueSendBatchRequest_Queue_Pattern = regexp.MustCompile("^\\w+([.\\-]\\w+)*$")

// Validate checks the field values on QueueSendBatchResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *QueueSendBatchResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueSendBatchResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueSendBatchResponseMultiError, or nil if none found.
func (m *QueueSendBatchResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueSendBatchResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetFailedTasks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, QueueSendBatchResponseValidationError{
						field:  fmt.Sprintf("FailedTasks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, QueueSendBatchResponseValidationError{
						field:  fmt.Sprintf("FailedTasks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return QueueSendBatchResponseValidationError{
					field:  fmt.Sprintf("FailedTasks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return QueueSendBatchResponseMultiError(errors)
	}

	return nil
}

// QueueSendBatchResponseMultiError is an error wrapping multiple validation
// errors returned by QueueSendBatchResponse.ValidateAll() if the designated
// constraints aren't met.
type QueueSendBatchResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueSendBatchResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueSendBatchResponseMultiError) AllErrors() []error { return m }

// QueueSendBatchResponseValidationError is the validation error returned by
// QueueSendBatchResponse.Validate if the designated constraints aren't met.
type QueueSendBatchResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueSendBatchResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueSendBatchResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueSendBatchResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueSendBatchResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueSendBatchResponseValidationError) ErrorName() string {
	return "QueueSendBatchResponseValidationError"
}

// Error satisfies the builtin error interface
func (e QueueSendBatchResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueSendBatchResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueSendBatchResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueSendBatchResponseValidationError{}

// Validate checks the field values on QueueReceiveRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *QueueReceiveRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueReceiveRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueReceiveRequestMultiError, or nil if none found.
func (m *QueueReceiveRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueReceiveRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetQueue()) > 256 {
		err := QueueReceiveRequestValidationError{
			field:  "Queue",
			reason: "value length must be at most 256 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_QueueReceiveRequest_Queue_Pattern.MatchString(m.GetQueue()) {
		err := QueueReceiveRequestValidationError{
			field:  "Queue",
			reason: "value does not match regex pattern \"^\\\\w+([.\\\\-]\\\\w+)*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Depth

	if len(errors) > 0 {
		return QueueReceiveRequestMultiError(errors)
	}

	return nil
}

// QueueReceiveRequestMultiError is an error wrapping multiple validation
// errors returned by QueueReceiveRequest.ValidateAll() if the designated
// constraints aren't met.
type QueueReceiveRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueReceiveRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueReceiveRequestMultiError) AllErrors() []error { return m }

// QueueReceiveRequestValidationError is the validation error returned by
// QueueReceiveRequest.Validate if the designated constraints aren't met.
type QueueReceiveRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueReceiveRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueReceiveRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueReceiveRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueReceiveRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueReceiveRequestValidationError) ErrorName() string {
	return "QueueReceiveRequestValidationError"
}

// Error satisfies the builtin error interface
func (e QueueReceiveRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueReceiveRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueReceiveRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueReceiveRequestValidationError{}

var _QueueReceiveRequest_Queue_Pattern = regexp.MustCompile("^\\w+([.\\-]\\w+)*$")

// Validate checks the field values on QueueReceiveResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *QueueReceiveResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueReceiveResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueReceiveResponseMultiError, or nil if none found.
func (m *QueueReceiveResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueReceiveResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetTasks() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, QueueReceiveResponseValidationError{
						field:  fmt.Sprintf("Tasks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, QueueReceiveResponseValidationError{
						field:  fmt.Sprintf("Tasks[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return QueueReceiveResponseValidationError{
					field:  fmt.Sprintf("Tasks[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return QueueReceiveResponseMultiError(errors)
	}

	return nil
}

// QueueReceiveResponseMultiError is an error wrapping multiple validation
// errors returned by QueueReceiveResponse.ValidateAll() if the designated
// constraints aren't met.
type QueueReceiveResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueReceiveResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueReceiveResponseMultiError) AllErrors() []error { return m }

// QueueReceiveResponseValidationError is the validation error returned by
// QueueReceiveResponse.Validate if the designated constraints aren't met.
type QueueReceiveResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueReceiveResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueReceiveResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueReceiveResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueReceiveResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueReceiveResponseValidationError) ErrorName() string {
	return "QueueReceiveResponseValidationError"
}

// Error satisfies the builtin error interface
func (e QueueReceiveResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueReceiveResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueReceiveResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueReceiveResponseValidationError{}

// Validate checks the field values on QueueCompleteRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *QueueCompleteRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueCompleteRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueCompleteRequestMultiError, or nil if none found.
func (m *QueueCompleteRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueCompleteRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetQueue()) > 256 {
		err := QueueCompleteRequestValidationError{
			field:  "Queue",
			reason: "value length must be at most 256 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_QueueCompleteRequest_Queue_Pattern.MatchString(m.GetQueue()) {
		err := QueueCompleteRequestValidationError{
			field:  "Queue",
			reason: "value does not match regex pattern \"^\\\\w+([.\\\\-]\\\\w+)*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetLeaseId()) < 1 {
		err := QueueCompleteRequestValidationError{
			field:  "LeaseId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return QueueCompleteRequestMultiError(errors)
	}

	return nil
}

// QueueCompleteRequestMultiError is an error wrapping multiple validation
// errors returned by QueueCompleteRequest.ValidateAll() if the designated
// constraints aren't met.
type QueueCompleteRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueCompleteRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueCompleteRequestMultiError) AllErrors() []error { return m }

// QueueCompleteRequestValidationError is the validation error returned by
// QueueCompleteRequest.Validate if the designated constraints aren't met.
type QueueCompleteRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueCompleteRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueCompleteRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueCompleteRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueCompleteRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueCompleteRequestValidationError) ErrorName() string {
	return "QueueCompleteRequestValidationError"
}

// Error satisfies the builtin error interface
func (e QueueCompleteRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueCompleteRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueCompleteRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueCompleteRequestValidationError{}

var _QueueCompleteRequest_Queue_Pattern = regexp.MustCompile("^\\w+([.\\-]\\w+)*$")

// Validate checks the field values on QueueCompleteResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *QueueCompleteResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on QueueCompleteResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// QueueCompleteResponseMultiError, or nil if none found.
func (m *QueueCompleteResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *QueueCompleteResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return QueueCompleteResponseMultiError(errors)
	}

	return nil
}

// QueueCompleteResponseMultiError is an error wrapping multiple validation
// errors returned by QueueCompleteResponse.ValidateAll() if the designated
// constraints aren't met.
type QueueCompleteResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m QueueCompleteResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m QueueCompleteResponseMultiError) AllErrors() []error { return m }

// QueueCompleteResponseValidationError is the validation error returned by
// QueueCompleteResponse.Validate if the designated constraints aren't met.
type QueueCompleteResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e QueueCompleteResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e QueueCompleteResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e QueueCompleteResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e QueueCompleteResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e QueueCompleteResponseValidationError) ErrorName() string {
	return "QueueCompleteResponseValidationError"
}

// Error satisfies the builtin error interface
func (e QueueCompleteResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sQueueCompleteResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = QueueCompleteResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = QueueCompleteResponseValidationError{}

// Validate checks the field values on FailedTask with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *FailedTask) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FailedTask with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in FailedTaskMultiError, or
// nil if none found.
func (m *FailedTask) ValidateAll() error {
	return m.validate(true)
}

func (m *FailedTask) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTask()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FailedTaskValidationError{
					field:  "Task",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FailedTaskValidationError{
					field:  "Task",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTask()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FailedTaskValidationError{
				field:  "Task",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Message

	if len(errors) > 0 {
		return FailedTaskMultiError(errors)
	}

	return nil
}

// FailedTaskMultiError is an error wrapping multiple validation errors
// returned by FailedTask.ValidateAll() if the designated constraints aren't met.
type FailedTaskMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FailedTaskMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FailedTaskMultiError) AllErrors() []error { return m }

// FailedTaskValidationError is the validation error returned by
// FailedTask.Validate if the designated constraints aren't met.
type FailedTaskValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FailedTaskValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FailedTaskValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FailedTaskValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FailedTaskValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FailedTaskValidationError) ErrorName() string { return "FailedTaskValidationError" }

// Error satisfies the builtin error interface
func (e FailedTaskValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFailedTask.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FailedTaskValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FailedTaskValidationError{}

// Validate checks the field values on NitricTask with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *NitricTask) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on NitricTask with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in NitricTaskMultiError, or
// nil if none found.
func (m *NitricTask) ValidateAll() error {
	return m.validate(true)
}

func (m *NitricTask) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for LeaseId

	// no validation rules for PayloadType

	if all {
		switch v := interface{}(m.GetPayload()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, NitricTaskValidationError{
					field:  "Payload",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, NitricTaskValidationError{
					field:  "Payload",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetPayload()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return NitricTaskValidationError{
				field:  "Payload",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return NitricTaskMultiError(errors)
	}

	return nil
}

// NitricTaskMultiError is an error wrapping multiple validation errors
// returned by NitricTask.ValidateAll() if the designated constraints aren't met.
type NitricTaskMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NitricTaskMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NitricTaskMultiError) AllErrors() []error { return m }

// NitricTaskValidationError is the validation error returned by
// NitricTask.Validate if the designated constraints aren't met.
type NitricTaskValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NitricTaskValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NitricTaskValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NitricTaskValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NitricTaskValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NitricTaskValidationError) ErrorName() string { return "NitricTaskValidationError" }

// Error satisfies the builtin error interface
func (e NitricTaskValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNitricTask.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NitricTaskValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NitricTaskValidationError{}