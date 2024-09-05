package mailgun

import (
	"testing"

	"github.com/theopenlane/newman"
)

// TestEmailSenderImplementation checks if mailgunEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*mailgunEmailSender)(nil)
}
