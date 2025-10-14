// Package apperror package apperror
package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	Code    string
	Message string
)

func NewMessage(v string) Message {
	return Message(v)
}

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

var messages = map[Code]Message{
	ErrInternalError:     "internal server error. Please contact admin support",
	ErrRequiredField:     "missing required field",
	ErrUnauthorized:      "unauthorized",
	ErrPermissionDenied:  "permission denied",
	ErrInvalidFieldValue: "invalid value filed",
	ErrInvalidFieldType:  "invalid type filed",
	ErrDuplicatedRecord:  "duplicated value field",
	ErrRecordNotFound:    "not found",
	ErrSuspended:         "suspend",
	ErrNotActivated:      "not activated",
	ErrNotMatched:        "not matched",
}

func getMessage(code Code) Message {
	switch code {
	case ErrInternalError:
		return messages[ErrInternalError]
	case ErrRequiredField:
		return messages[ErrRequiredField]
	case ErrUnauthorized:
		return messages[ErrUnauthorized]
	case ErrInvalidFieldType:
		return messages[ErrInvalidFieldType]
	case ErrPermissionDenied:
		return messages[ErrPermissionDenied]
	case ErrInvalidFieldValue:
		return messages[ErrInvalidFieldValue]
	case ErrDuplicatedRecord:
		return messages[ErrDuplicatedRecord]
	case ErrSuspended:
		return messages[ErrSuspended]
	case ErrNotActivated:
		return messages[ErrNotActivated]
	case ErrRecordNotFound:
		return messages[ErrRecordNotFound]
	case ErrNotMatched:
		return messages[ErrNotMatched]
	default:
		return messages[ErrInternalError]
	}
}

type AppError struct {
	Code     Code    `json:"code"`
	Message  Message `json:"message"`
	HTTPCode int     `json:"-"`
	Cause    error   `json:"-"`
}

func (e *AppError) GetCode() int {
	if e == nil {
		return 0
	}

	return e.HTTPCode
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", string(e.Message), e.Cause)
	}
	return string(e.Message)
}

func (e *AppError) Unwrap() error { return e.Cause }

func (e *AppError) Is(target error) bool {
	var t *AppError
	ok := errors.As(target, &t)
	return ok && e.Code == t.Code
}

type spec struct {
	status int
	tmpl   string // template bổ sung sau getErrorMessage(code)
}

const (
	tmplBaseTwoPlaceholders = "%s `%s`"
)

// Khai báo mapping code -> http status + template bổ sung
var errorSpec = map[Code]spec{
	ErrInternalError:     {http.StatusInternalServerError, tmplBaseTwoPlaceholders},
	ErrUnauthorized:      {http.StatusUnauthorized, tmplBaseTwoPlaceholders},
	ErrPermissionDenied:  {http.StatusForbidden, tmplBaseTwoPlaceholders},
	ErrRequiredField:     {http.StatusBadRequest, tmplBaseTwoPlaceholders},
	ErrInvalidFieldValue: {http.StatusBadRequest, tmplBaseTwoPlaceholders},
	ErrInvalidFieldType:  {http.StatusBadRequest, tmplBaseTwoPlaceholders},
	ErrRecordNotFound:    {http.StatusNotFound, tmplBaseTwoPlaceholders},

	ErrNotActivated:     {http.StatusForbidden, "%s `%s`. Please activate it"},
	ErrNotMatched:       {http.StatusBadRequest, "`%s` and `%s` do not match"},
	ErrSuspended:        {http.StatusBadRequest, "`%s` is suspended"},
	ErrDuplicatedRecord: {http.StatusConflict, "%s `%s` already used"},
}

func newAppError(code Code, cause error, args ...any) *AppError {
	sp, ok := errorSpec[code]
	if !ok {
		sp = spec{status: http.StatusInternalServerError, tmpl: "%s"}
	}

	baseStr := string(getMessage(code))

	if sp.tmpl == "%s" && len(args) == 0 {
		return &AppError{
			Code:     code,
			Message:  NewMessage(baseStr),
			HTTPCode: sp.status,
			Cause:    cause,
		}
	}

	var small [8]any
	n := 0
	small[n] = baseStr
	n++
	for i := 0; i < len(args) && n < len(small); i++ {
		small[n] = args[i]
		n++
	}

	var msg string
	if len(args) <= len(small)-1 {
		msg = fmt.Sprintf(sp.tmpl, small[:n]...)
	} else {
		all := make([]any, 0, len(args)+1)
		all = append(all, baseStr)
		all = append(all, args...)
		msg = fmt.Sprintf(sp.tmpl, all...)
	}

	return &AppError{
		Code:     code,
		Message:  NewMessage(msg),
		HTTPCode: sp.status,
		Cause:    cause,
	}
}

func NewErrInternalServer(msg string, causes ...error) error {
	return newAppError(ErrInternalError, firstErr(causes...), msg)
}

func NewErrUnauthorized(msg string) error { // fix typo: UnUnauthorized -> Unauthorized
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
	return newAppError(ErrInvalidFieldValue, nil, fmt.Sprintf("%s. It should be ≥ %d", fieldName, val))
}

func NewErrInvalidMaxValue(fieldName string, val int) error {
	return newAppError(ErrInvalidFieldValue, nil, fmt.Sprintf("%s. It should be ≤ %d", fieldName, val))
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

func firstErr(es ...error) error {
	for _, e := range es {
		if e != nil {
			return e
		}
	}
	return nil
}

func AsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
