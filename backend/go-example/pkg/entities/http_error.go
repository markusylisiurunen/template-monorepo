package entities

type HttpError struct {
	Code    string
	Message string
}

func NewHttpError(code string, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

var (
	HttpErrorUnknown        = NewHttpError("UNKNOWN", "Unknown error")
	HttpErrorInvalidPayload = NewHttpError("INVALID_PAYLOAD", "Invalid payload")
)
