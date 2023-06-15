package apperror

type (
	ErrorCode    string
	ErrorMessage string
)

func NewErrorMessage(v string) ErrorMessage {
	return ErrorMessage(v)
}

const (
	ErrRequiredField     ErrorCode = "V4000"
	ErrUnauthorized      ErrorCode = "V4001"
	ErrInvalidFieldValue ErrorCode = "V4002"
	ErrPermissionDenied  ErrorCode = "V4003"
	ErrRecordNotFound    ErrorCode = "V4004"
	ErrInvalidFieldType  ErrorCode = "V4005"
	ErrDuplicatedRecord  ErrorCode = "V4009"
	ErrNotMatched        ErrorCode = "V4010"
	ErrSuspended         ErrorCode = "V5001"
	ErrorInternalError   ErrorCode = "S5003"
)

var (
	messages = map[ErrorCode]ErrorMessage{
		ErrorInternalError:   "internal server error. Please contact admin support",
		ErrRequiredField:     "missing required field",
		ErrUnauthorized:      "unauthorized",
		ErrPermissionDenied:  "permission denied",
		ErrInvalidFieldValue: "invalid value filed",
		ErrInvalidFieldType:  "invalid type filed",
		ErrDuplicatedRecord:  "duplicated value field",
		ErrRecordNotFound:    "not found",
		ErrSuspended:         "suspend",
		ErrNotMatched:        "not matched",
	}
)

func getErrorMessage(code ErrorCode) ErrorMessage {
	switch code {
	case ErrorInternalError:
		return messages[ErrorInternalError]
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
	case ErrRecordNotFound:
		return messages[ErrRecordNotFound]
	case ErrNotMatched:
		return messages[ErrNotMatched]
	default:
		return messages[ErrorInternalError]
	}
}
