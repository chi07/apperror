package apperror_test

import (
	"net/http"
	"testing"

	"github.com/chi07/apperror"
)

func TestNewErrMissingField(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrRequiredField,
		Message:  apperror.NewErrorMessage("missing required field `name`"),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrMissingField(fieldName)

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error message: %s, but got: %s", expectedError.Error(), err.Error())
	}

	if err.(apperror.AppError).Code != expectedError.Code {
		t.Errorf("Expected error code: %s, but got: %s", expectedError.Code, err.(apperror.AppError).Code)
	}

	if err.(apperror.AppError).HttpCode != expectedError.HttpCode {
		t.Errorf("Expected HTTP code: %d, but got: %d", expectedError.HttpCode, err.(apperror.AppError).HttpCode)
	}
}

func TestNewErrDuplicatedValue(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrDuplicatedRecord,
		Message:  apperror.NewErrorMessage("duplicated value field `name` already used."),
		HttpCode: http.StatusConflict,
	}

	err := apperror.NewErrDuplicatedValue(fieldName)

	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error message: %s, but got: %s", expectedError.Error(), err.Error())
	}

	if err.(apperror.AppError).Code != expectedError.Code {
		t.Errorf("Expected error code: %s, but got: %s", expectedError.Code, err.(apperror.AppError).Code)
	}

	if err.(apperror.AppError).HttpCode != expectedError.HttpCode {
		t.Errorf("Expected HTTP code: %d, but got: %d", expectedError.HttpCode, err.(apperror.AppError).HttpCode)
	}
}
