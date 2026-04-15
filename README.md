# apperror

Structured application errors with HTTP status codes for Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/chi07/apperror?v=1)](https://goreportcard.com/report/github.com/chi07/apperror)
[![codecov](https://codecov.io/gh/chi07/apperror/branch/main/graph/badge.svg)](https://codecov.io/gh/chi07/apperror)
[![CI](https://github.com/chi07/apperror/actions/workflows/ci.yml/badge.svg)](https://github.com/chi07/apperror/actions/workflows/ci.yml)

## Install

```sh
go get github.com/chi07/apperror
```

## Usage

```go
// Built-in constructors
err := apperror.NewErrMissingField("email")
err := apperror.NewErrInvalidValue("age")
err := apperror.NewErrInvalidMinValue("price", 10)
err := apperror.NewErrInvalidMaxValue("qty", 100)
err := apperror.NewErrRecordNotfound("user")
err := apperror.NewErrDuplicatedValue("email")
err := apperror.NewErrNotMatched("password", "confirm_password")
err := apperror.NewErrUnauthorized("token expired")
err := apperror.NewErrStatusForbidden("admin only")
err := apperror.NewErrNotActivated("account")
err := apperror.NewErrSuspended("account")
err := apperror.NewErrInternalServer("db query failed", cause)

// Custom error
err := apperror.NewAppError(422, "C422", "unprocessable entity", cause)

// Unwrap
ae, ok := apperror.AsAppError(err)
if ok {
    fmt.Println(ae.GetHTTPCode()) // e.g. 400
    fmt.Println(ae.ErrCode())     // e.g. "V4000"
    fmt.Println(ae.Error())       // human-readable message
}

// Attach a cause after the fact
ae2 := ae.WithCause(errors.New("root cause"))

// errors.Is by code
errors.Is(err, &apperror.AppError{Code: apperror.ErrRequiredField})
```

## Error codes

| Code | Constant | HTTP |
|------|----------|------|
| V4000 | `ErrRequiredField` | 400 |
| V4001 | `ErrUnauthorized` | 401 |
| V4002 | `ErrInvalidFieldValue` | 400 |
| V4003 | `ErrPermissionDenied` | 403 |
| V4004 | `ErrRecordNotFound` | 404 |
| V4005 | `ErrInvalidFieldType` | 400 |
| V4009 | `ErrDuplicatedRecord` | 409 |
| V4010 | `ErrNotMatched` | 400 |
| V5001 | `ErrSuspended` | 400 |
| S4013 | `ErrNotActivated` | 403 |
| S5003 | `ErrInternalError` | 500 |

## Development

```sh
make test        # run tests
make lint        # run linter
make coverage    # open coverage report
```
