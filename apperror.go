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

func NewMissingField(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", messages[ErrRequiredField], fieldName)
	return AppError{
		Code:     ErrRequiredField,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewInvalidValue(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", messages[ErrInvalidFieldValue], fieldName)
	return AppError{
		Code:     ErrInvalidFieldValue,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewInvalidMinValue(fieldName string, val int) error {
	message := fmt.Sprintf("%s `%s`. Its should be bigger than %d", messages[ErrInvalidFieldValue], fieldName, val)
	return AppError{
		Code:     ErrInvalidFieldValue,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewInvalidMaxValue(fieldName string, val int) error {
	message := fmt.Sprintf("%s `%s`. Its should be less than %d", messages[ErrInvalidFieldValue], fieldName, val)
	return AppError{
		Code:     ErrInvalidFieldValue,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewInvalidType(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", messages[ErrInvalidFieldType], fieldName)
	return AppError{
		Code:     ErrInvalidFieldType,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewNotMatched(fieldName string, fieldName2 string) error {
	message := fmt.Sprintf("%s `%s` are not match", fieldName, fieldName2)
	return AppError{
		Code:     ErrInvalidFieldType,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewSuspended(fieldName string) error {
	message := fmt.Sprintf("`%s` is suspended.", fieldName)
	return AppError{
		Code:     ErrInvalidFieldType,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewRecordNotfound(fieldName string) error {
	message := fmt.Sprintf("%s `%s`", messages[ErrRecordNotFound], fieldName)
	return AppError{
		Code:     ErrRecordNotFound,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusBadRequest,
	}
}

func NewDuplicatedValue(fieldName string) error {
	message := fmt.Sprintf("%s `%s` already used.", messages[ErrDuplicatedRecord], fieldName)
	return AppError{
		Code:     ErrDuplicatedRecord,
		Message:  NewErrorMessage(message),
		HttpCode: http.StatusConflict,
	}
}
