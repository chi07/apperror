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
	format  func(args ...any) string
	format1 func(string) string
	format2 func(string, string) string
}

// errorSpecs maps each error code to its HTTP status, base message, and format template.
// Templates are self-contained: args are passed directly without prepending sp.message.
var errorSpecs = map[Code]errorSpec{
	ErrInternalError: {
		status:  http.StatusInternalServerError,
		message: "internal server error. Please contact admin support",
		format1: func(v string) string {
			return "internal server error. Please contact admin support `" + v + "`"
		},
	},
	ErrUnauthorized: {
		status:  http.StatusUnauthorized,
		message: "unauthorized",
		format1: func(v string) string {
			return "unauthorized `" + v + "`"
		},
	},
	ErrPermissionDenied: {
		status:  http.StatusForbidden,
		message: "permission denied",
		format1: func(v string) string {
			return "permission denied `" + v + "`"
		},
	},
	ErrRequiredField: {
		status:  http.StatusBadRequest,
		message: "missing required field",
		format1: func(v string) string {
			return "missing required field `" + v + "`"
		},
	},
	ErrInvalidFieldValue: {
		status:  http.StatusBadRequest,
		message: "invalid value filed",
		format1: func(v string) string {
			return "invalid value filed `" + v + "`"
		},
	},
	ErrInvalidFieldType: {
		status:  http.StatusBadRequest,
		message: "invalid type filed",
		format1: func(v string) string {
			return "invalid type filed `" + v + "`"
		},
	},
	ErrRecordNotFound: {
		status:  http.StatusNotFound,
		message: "not found",
		format1: func(v string) string {
			return "not found `" + v + "`"
		},
	},
	ErrNotActivated: {
		status:  http.StatusForbidden,
		message: "not activated",
		format1: func(v string) string {
			return "`" + v + "` is not activated. Please activate it"
		},
	},
	ErrNotMatched: {
		status:  http.StatusBadRequest,
		message: "not matched",
		format2: func(v1, v2 string) string {
			return "`" + v1 + "` and `" + v2 + "` do not match"
		},
	},
	ErrSuspended: {
		status:  http.StatusBadRequest,
		message: "suspended",
		format1: func(v string) string {
			return "`" + v + "` is suspended"
		},
	},
	ErrDuplicatedRecord: {
		status:  http.StatusConflict,
		message: "duplicated value field",
		format1: func(v string) string {
			return "duplicated value field `" + v + "` already used"
		},
	},
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

func newAppError(code Code, cause error, args ...any) *AppError {
	sp, ok := errorSpecs[code]
	if !ok {
		code = ErrInternalError
		sp = errorSpecs[ErrInternalError]
	}

	return &AppError{
		Code:     code,
		Message:  Message(formatMessage(sp, args...)),
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
		Message:  Message("invalid value filed `" + fieldName + ". It should be ≥ " + strconv.Itoa(val) + "`"),
		HTTPCode: sp.status,
	}
}

func NewErrInvalidMaxValue(fieldName string, val int) error {
	sp := errorSpecs[ErrInvalidFieldValue]
	return &AppError{
		Code:     ErrInvalidFieldValue,
		Message:  Message("invalid value filed `" + fieldName + ". It should be ≤ " + strconv.Itoa(val) + "`"),
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

func formatMessage(sp errorSpec, args ...any) string {
	switch len(args) {
	case 0:
		return sp.message
	case 1:
		if sp.format1 != nil {
			if v, ok := args[0].(string); ok {
				return sp.format1(v)
			}
		}
	case 2:
		if sp.format2 != nil {
			v1, ok1 := args[0].(string)
			v2, ok2 := args[1].(string)
			if ok1 && ok2 {
				return sp.format2(v1, v2)
			}
		}
	}
	if sp.format != nil {
		return sp.format(args...)
	}
	return sp.message
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
