package postmark

import "errors"

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrFailedToCreateHTTPRequest is returned when an http request fails to be created
	ErrFailedToCreateHTTPRequest = errors.New("failed to create http request")
	// ErrFailedToMarshallEmailData is returned when email data fails to be marshalled
	ErrFailedToMarshallEmailData = errors.New("failed to marshall email data")
)
