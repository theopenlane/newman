package sendgrid

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/theopenlane/newman"
)

// sendGridEmailSender defines a struct for sending emails using the SendGrid API
type sendGridEmailSender struct {
	client *sendgrid.Client
}

// NewSendGridEmailSender creates a new instance of sendGridEmailSender
func NewSendGridEmailSender(apiKey string) (*sendGridEmailSender, error) {
	return &sendGridEmailSender{
		client: sendgrid.NewSendClient(apiKey),
	}, nil
}

// SendEmail sends an email using the SendGrid API
func (s *sendGridEmailSender) SendEmail(message *newman.EmailMessage) error {
	from := mail.NewEmail("", message.GetFrom())
	toRecipients := []*mail.Email{}

	for _, to := range message.GetTo() {
		toRecipients = append(toRecipients, mail.NewEmail("", to))
	}

	v3Mail := mail.NewV3Mail()
	v3Mail.SetFrom(from)
	v3Mail.Subject = message.GetSubject()

	// Create personalization for To recipients
	personalization := mail.NewPersonalization()
	for _, to := range toRecipients {
		personalization.AddTos(to)
	}

	// Add BCC recipients
	for _, bcc := range message.GetBCC() {
		personalization.AddBCCs(mail.NewEmail("", bcc))
	}

	// Add Reply-To if specified
	if message.GetReplyTo() != "" {
		replyTo := mail.NewEmail("", message.GetReplyTo())
		v3Mail.SetReplyTo(replyTo)
	}

	v3Mail.AddPersonalizations(personalization)

	// Add plain text content
	if message.GetText() != "" {
		v3Mail.AddContent(mail.NewContent("text/plain", message.GetText()))
	}

	// Add HTML content
	if message.GetHTML() != "" {
		v3Mail.AddContent(mail.NewContent("text/html", message.GetHTML()))
	}

	// Add attachments
	for _, attachment := range message.GetAttachments() {
		a := mail.NewAttachment()
		a.SetContent(attachment.GetBase64StringContent())
		a.SetType(newman.GetMimeType(attachment.GetFilename()))
		a.SetFilename(attachment.GetFilename())
		a.SetDisposition("attachment")
		v3Mail.AddAttachment(a)
	}

	response, err := s.client.Send(v3Mail)
	if err != nil {
		return ErrFailedToSendEmail
	}

	if response.StatusCode >= 400 { // nolint: mnd
		return ErrFailedToSendEmail
	}

	return nil
}
