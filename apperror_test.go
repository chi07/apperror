// apperror_test.go
package apperror_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/chi07/apperror"
)

func mustAsAppError(t *testing.T, err error) *apperror.AppError {
	t.Helper()
	ae, ok := apperror.AsAppError(err)
	if !ok {
		t.Fatalf("expected *AppError via AsAppError, got %T (%v)", err, err)
	}
	return ae
}

func mustContainAll(t *testing.T, s string, subs ...string) {
	t.Helper()
	for _, sub := range subs {
		if sub == "" {
			continue
		}
		if !strings.Contains(s, sub) {
			t.Fatalf("expected %q to contain %q", s, sub)
		}
	}
}

/* -------------------- table-driven tests -------------------- */

func TestAppErrors_TableDriven(t *testing.T) {
	type tc struct {
		name     string
		makeErr  func() error
		wantCode apperror.Code
		wantHTTP int
		contains []string
		checkIs  bool // check errors.Is with code
	}

	tests := []tc{
		{
			name:     "Unauthorized",
			makeErr:  func() error { return apperror.NewErrUnauthorized("token") },
			wantCode: apperror.ErrUnauthorized,
			wantHTTP: 401,
			contains: []string{"token", "`"},
			checkIs:  true,
		},
		{
			name:     "Forbidden (permission denied)",
			makeErr:  func() error { return apperror.NewErrStatusForbidden("no access") },
			wantCode: apperror.ErrPermissionDenied,
			wantHTTP: 403,
			contains: []string{"no access"},
			checkIs:  true,
		},
		{
			name:     "Missing field",
			makeErr:  func() error { return apperror.NewErrMissingField("email") },
			wantCode: apperror.ErrRequiredField,
			wantHTTP: 400,
			contains: []string{"email", "`"},
			checkIs:  true,
		},
		{
			name:     "Invalid value",
			makeErr:  func() error { return apperror.NewErrInvalidValue("age") },
			wantCode: apperror.ErrInvalidFieldValue,
			wantHTTP: 400,
			contains: []string{"age", "`"},
			checkIs:  true,
		},
		{
			name:     "Invalid min value",
			makeErr:  func() error { return apperror.NewErrInvalidMinValue("price", 10) },
			wantCode: apperror.ErrInvalidFieldValue,
			wantHTTP: 400,
			contains: []string{"price", "10"},
			checkIs:  true,
		},
		{
			name:     "Invalid max value",
			makeErr:  func() error { return apperror.NewErrInvalidMaxValue("qty", 99) },
			wantCode: apperror.ErrInvalidFieldValue,
			wantHTTP: 400,
			contains: []string{"qty", "99"},
			checkIs:  true,
		},
		{
			name:     "Invalid type",
			makeErr:  func() error { return apperror.NewErrInvalidType("file") },
			wantCode: apperror.ErrInvalidFieldType,
			wantHTTP: 400,
			contains: []string{"file", "`"},
			checkIs:  true,
		},
		{
			name:     "Not matched",
			makeErr:  func() error { return apperror.NewErrNotMatched("password", "confirm_password") },
			wantCode: apperror.ErrNotMatched,
			wantHTTP: 400,
			contains: []string{"password", "confirm_password", "not match"},
			checkIs:  true,
		},
		{
			name:     "Suspended",
			makeErr:  func() error { return apperror.NewErrSuspended("account") },
			wantCode: apperror.ErrSuspended,
			wantHTTP: 400,
			contains: []string{"account", "suspended"},
			checkIs:  true,
		},
		{
			name:     "Record not found",
			makeErr:  func() error { return apperror.NewErrRecordNotfound("course") },
			wantCode: apperror.ErrRecordNotFound,
			wantHTTP: 404,
			contains: []string{"course", "`"},
			checkIs:  true,
		},
		{
			name:     "Duplicated record",
			makeErr:  func() error { return apperror.NewErrDuplicatedValue("email") },
			wantCode: apperror.ErrDuplicatedRecord,
			wantHTTP: 409,
			contains: []string{"email", "already used"},
			checkIs:  true,
		},
		{
			name:     "Not activated",
			makeErr:  func() error { return apperror.NewErrNotActivated("account") },
			wantCode: apperror.ErrNotActivated,
			wantHTTP: 403,
			contains: []string{"account", "activate"},
			checkIs:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.makeErr()
			ae := mustAsAppError(t, err)

			if ae.Code != tt.wantCode {
				t.Fatalf("Code: got %v, want %v", ae.Code, tt.wantCode)
			}
			if ae.HTTPCode != tt.wantHTTP || ae.GetCode() != tt.wantHTTP {
				t.Fatalf("HTTPCode/GetCode: got (%d,%d), want (%d,%d)", ae.HTTPCode, ae.GetCode(), tt.wantHTTP, tt.wantHTTP)
			}
			mustContainAll(t, ae.Error(), tt.contains...)

			if tt.checkIs {
				// errors.Is match theo Code
				if !errors.Is(err, &apperror.AppError{Code: tt.wantCode}) {
					t.Fatalf("errors.Is should match by code %v", tt.wantCode)
				}
			}
		})
	}
}

func TestInternalServer_WithAndWithoutCause(t *testing.T) {
	// With cause
	sentinel := errors.New("db down")
	err := apperror.NewErrInternalServer("processing", sentinel)
	ae := mustAsAppError(t, err)

	if ae.Code != apperror.ErrInternalError {
		t.Fatalf("unexpected code: got %v want %v", ae.Code, apperror.ErrInternalError)
	}
	if ae.HTTPCode != 500 || ae.GetCode() != 500 {
		t.Fatalf("unexpected http code: got (%d,%d) want 500", ae.HTTPCode, ae.GetCode())
	}
	// Error() phải chứa cả message và cause
	mustContainAll(t, ae.Error(), "processing", "db down")
	// errors.Is theo Code
	if !errors.Is(err, &apperror.AppError{Code: apperror.ErrInternalError}) {
		t.Fatalf("errors.Is should match by internal code")
	}
	// errors.Is theo cause (Unwrap)
	if !errors.Is(err, sentinel) {
		t.Fatalf("errors.Is should find original cause via Unwrap")
	}

	// Without cause
	err2 := apperror.NewErrInternalServer("panic-free")
	ae2 := mustAsAppError(t, err2)
	if ae2.Code != apperror.ErrInternalError || ae2.HTTPCode != 500 {
		t.Fatalf("unexpected (code,http): got (%v,%d)", ae2.Code, ae2.HTTPCode)
	}
	mustContainAll(t, ae2.Error(), "panic-free")
}

func TestAsAppError_NotAppError(t *testing.T) {
	ae, ok := apperror.AsAppError(errors.New("plain error"))
	if ok || ae != nil {
		t.Fatalf("expected (nil,false) for non-AppError, got (%v,%v)", ae, ok)
	}
}

func TestAsAppError_NilError(t *testing.T) {
	var err error = nil
	ae, ok := apperror.AsAppError(err)
	if ok || ae != nil {
		t.Fatalf("expected (nil,false) for nil error, got (%v,%v)", ae, ok)
	}
}

func TestErrCode_GetHTTPCode(t *testing.T) {
	err := apperror.NewErrUnauthorized("token")
	ae := mustAsAppError(t, err)

	if ae.ErrCode() != apperror.ErrUnauthorized {
		t.Fatalf("ErrCode: got %v, want %v", ae.ErrCode(), apperror.ErrUnauthorized)
	}
	if ae.GetHTTPCode() != 401 {
		t.Fatalf("GetHTTPCode: got %d, want 401", ae.GetHTTPCode())
	}
}

func TestNewAppError_Custom(t *testing.T) {
	cause := errors.New("root cause")
	err := apperror.NewAppError(422, "C422", "custom error", cause)

	if err.Code != "C422" {
		t.Fatalf("Code: got %v, want C422", err.Code)
	}
	if err.HTTPCode != 422 {
		t.Fatalf("HTTPCode: got %d, want 422", err.HTTPCode)
	}
	if err.Message != "custom error" {
		t.Fatalf("Message: got %v, want 'custom error'", err.Message)
	}
	if !errors.Is(err, cause) {
		t.Fatal("errors.Is should find the cause")
	}
}

func TestWithCause(t *testing.T) {
	base := apperror.NewErrMissingField("email")
	ae1 := mustAsAppError(t, base)
	if ae1.Cause != nil {
		t.Fatal("base error should not have a cause")
	}

	cause := errors.New("underlying db error")
	ae2 := ae1.WithCause(cause)

	if ae2.Cause != cause {
		t.Fatal("WithCause did not set the cause")
	}
	// Original should be unchanged
	if ae1.Cause != nil {
		t.Fatal("WithCause mutated the original error")
	}
	// Both should have the same code and HTTP status
	if ae1.Code != ae2.Code || ae1.HTTPCode != ae2.HTTPCode {
		t.Fatal("WithCause should preserve code and http status")
	}
}

func TestGetters_NilSafety(t *testing.T) {
	var ae *apperror.AppError
	if ae.ErrCode() != "" {
		t.Fatalf("ErrCode on nil should return empty, got %v", ae.ErrCode())
	}
	if ae.GetHTTPCode() != 0 {
		t.Fatalf("GetHTTPCode on nil should return 0, got %d", ae.GetHTTPCode())
	}
	if ae.GetCode() != 0 {
		t.Fatalf("GetCode on nil should return 0, got %d", ae.GetCode())
	}
	if ae.Error() != "" {
		t.Fatalf("Error on nil should return empty string, got %q", ae.Error())
	}
}

func BenchmarkNewErrMissingField(b *testing.B) {
	for b.Loop() {
		_ = apperror.NewErrMissingField("email")
	}
}

func BenchmarkNewErrNotMatched(b *testing.B) {
	for b.Loop() {
		_ = apperror.NewErrNotMatched("password", "confirm_password")
	}
}

func BenchmarkAppErrorErrorWithoutCause(b *testing.B) {
	err := apperror.NewErrMissingField("email")
	b.ResetTimer()

	for b.Loop() {
		_ = err.Error()
	}
}

func BenchmarkAppErrorErrorWithCause(b *testing.B) {
	err := apperror.NewErrInternalServer("processing", errors.New("db down"))
	b.ResetTimer()

	for b.Loop() {
		_ = err.Error()
	}
}
