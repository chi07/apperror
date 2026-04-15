// Package apperror provides structured application error types with HTTP status codes.
package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// Code represents an application-level error code.
type Code string

// Message represents a human-readable error message.
type Message string

const (
	ErrRequiredField     Code = "V4000"
	ErrUnauthorized      Code = "V4001"
	ErrInvalidFieldValue Code = "V4002"
	ErrPermissionDenied  Code = "V4003"
	ErrRecordNotFound    Code = "V4004"
	ErrInvalidFieldType  Code = "V4005"
	ErrDuplicatedRecord  Code = "V4009"
	ErrNotMatched        Code = "V4010"
	ErrSuspended         Code = "V5001"
	ErrInternalError     Code = "S5003"
	ErrNotActivated      Code = "S4013"
)

type errorSpec struct {
	status  int
	message string
	tmpl    string
}

// errorSpecs maps each error code to its HTTP status, base message, and format template.
// Templates are self-contained: args are passed directly without prepending sp.message.
var errorSpecs = map[Code]errorSpec{
	ErrInternalError:     {http.StatusInternalServerError, "internal server error. Please contact admin support", "internal server error. Please contact admin support `%s`"},
	ErrUnauthorized:      {http.StatusUnauthorized, "unauthorized", "unauthorized `%s`"},
	ErrPermissionDenied:  {http.StatusForbidden, "permission denied", "permission denied `%s`"},
	ErrRequiredField:     {http.StatusBadRequest, "missing required field", "missing required field `%s`"},
	ErrInvalidFieldValue: {http.StatusBadRequest, "invalid value filed", "invalid value filed `%s`"},
	ErrInvalidFieldType:  {http.StatusBadRequest, "invalid type filed", "invalid type filed `%s`"},
	ErrRecordNotFound:    {http.StatusNotFound, "not found", "not found `%s`"},
	ErrNotActivated:      {http.StatusForbidden, "not activated", "`%s` is not activated. Please activate it"},
	ErrNotMatched:        {http.StatusBadRequest, "not matched", "`%s` and `%s` do not match"},
	ErrSuspended:         {http.StatusBadRequest, "suspended", "`%s` is suspended"},
	ErrDuplicatedRecord:  {http.StatusConflict, "duplicated value field", "duplicated value field `%s` already used"},
}

// AppError is a structured error with an error code, message, HTTP status, and optional cause.
type AppError struct {
	Code     Code    `json:"code"`
	Message  Message `json:"message"`
	HTTPCode int     `json:"-"`
	Cause    error   `json:"-"`
}

// ErrCode returns the error code.
func (e *AppError) ErrCode() Code {
	if e == nil {
		return ""
	}
	return e.Code
}

// GetHTTPCode returns the HTTP status code.
func (e *AppError) GetHTTPCode() int {
	if e == nil {
		return 0
	}
	return e.HTTPCode
}

// GetCode returns the HTTP status code.
// Deprecated: use GetHTTPCode instead.
func (e *AppError) GetCode() int {
	return e.GetHTTPCode()
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return string(e.Message)
}

// Unwrap returns the underlying cause for errors.Is/As chaining.
func (e *AppError) Unwrap() error { return e.Cause }

// Is reports whether target is an AppError with the same Code.
func (e *AppError) Is(target error) bool {
	var t *AppError
	return errors.As(target, &t) && e.Code == t.Code
}

func newAppError(code Code, cause error, args ...any) *AppError {
	sp, ok := errorSpecs[code]
	if !ok {
		sp = errorSpecs[ErrInternalError]
	}

	var msg string
	if len(args) == 0 {
		msg = sp.message
	} else {
		msg = fmt.Sprintf(sp.tmpl, args...)
	}

	return &AppError{
		Code:     code,
		Message:  Message(msg),
		HTTPCode: sp.status,
		Cause:    cause,
	}
}

func NewErrInternalServer(msg string, causes ...error) error {
	return newAppError(ErrInternalError, firstErr(causes...), msg)
}

func NewErrUnauthorized(msg string) error {
	return newAppError(ErrUnauthorized, nil, msg)
}

func NewErrNotActivated(field string) error {
	return newAppError(ErrNotActivated, nil, field)
}

func NewErrStatusForbidden(msg string) error {
	return newAppError(ErrPermissionDenied, nil, msg)
}

func NewErrMissingField(fieldName string) error {
	return newAppError(ErrRequiredField, nil, fieldName)
}

func NewErrInvalidValue(fieldName string) error {
	return newAppError(ErrInvalidFieldValue, nil, fieldName)
}

func NewErrInvalidMinValue(fieldName string, val int) error {
	sp := errorSpecs[ErrInvalidFieldValue]
	return &AppError{
		Code:     ErrInvalidFieldValue,
		Message:  Message(fmt.Sprintf("invalid value filed `%s. It should be ≥ %d`", fieldName, val)),
		HTTPCode: sp.status,
	}
}

func NewErrInvalidMaxValue(fieldName string, val int) error {
	sp := errorSpecs[ErrInvalidFieldValue]
	return &AppError{
		Code:     ErrInvalidFieldValue,
		Message:  Message(fmt.Sprintf("invalid value filed `%s. It should be ≤ %d`", fieldName, val)),
		HTTPCode: sp.status,
	}
}

func NewErrInvalidType(fieldName string) error {
	return newAppError(ErrInvalidFieldType, nil, fieldName)
}

func NewErrNotMatched(field1, field2 string) error {
	return newAppError(ErrNotMatched, nil, field1, field2)
}

func NewErrSuspended(field string) error {
	return newAppError(ErrSuspended, nil, field)
}

func NewErrRecordNotfound(field string) error {
	return newAppError(ErrRecordNotFound, nil, field)
}

func NewErrDuplicatedValue(field string) error {
	return newAppError(ErrDuplicatedRecord, nil, field)
}

// firstErr returns the first non-nil error from the provided list.
func firstErr(es ...error) error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

// AsAppError unwraps err into *AppError. Returns (nil, false) if not an AppError.
func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// NewAppError creates a custom AppError with the given HTTP status, code, and message.
// Useful for callers who need error codes not covered by the built-in constructors.
func NewAppError(httpCode int, code Code, message string, cause ...error) *AppError {
	return &AppError{
		Code:     code,
		Message:  Message(message),
		HTTPCode: httpCode,
		Cause:    firstErr(cause...),
	}
}

// WithCause returns a new AppError with the given cause attached.
// The original error is not modified.
func (e *AppError) WithCause(cause error) *AppError {
	if e == nil {
		return nil
	}
	newErr := *e
	newErr.Cause = cause
	return &newErr
}
