package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code     ErrorCode    `json:"code"`
	Message  ErrorMessage `json:"message"`
	HttpCode int          `json:"-"`
}

func (e AppError) Error() string {
	return string(e.Message)
}

func NewErrInternalServer(msg string) error {
	message := fmt.Sprintf("%s `%s`", getErrorMessage(ErrInternalError), msg)
	return AppError{
		Code:     ErrInternalError,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrMissingField(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", getErrorMessage(ErrRequiredField), fieldName)
	return AppError{
		Code:     ErrRequiredField,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrInvalidValue(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", getErrorMessage(ErrInvalidFieldValue), fieldName)
	return AppError{
		Code:     ErrInvalidFieldValue,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrInvalidMinValue(fieldName string, val int) error {
	message := fmt.Sprintf("%s `%s`. Its should be less than %d", getErrorMessage(ErrInvalidFieldValue), fieldName, val)
	return AppError{
		Code:     ErrInvalidFieldValue,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrInvalidMaxValue(fieldName string, val int) error {
	message := fmt.Sprintf("%s `%s`. Its should be bigger than %d", getErrorMessage(ErrInvalidFieldValue), fieldName, val)
	return AppError{
		Code:     ErrInvalidFieldValue,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrInvalidType(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", getErrorMessage(ErrInvalidFieldType), fieldName)
	return AppError{
		Code:     ErrInvalidFieldType,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrNotMatched(fieldName string, fieldName2 string) error {
	message := fmt.Sprintf("`%s` and `%s` are not match", fieldName, fieldName2)
	return AppError{
		Code:     ErrNotMatched,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrSuspended(fieldName string) error {
	message := fmt.Sprintf("`%s` is suspended.", fieldName)
	return AppError{
		Code:     ErrSuspended,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewErrRecordNotfound(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", getErrorMessage(ErrRecordNotFound), fieldName)
	return AppError{
		Code:     ErrRecordNotFound,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusNotFound,
	}
}

func NewErrDuplicatedValue(fieldName string) error {
	message := fmt.Sprintf("%s `%s` already used.", getErrorMessage(ErrDuplicatedRecord), fieldName)
	return AppError{
		Code:     ErrDuplicatedRecord,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusConflict,
	}
}
