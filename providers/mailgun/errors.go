package mailgun

import "errors"

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrMissingAPIKey is returned when an API key is missing
	ErrMissingAPIKey = errors.New("missing API key")
	// ErrDomainMissing is returned when no sending domain is configured
	ErrDomainMissing = errors.New("missing domain")
	// ErrVerifyFailed is returned when Mailgun rejects the configured API key or domain
	ErrVerifyFailed = errors.New("failed to verify mailgun credentials")
)
