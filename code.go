package apperror

type (
	ErrorCode    string
	ErrorMessage string
)

func NewErrorMessage(v string) ErrorMessage {
	return ErrorMessage(v)
}

const (
	ErrRequiredField     ErrorCode = "V4003"
	ErrInvalidFieldType  ErrorCode = "V4001"
	ErrInvalidFieldValue ErrorCode = "V4002"

	ErrDuplicatedRecord ErrorCode = "V4009"
	ErrRecordNotFound   ErrorCode = "V4000"

	ErrorInternalError ErrorCode = "S5003"
)

var (
	messages = map[ErrorCode]ErrorMessage{
		ErrorInternalError:   "internal server error. Please contact admin support",
		ErrRequiredField:     "missing required field",
		ErrInvalidFieldType:  "Invalid type field",
		ErrInvalidFieldValue: "Invalid value filed",
		ErrDuplicatedRecord:  "Duplicated value field",
		ErrRecordNotFound:    "Not found value field",
	}
)
