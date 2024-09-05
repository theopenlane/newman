package newman

import (
	"context"

	"github.com/theopenlane/newman/shared"
)

// EmailSender interface defines the method to send an email
type EmailSender interface {
	// SendEmail sends an email with the given message
	SendEmail(message *EmailMessage) error
	// SendEmailWithContext sends an email with the given message and context
	SendEmailWithContext(ctx context.Context, message *EmailMessage) error
}

// EmailMessage represents an email message
type EmailMessage = shared.EmailMessage

// Attachment represents an email attachment with its filename and content
type Attachment = shared.Attachment

// Tag is used to define custom metadata for message
type Tag = shared.Tag

// NewEmailMessage creates a new EmailMessage with the required fields
func NewEmailMessage(from string, to []string, subject string, body string) *EmailMessage {
	return shared.NewEmailMessage(from, to, subject, body)
}

// NewEmailMessageWithOptions creates a new EmailMessage with the specified options.
func NewEmailMessageWithOptions(options ...MessageOption) *EmailMessage {
	s := EmailMessage{}

	for _, option := range options {
		option(&s)
	}

	return &s
}

// NewAttachment creates a new Attachment instance with the specified filename and content
func NewAttachment(filename string, content []byte) *Attachment {
	return shared.NewAttachment(filename, content)
}

// NewAttachmentFromFile creates a new Attachment instance from the specified file path
func NewAttachmentFromFile(filePath string) (*Attachment, error) {
	return shared.NewAttachmentFromFile(filePath)
}

// BuildMimeMessage constructs the MIME message for the email, including text, HTML, and attachments
func BuildMimeMessage(message *EmailMessage) ([]byte, error) {
	return shared.BuildMimeMessage(message)
}

// ValidateEmail validates and sanitizes an email address
func ValidateEmail(email string) string {
	return shared.ValidateEmailAddress(email)
}

// ValidateEmailSlice validates and sanitizes a slice of email addresses
func ValidateEmailSlice(emails []string) []string {
	return shared.ValidateEmailAddresses(emails)
}

// GetMimeType returns the MIME type based on the file extension
func GetMimeType(filename string) string {
	return shared.GetMimeType(filename)
}
