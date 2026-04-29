package errors

import (
	"fmt"
)

type AppError struct {
	HTTPCode	int
	Code		string
	Message		string
	Err		error
}

var _ error = (*AppError)(nil)

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(httpCode int, code, message string) *AppError {
	return &AppError{
		HTTPCode:	httpCode,
		Code:		code,
		Message:	message,
	}
}

func Wrap(err error, httpCode int, code, message string) *AppError {
	return &AppError{
		HTTPCode:	httpCode,
		Code:		code,
		Message:	message,
		Err:		err,
	}
}
