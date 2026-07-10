package apperr

type Code string

const (
	CodeInvalidInput  = "invalid_input"
	CodeNotFound      = "not_found"
	CodeAlreadyExists = "already_exists"
	CodeInternal      = "internal_error"
)

const (
	MessageInvalidInput  = "invalid input"
	MessageNotFound      = "not found"
	MessageAlreadyExists = "already exists"
	MessageInternal      = "internal server error"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewError(code Code, message string, err error) error {
	return &Error{Code: code, Message: message, Err: err}
}

func InvalidInput(message string) error {
	return NewError(CodeInvalidInput, message, nil)
}

func NotFound(message string) error {
	return NewError(CodeNotFound, message, nil)
}

func AlreadyExists(message string) error {
	return NewError(CodeAlreadyExists, message, nil)
}

func Internal(message string, err error) error {
	return NewError(CodeInternal, message, err)
}
