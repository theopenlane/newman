package resend

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"slices"

	"github.com/resend/resend-go/v2"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/providers/mock"
	"github.com/theopenlane/newman/shared"
)

// resendEmailSender represents a type that is responsible for sending email messages using the Resend service
type resendEmailSender struct {
	client  *resend.Client
	testDir string
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
	return func(_ *resendEmailSender) {
		func() *resend.SendEmailRequest {
			return &resend.SendEmailRequest{
				Attachments: []*resend.Attachment{
					{
						Path: filepath,
					},
				},
			}
		}()
	}
}

// WithAPIKey is an option that allows to set a custom API key for the Resend client
func WithAPIKey(apiKey string) Option {
	return func(s *resendEmailSender) {
		s.client.ApiKey = apiKey
	}
}

// WithHeaders is an option that allows to set a custom headers for the Resend client
func WithHeaders(headers map[string]string) Option {
	return func(_ *resendEmailSender) {
		func() *resend.SendEmailRequest {
			return &resend.SendEmailRequest{
				Headers: maps.Clone(headers),
			}
		}()
	}
}

// SendEmail satisfies the EmailSender interface
func (s *resendEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *resendEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	if err := shared.ValidateEmailMessage(message); err != nil {
		return err
	}

	msgToSend := resend.SendEmailRequest{
		From:        message.From,
		To:          slices.Clone(message.To),
		Subject:     message.Subject,
		Bcc:         slices.Clone(message.Bcc),
		Cc:          slices.Clone(message.Cc),
		ReplyTo:     message.ReplyTo,
		Html:        message.HTML,
		Text:        message.Text,
		Tags:        make([]resend.Tag, 0, len(message.Tags)),
		Attachments: make([]*resend.Attachment, 0, len(message.Attachments)),
		Headers:     maps.Clone(message.Headers),
	}

	for _, attachment := range message.Attachments {
		resendAttachment := &resend.Attachment{
			Content:     attachment.Content,
			Filename:    attachment.Filename,
			Path:        attachment.FilePath,
			ContentType: attachment.ContentType,
		}

		msgToSend.Attachments = append(msgToSend.Attachments, resendAttachment)
	}

	for _, tag := range message.Tags {
		resendTag := resend.Tag{
			Name:  tag.Name,
			Value: tag.Value,
		}

		msgToSend.Tags = append(msgToSend.Tags, resendTag)
	}

	if _, err := s.client.Emails.SendWithContext(ctx, &msgToSend); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSendEmail, err)
	}

	return nil
}
