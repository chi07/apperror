// Package apperror provides structured application error types with HTTP status codes.
package apperror

import (
	"errors"
	"net/http"
	"strconv"
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
}

var (
	specInternalError     = errorSpec{status: http.StatusInternalServerError, message: "internal server error. Please contact admin support"}
	specUnauthorized      = errorSpec{status: http.StatusUnauthorized, message: "unauthorized"}
	specPermissionDenied  = errorSpec{status: http.StatusForbidden, message: "permission denied"}
	specRequiredField     = errorSpec{status: http.StatusBadRequest, message: "missing required field"}
	specInvalidFieldValue = errorSpec{status: http.StatusBadRequest, message: "invalid value filed"}
	specInvalidFieldType  = errorSpec{status: http.StatusBadRequest, message: "invalid type filed"}
	specRecordNotFound    = errorSpec{status: http.StatusNotFound, message: "not found"}
	specDuplicatedRecord  = errorSpec{status: http.StatusConflict, message: "duplicated value field"}
	specNotMatched        = errorSpec{status: http.StatusBadRequest, message: "not matched"}
	specSuspended         = errorSpec{status: http.StatusBadRequest, message: "suspended"}
	specNotActivated      = errorSpec{status: http.StatusForbidden, message: "not activated"}
)

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
		return string(e.Message) + ": " + e.Cause.Error()
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

func newAppError1(code Code, cause error, sp errorSpec, arg string) *AppError {
	return &AppError{
		Code:     code,
		Message:  Message(formatMessage1(code, arg)),
		HTTPCode: sp.status,
		Cause:    cause,
	}
}

func newAppError2(code Code, cause error, sp errorSpec, arg1, arg2 string) *AppError {
	return &AppError{
		Code:     code,
		Message:  Message(formatMessage2(code, arg1, arg2)),
		HTTPCode: sp.status,
		Cause:    cause,
	}
}

func NewErrInternalServer(msg string, causes ...error) error {
	return newAppError1(ErrInternalError, firstErr(causes...), specInternalError, msg)
}

func NewErrUnauthorized(msg string) error {
	return newAppError1(ErrUnauthorized, nil, specUnauthorized, msg)
}

func NewErrNotActivated(field string) error {
	return newAppError1(ErrNotActivated, nil, specNotActivated, field)
}

func NewErrStatusForbidden(msg string) error {
	return newAppError1(ErrPermissionDenied, nil, specPermissionDenied, msg)
}

func NewErrMissingField(fieldName string) error {
	return newAppError1(ErrRequiredField, nil, specRequiredField, fieldName)
}

func NewErrInvalidValue(fieldName string) error {
	return newAppError1(ErrInvalidFieldValue, nil, specInvalidFieldValue, fieldName)
}

func NewErrInvalidMinValue(fieldName string, val int) error {
	return &AppError{
		Code:     ErrInvalidFieldValue,
		Message:  Message("invalid value filed `" + fieldName + ". It should be ≥ " + strconv.Itoa(val) + "`"),
		HTTPCode: specInvalidFieldValue.status,
	}
}

func NewErrInvalidMaxValue(fieldName string, val int) error {
	return &AppError{
		Code:     ErrInvalidFieldValue,
		Message:  Message("invalid value filed `" + fieldName + ". It should be ≤ " + strconv.Itoa(val) + "`"),
		HTTPCode: specInvalidFieldValue.status,
	}
}

func NewErrInvalidType(fieldName string) error {
	return newAppError1(ErrInvalidFieldType, nil, specInvalidFieldType, fieldName)
}

func NewErrNotMatched(field1, field2 string) error {
	return newAppError2(ErrNotMatched, nil, specNotMatched, field1, field2)
}

func NewErrSuspended(field string) error {
	return newAppError1(ErrSuspended, nil, specSuspended, field)
}

func NewErrRecordNotfound(field string) error {
	return newAppError1(ErrRecordNotFound, nil, specRecordNotFound, field)
}

func NewErrDuplicatedValue(field string) error {
	return newAppError1(ErrDuplicatedRecord, nil, specDuplicatedRecord, field)
}

func formatMessage1(code Code, arg string) string {
	switch code {
	case ErrInternalError:
		return "internal server error. Please contact admin support `" + arg + "`"
	case ErrUnauthorized:
		return "unauthorized `" + arg + "`"
	case ErrPermissionDenied:
		return "permission denied `" + arg + "`"
	case ErrRequiredField:
		return "missing required field `" + arg + "`"
	case ErrInvalidFieldValue:
		return "invalid value filed `" + arg + "`"
	case ErrInvalidFieldType:
		return "invalid type filed `" + arg + "`"
	case ErrRecordNotFound:
		return "not found `" + arg + "`"
	case ErrNotActivated:
		return "`" + arg + "` is not activated. Please activate it"
	case ErrSuspended:
		return "`" + arg + "` is suspended"
	case ErrDuplicatedRecord:
		return "duplicated value field `" + arg + "` already used"
	default:
		return specInternalError.message
	}
}

func formatMessage2(code Code, arg1, arg2 string) string {
	if code == ErrNotMatched {
		return "`" + arg1 + "` and `" + arg2 + "` do not match"
	}
	return formatMessage1(code, arg1)
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
