package mailgun

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/theopenlane/newman"
)

type MailgunEmailSender struct {
	client mailgun.Mailgun
}

type MGOption func(*MailgunEmailSender)

// WithEurope sets the API Mailgun base url to Europe region.
func WithEurope() MGOption {
	return func(m *MailgunEmailSender) {
		m.client.SetAPIBase(mailgun.APIBaseEU)
	}
}

func NewMailgunEmailSender(domain, apikey string, opts ...MGOption) *MailgunEmailSender {
	mg := &MailgunEmailSender{
		client: mailgun.NewMailgun(domain, apikey),
	}

	for _, opt := range opts {
		opt(mg)
	}

	return mg
}

// SendEmail sends an email using the Mailgun API
func (s *MailgunEmailSender) SendEmail(ctx context.Context, message *newman.EmailMessage) error {
	mailMessage := s.client.NewMessage(message.From, message.Subject, message.Text, message.To...)

	_, _, err := s.client.Send(ctx, mailMessage)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	return nil
}
