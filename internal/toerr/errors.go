package toerr

import (
	"errors"
)

type ValidationError struct {
	HttpErrorCode int
	Message       string
}

func NewValidationError(code int, message string) *ValidationError {
	return &ValidationError{code, message}
}

func (v *ValidationError) Error() string {
	return v.Message
}

func ExtractValidationError(err error) *ValidationError {
	verr := &ValidationError{}

	isverr := errors.As(err, &verr)
	if isverr {
		return verr
	}
	return nil
}
