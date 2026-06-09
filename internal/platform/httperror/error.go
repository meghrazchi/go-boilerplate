package httperror

import "fmt"

type Error struct {
	StatusCode int
	Message    string
	Err        error
	Errors     any
}

func New(statusCode int, message string, err error) *Error {
	return &Error{StatusCode: statusCode, Message: message, Err: err}
}

func (e *Error) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}
