package newman

import (
	"errors"
)

type retryableError struct {
	reason error
}

func (e retryableError) Error() string { return e.reason.Error() }

// NewRetryableError creates a new retryable error with a given reason.
func NewRetryableError(reason error) error {
	return retryableError{reason: reason}
}

// IsRetryableError checks if the error is retryable.
func IsRetryableError(err error) bool {
	var re retryableError
	return errors.As(err, &re)
}
