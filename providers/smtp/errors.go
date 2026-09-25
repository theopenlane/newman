package smtp

import "errors"

var (
	// ErrVerifyFailed is returned when the SMTP server rejects the connection or the configured credentials
	ErrVerifyFailed = errors.New("failed to verify smtp credentials")
)
