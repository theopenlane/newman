package mailgun

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/theopenlane/newman"
)

type mailgunEmailSender struct {
	client mailgun.Mailgun
}

// Option is a type representing a function that modifies a mailgunEmailSender
type Option func(*mailgunEmailSender)

// WithEurope sets the API Mailgun base url to Europe region.
func WithEurope() Option {
	return func(m *mailgunEmailSender) {
		m.client.SetAPIBase(mailgun.APIBaseEU)
	}
}

// New creates a new mailgunEmailSender
func New(domain, apiKey string, opts ...Option) (newman.EmailSender, error) {
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}

	mg := &mailgunEmailSender{
		client: mailgun.NewMailgun(domain, apiKey),
	}

	for _, opt := range opts {
		opt(mg)
	}

	return mg, nil
}

// SendEmail satisfies the EmailSender interface
func (s *mailgunEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *mailgunEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	mailMessage := mailgun.NewMessage(message.From, message.Subject, message.Text, message.To...)

	if _, _, err := s.client.Send(ctx, mailMessage); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
