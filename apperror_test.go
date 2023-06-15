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

func TestNewErrInvalidMaxValue(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrInvalidFieldValue,
		Message:  apperror.NewErrorMessage("invalid value filed `name`. Its should be bigger than 1"),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrInvalidMaxValue(fieldName, 1)

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

func TestNewErrInvalidMinValue(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrInvalidFieldValue,
		Message:  apperror.NewErrorMessage("invalid value filed `name`. Its should be less than 1"),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrInvalidMinValue(fieldName, 1)

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

func TestNewErrInvalidType(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrInvalidFieldType,
		Message:  apperror.NewErrorMessage("invalid type filed `name`"),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrInvalidType(fieldName)

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

func TestNewErrInvalidValue(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrInvalidFieldValue,
		Message:  apperror.NewErrorMessage("invalid value filed `name`"),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrInvalidValue(fieldName)

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

func TestNewErrNotMatched(t *testing.T) {
	field1, field2 := "username", "password"
	expectedError := apperror.AppError{
		Code:     apperror.ErrNotMatched,
		Message:  apperror.NewErrorMessage("`username` and `password` are not match"),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrNotMatched(field1, field2)

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

func TestNewErrRecordNotfound(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrRecordNotFound,
		Message:  apperror.NewErrorMessage("not found `name`"),
		HttpCode: http.StatusNotFound,
	}

	err := apperror.NewErrRecordNotfound(fieldName)

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

func TestNewErrSuspended(t *testing.T) {
	fieldName := "name"
	expectedError := apperror.AppError{
		Code:     apperror.ErrSuspended,
		Message:  apperror.NewErrorMessage("`name` is suspended."),
		HttpCode: http.StatusBadRequest,
	}

	err := apperror.NewErrSuspended(fieldName)

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

func TestAppError_Error(t *testing.T) {
	err := apperror.AppError{
		Code:     "V40001",
		Message:  "test failed",
		HttpCode: 400,
	}

	if err.Error() != "test failed" {
		t.Errorf("Expected error message: %s, but got: %s", "test failed", err.Error())
	}
}

func TestNewErrorMessage(t *testing.T) {
	message := "test failed"
	err := apperror.NewErrorMessage(message)

	if err != apperror.ErrorMessage(message) {
		t.Errorf("Expected error message: %s, but got: %s", apperror.ErrorMessage(message), err)
	}
}
