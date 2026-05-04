package postmark

import (
	"context"
	"strings"
	"time"

	"github.com/theopenlane/httpsling"
	"github.com/theopenlane/httpsling/httpclient"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/scrubber"
)

const (
	requestURL    = "https://api.postmarkapp.com"
	endpoint      = "/email"
	clientTimeout = time.Millisecond * 100
	tokenHeader   = "X-Postmark-Server-Token"
)

// postmarkEmailSender defines a struct for sending emails using the Postmark API
type postmarkEmailSender struct {
	serverToken  string
	endpoint     string
	url          string
	htmlScrubber scrubber.Scrubber
}

// Option configures a postmarkEmailSender
type Option func(*postmarkEmailSender)

// WithHTMLScrubber sets a scrubber applied to HTML content before sending.
// When set, every outbound message has its HTML sanitized by this scrubber
func WithHTMLScrubber(s scrubber.Scrubber) Option {
	return func(pm *postmarkEmailSender) {
		pm.htmlScrubber = s
	}
}

// email represents an email for Postmark
type email struct {
	From        string       `json:"From"`
	To          string       `json:"To"`
	CC          string       `json:"Cc,omitempty"`
	Subject     string       `json:"Subject"`
	TextBody    string       `json:"TextBody,omitempty"`
	HTMLBody    string       `json:"HTMLBody,omitempty"`
	ReplyTo     string       `json:"ReplyTo,omitempty"`
	Bcc         string       `json:"Bcc,omitempty"`
	Attachments []attachment `json:"Attachments,omitempty"`
}

// attachment represents an attachment for a Postmark email
type attachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
}

// New creates a new instance of postmarkEmailSender
func New(serverToken string, opts ...Option) (newman.EmailSender, error) {
	pm := &postmarkEmailSender{
		serverToken: serverToken,
		endpoint:    endpoint,
		url:         requestURL,
	}

	for _, opt := range opts {
		opt(pm)
	}

	return pm, nil
}

// SendEmail satisfies the EmailSender interface
func (s *postmarkEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendBatchEmail satisfies the EmailSender interface
func (s *postmarkEmailSender) SendBatchEmail(_ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendBatchEmailWithContext satisfies the EmailSender interface
func (s *postmarkEmailSender) SendBatchEmailWithContext(_ context.Context, _ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *postmarkEmailSender) SendEmailWithContext(_ context.Context, message *newman.EmailMessage) error {
	htmlContent := message.GetHTML()
	if s.htmlScrubber != nil {
		htmlContent = s.htmlScrubber.Scrub(htmlContent)
	}

	emailStruct := email{
		From:     message.GetFrom(),
		To:       strings.Join(message.GetTo(), ","),
		CC:       strings.Join(message.GetCC(), ","),
		Subject:  message.GetSubject(),
		TextBody: message.GetText(),
		HTMLBody: htmlContent,
		ReplyTo:  message.GetReplyTo(),
		Bcc:      strings.Join(message.GetBCC(), ","),
	}

	// Add attachments
	for _, a := range message.GetAttachments() {
		emailStruct.Attachments = append(emailStruct.Attachments, attachment{
			Name:        a.GetFilename(),
			Content:     a.GetBase64StringContent(),
			ContentType: newman.GetMimeType(a.GetFilename()),
		})
	}

	requester, err := httpsling.New(
		httpsling.Client(httpclient.Timeout(clientTimeout)),
		httpsling.URL(s.url),
		httpsling.Header(tokenHeader, s.serverToken),
	)
	if err != nil {
		return err
	}

	resp, err := requester.Receive(
		httpsling.Post(s.endpoint),
		httpsling.Body(emailStruct),
	)
	if err != nil {
		return ErrFailedToSendEmail
	}

	defer resp.Body.Close()

	if !httpsling.IsSuccess(resp) {
		return ErrFailedToSendEmail
	}

	return nil
}
