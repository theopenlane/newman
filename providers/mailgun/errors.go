package mailgun

import "errors"

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrMissingAPIKey is returned when an API key is missing
	ErrMissingAPIKey = errors.New("missing API key")
)
