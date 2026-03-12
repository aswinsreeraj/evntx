package errors

import (
	"fmt"
)

// AppError is a custom error type that implements the error interface.
// It holds an HTTP status code, an application-specific code, a user-facing message,
// and optionally wraps an underlying standard error.
type AppError struct {
	HTTPCode int
	Code     string
	Message  string
	Err      error
}

// Ensure AppError implements the standard error interface.
var _ error = (*AppError)(nil)

// Error returns a string representation of the AppError.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap allows errors.Is and errors.As to work with the underlying error.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError without an underlying error.
func New(httpCode int, code, message string) *AppError {
	return &AppError{
		HTTPCode: httpCode,
		Code:     code,
		Message:  message,
	}
}

// Wrap creates a new AppError by wrapping an existing standard error.
func Wrap(err error, httpCode int, code, message string) *AppError {
	return &AppError{
		HTTPCode: httpCode,
		Code:     code,
		Message:  message,
		Err:      err,
	}
}
