package resend

import (
	"errors"
)

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrFailedToSendBatchEmail is returned when a batch email fails to send
	ErrFailedToSendBatchEmail = errors.New("failed to send batch email")
	// ErrMissingAPIKey is returned when an API key is missing
	ErrMissingAPIKey = errors.New("missing API key")
	// ErrEmptyBatch is returned when an empty batch is provided
	ErrEmptyBatch = errors.New("batch must contain at least one message")
)
