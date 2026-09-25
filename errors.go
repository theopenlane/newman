package newman

import (
	"errors"
)

// ErrBatchNotImplemented is returned by providers that do not support native batch sending
var ErrBatchNotImplemented = errors.New("batch email sending is not implemented for this provider")

type retryableError struct {
	reason error
}

func (e retryableError) Error() string { return e.reason.Error() }

// Unwrap returns the underlying reason so errors.Is and errors.As can match it
func (e retryableError) Unwrap() error { return e.reason }

// NewRetryableError creates a new retryable error with a given reason.
func NewRetryableError(reason error) error {
	return retryableError{reason: reason}
}

// IsRetryableError checks if the error is retryable.
func IsRetryableError(err error) bool {
	var re retryableError
	return errors.As(err, &re)
}
