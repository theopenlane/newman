package resend

import (
	"errors"
)

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
)
