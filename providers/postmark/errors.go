package postmark

import "errors"

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrFailedToCreateHTTPRequest is returned when an http request fails to be created
	ErrFailedToCreateHTTPRequest = errors.New("failed to create http request")
	// ErrFailedToMarshallEmailData is returned when email data fails to be marshalled
	ErrFailedToMarshallEmailData = errors.New("failed to marshall email data")
	// ErrMissingServerToken is returned when the Postmark server token is empty
	ErrMissingServerToken = errors.New("missing server token")
	// ErrVerifyFailed is returned when Postmark rejects the configured server token
	ErrVerifyFailed = errors.New("failed to verify postmark credentials")
)
