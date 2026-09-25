package sendgrid

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/scrubber"
)

const (
	scopesEndpoint = "/v3/scopes"
)

// sendGridEmailSender defines a struct for sending emails using the SendGrid API
type sendGridEmailSender struct {
	client       *sendgrid.Client
	htmlScrubber scrubber.Scrubber
}

// Option configures a sendGridEmailSender
type Option func(*sendGridEmailSender)

// WithHTMLScrubber sets a scrubber applied to HTML content before sending.
// When set, every outbound message has its HTML sanitized by this scrubber
func WithHTMLScrubber(s scrubber.Scrubber) Option {
	return func(sg *sendGridEmailSender) {
		sg.htmlScrubber = s
	}
}

// New creates a new instance of sendGridEmailSender
func New(apiKey string, opts ...Option) (newman.EmailSender, error) {
	sg := &sendGridEmailSender{
		client: sendgrid.NewSendClient(apiKey),
	}

	for _, opt := range opts {
		opt(sg)
	}

	return sg, nil
}

// SendEmail satisfies the EmailSender interface
func (s *sendGridEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendBatchEmail satisfies the EmailSender interface
func (s *sendGridEmailSender) SendBatchEmail(_ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendBatchEmailWithContext satisfies the EmailSender interface
func (s *sendGridEmailSender) SendBatchEmailWithContext(_ context.Context, _ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *sendGridEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
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
	htmlContent := message.GetHTML()
	if s.htmlScrubber != nil {
		htmlContent = s.htmlScrubber.Scrub(htmlContent)
	}

	if htmlContent != "" {
		v3Mail.AddContent(mail.NewContent("text/html", htmlContent))
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

	response, err := s.client.SendWithContext(ctx, v3Mail)
	if err != nil {
		return ErrFailedToSendEmail
	}

	if response.StatusCode >= http.StatusBadRequest {
		return ErrFailedToSendEmail
	}

	return nil
}

// Verify satisfies the EmailSender interface by listing the scopes granted to the API key
func (s *sendGridEmailSender) Verify(ctx context.Context) error {
	scopesURL, err := url.Parse(s.client.BaseURL)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
	}

	scopesURL.Path = scopesEndpoint

	response, err := sendgrid.MakeRequestWithContext(ctx, rest.Request{
		Method:  rest.Get,
		BaseURL: scopesURL.String(),
		Headers: s.client.Headers,
	})

	switch {
	case err != nil:
		return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
	case response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices:
		return fmt.Errorf("%w: http %d: %s", ErrVerifyFailed, response.StatusCode, response.Body)
	}

	return nil
}
