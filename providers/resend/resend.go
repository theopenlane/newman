package resend

import (
	"context"
	"maps"
	"net/url"
	"slices"

	"github.com/resend/resend-go/v2"

	"github.com/theopenlane/newman"
)

// ResendEmailSender represents a type that is responsible for sending email messages using the Resend service
type ResendEmailSender struct {
	client *resend.Client
}

// ResendOption is a type representing a function that modifies a ResendEmailSender
type ResendOption func(*ResendEmailSender)

// NewResendEmailSender is a function that creates a new ResendEmailSender instance.
func NewResendEmailSender(apikey string, options ...ResendOption) *ResendEmailSender {
	s := ResendEmailSender{
		client: resend.NewClient(apikey),
	}

	for _, option := range options {
		option(&s)
	}

	return &s
}

// WithClient is an option that allows to set a custom Resend client
func WithClient(client *resend.Client) ResendOption {
	return func(s *ResendEmailSender) {
		s.client = client
	}
}

// WithBaseURL is an option that allows to set a custom base URL for the Resend client
func WithBaseURL(baseURL url.URL) ResendOption {
	return func(s *ResendEmailSender) {
		s.client.BaseURL = &baseURL
	}
}

// WithUserAgent is an option that allows to set a custom user agent for the Resend client
func WithUserAgent(userAgent string) ResendOption {
	return func(s *ResendEmailSender) {
		s.client.UserAgent = userAgent
	}
}

// WithFilePath is an option that allows to set a custom file path for the Resend client
func WithFilepath(filepath string) ResendOption {
	return func(s *ResendEmailSender) {
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
func WithAPIKey(apikey string) ResendOption {
	return func(s *ResendEmailSender) {
		s.client.ApiKey = apikey
	}
}

// WithHeaders is an option that allows to set a custom headers for the Resend client
func WithHeaders(headers map[string]string) ResendOption {
	return func(s *ResendEmailSender) {
		func() *resend.SendEmailRequest {
			return &resend.SendEmailRequest{
				Headers: maps.Clone(headers),
			}
		}()
	}
}

// SendEmail sends an email message using the Resend service
func (s *ResendEmailSender) SendEmail(ctx context.Context, message newman.EmailMessage) error {
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
		return ErrFailedToSendEmail
	}

	return nil
}
