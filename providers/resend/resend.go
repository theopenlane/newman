package resend

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	"github.com/resend/resend-go/v3"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/providers/mock"
	"github.com/theopenlane/newman/scrubber"
	"github.com/theopenlane/newman/shared"
)

// resendEmailSender represents a type that is responsible for sending email messages using the Resend service
type resendEmailSender struct {
	client             *resend.Client
	testDir            string
	defaultAttachments []*resend.Attachment
	htmlScrubber       scrubber.Scrubber
}

// Option is a type representing a function that modifies a ResendEmailSender
type Option func(*resendEmailSender)

// New is a function that creates a new resend EmailSender instance.
func New(apiKey string, options ...Option) (newman.EmailSender, error) {
	// initialize the resendEmailSender
	s := &resendEmailSender{
		client: resend.NewClient(apiKey),
	}

	// apply the options
	for _, option := range options {
		option(s)
	}

	// if the testDir is set, we will use the mock provider
	if s.testDir != "" {
		return mock.New(s.testDir)
	}

	// ensure there is an API key when using the Resend client
	if s.client.ApiKey == "" {
		return nil, ErrMissingAPIKey
	}

	return s, nil
}

// WithClient is an option that allows to set a custom Resend client
func WithClient(client *resend.Client) Option {
	return func(s *resendEmailSender) {
		s.client = client
	}
}

// WithDevMode routes sends to the mock provider, writing MIME files to the given path
func WithDevMode(path string) Option {
	return func(s *resendEmailSender) {
		s.testDir = path
	}
}

// WithBaseURL is an option that allows to set a custom base URL for the Resend client
func WithBaseURL(baseURL url.URL) Option {
	return func(s *resendEmailSender) {
		s.client.BaseURL = &baseURL
	}
}

// WithUserAgent is an option that allows to set a custom user agent for the Resend client
func WithUserAgent(userAgent string) Option {
	return func(s *resendEmailSender) {
		s.client.UserAgent = userAgent
	}
}

// WithFilePath is an option that allows to set a custom file path for the Resend client
func WithFilepath(filepath string) Option {
	return func(s *resendEmailSender) {
		s.defaultAttachments = append(s.defaultAttachments, &resend.Attachment{
			Path: filepath,
		})
	}
}

// WithAPIKey is an option that allows to set a custom API key for the Resend client
func WithAPIKey(apiKey string) Option {
	return func(s *resendEmailSender) {
		s.client.ApiKey = apiKey
	}
}

// WithHTMLScrubber sets a scrubber applied to HTML content before sending.
// When set, every outbound message has its HTML sanitized by this scrubber
func WithHTMLScrubber(s scrubber.Scrubber) Option {
	return func(r *resendEmailSender) {
		r.htmlScrubber = s
	}
}

// SendEmail satisfies the EmailSender interface
func (s *resendEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendBatchEmail satisfies the EmailSender interface
func (s *resendEmailSender) SendBatchEmail(messages []*newman.EmailMessage) error {
	return s.SendBatchEmailWithContext(context.Background(), messages)
}

// toSendEmailRequest converts a newman EmailMessage to a resend SendEmailRequest.
// Resend's batch API does not support attachments, so withAttachments controls
// whether attachment fields are populated on the request
func (s *resendEmailSender) toSendEmailRequest(message *newman.EmailMessage, withAttachments bool) (*resend.SendEmailRequest, error) {
	if err := shared.ValidateEmailMessage(message); err != nil {
		return nil, err
	}

	htmlContent := message.GetHTML()
	if s.htmlScrubber != nil {
		htmlContent = s.htmlScrubber.Scrub(htmlContent)
	}

	req := &resend.SendEmailRequest{
		From:    message.GetFrom(),
		To:      message.GetTo(),
		Subject: message.GetSubject(),
		Bcc:     message.GetBCC(),
		Cc:      message.GetCC(),
		ReplyTo: message.GetReplyTo(),
		Html:    htmlContent,
		Text:    message.GetText(),
		Tags:    make([]resend.Tag, 0, len(message.Tags)),
		Headers: maps.Clone(message.Headers),
	}

	if withAttachments {
		req.Attachments = make([]*resend.Attachment, 0, len(message.Attachments))

		for _, attachment := range message.Attachments {
			req.Attachments = append(req.Attachments, &resend.Attachment{
				Content:     attachment.Content,
				Filename:    attachment.Filename,
				Path:        attachment.FilePath,
				ContentType: attachment.ContentType,
			})
		}

		req.Attachments = append(req.Attachments, slices.Clone(s.defaultAttachments)...)
	}

	for _, tag := range message.Tags {
		req.Tags = append(req.Tags, resend.Tag{
			Name:  tag.Name,
			Value: tag.Value,
		})
	}

	return req, nil
}

// handleSendError normalizes resend API errors into sentinel or retryable errors
func handleSendError(err error, sentinel error) error {
	if strings.Contains(strings.ToLower(err.Error()), "too many requests") {
		return newman.NewRetryableError(err)
	}

	if strings.Contains(err.Error(), "use our testing email address") {
		return nil
	}

	return fmt.Errorf("%w: %w", sentinel, err)
}

// SendBatchEmailWithContext satisfies the EmailSender interface
func (s *resendEmailSender) SendBatchEmailWithContext(ctx context.Context, messages []*newman.EmailMessage) error {
	if len(messages) == 0 {
		return ErrEmptyBatch
	}

	requests := make([]*resend.SendEmailRequest, 0, len(messages))

	for _, message := range messages {
		req, err := s.toSendEmailRequest(message, false)
		if err != nil {
			return err
		}

		requests = append(requests, req)
	}

	_, err := s.client.Batch.SendWithContext(ctx, requests)
	if err != nil {
		return handleSendError(err, ErrFailedToSendBatchEmail)
	}

	return nil
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *resendEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	req, err := s.toSendEmailRequest(message, true)
	if err != nil {
		return err
	}

	_, err = s.client.Emails.SendWithContext(ctx, req)
	if err != nil {
		return handleSendError(err, ErrFailedToSendEmail)
	}

	return nil
}
