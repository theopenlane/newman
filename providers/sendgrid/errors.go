package sendgrid

import "errors"

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrVerifyFailed is returned when SendGrid rejects the configured API key
	ErrVerifyFailed = errors.New("failed to verify sendgrid credentials")
)
