package postmark

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/theopenlane/httpsling"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/scrubber"
	"github.com/theopenlane/newman/shared"
)

const (
	requestURL           = "https://api.postmarkapp.com"
	sendEndpoint         = "/email"
	serverEndpoint       = "/server"
	defaultTimeout       = 30 * time.Second
	tokenHeader          = "X-Postmark-Server-Token"
	defaultMessageStream = "outbound"
)

// postmarkEmailSender defines a struct for sending emails using the Postmark API
type postmarkEmailSender struct {
	serverToken  string
	baseURL      string
	client       *http.Client
	requester    *httpsling.Requester
	htmlScrubber scrubber.Scrubber
	stream       string
	trackOpens   bool
	trackLinks   string
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

// WithBaseURL is an option that allows to set a custom base URL for the Postmark API
func WithBaseURL(baseURL string) Option {
	return func(pm *postmarkEmailSender) {
		pm.baseURL = baseURL
	}
}

// WithMessageStream sets the Postmark message stream ID used on every send; an empty stream keeps the default
func WithMessageStream(stream string) Option {
	return func(pm *postmarkEmailSender) {
		if stream != "" {
			pm.stream = stream
		}
	}
}

// WithTrackOpens enables Postmark open tracking on every send
func WithTrackOpens() Option {
	return func(pm *postmarkEmailSender) {
		pm.trackOpens = true
	}
}

// WithTrackLinks sets the Postmark link tracking mode used on every send
func WithTrackLinks(mode string) Option {
	return func(pm *postmarkEmailSender) {
		pm.trackLinks = mode
	}
}

// email represents an email for Postmark
type email struct {
	From          string            `json:"From"`
	To            string            `json:"To"`
	CC            string            `json:"Cc,omitempty"`
	Subject       string            `json:"Subject"`
	TextBody      string            `json:"TextBody,omitempty"`
	HTMLBody      string            `json:"HtmlBody,omitempty"`
	ReplyTo       string            `json:"ReplyTo,omitempty"`
	Bcc           string            `json:"Bcc,omitempty"`
	Headers       []header          `json:"Headers,omitempty"`
	Metadata      map[string]string `json:"Metadata,omitempty"`
	Attachments   []attachment      `json:"Attachments,omitempty"`
	MessageStream string            `json:"MessageStream"`
	TrackOpens    bool              `json:"TrackOpens,omitempty"`
	TrackLinks    string            `json:"TrackLinks,omitempty"`
}

// header represents a custom email header for a Postmark email
type header struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

// attachment represents an attachment for a Postmark email
type attachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
}

// apiResponse represents the error fields of a Postmark response body
type apiResponse struct {
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}

// New creates a new instance of postmarkEmailSender
func New(serverToken string, opts ...Option) (newman.EmailSender, error) {
	if serverToken == "" {
		return nil, ErrMissingServerToken
	}

	pm := &postmarkEmailSender{
		serverToken: serverToken,
		baseURL:     requestURL,
		client:      &http.Client{Timeout: defaultTimeout},
		stream:      defaultMessageStream,
	}

	for _, opt := range opts {
		opt(pm)
	}

	requester, err := httpsling.New(
		httpsling.WithHTTPClient(pm.client),
		httpsling.URL(pm.baseURL),
		httpsling.Header(tokenHeader, pm.serverToken),
		httpsling.Accept(httpsling.ContentTypeJSON),
		httpsling.WithUnmarshaler(&httpsling.JSONMarshaler{}),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateHTTPRequest, err)
	}

	pm.requester = requester

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
func (s *postmarkEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	if err := shared.ValidateEmailMessage(message); err != nil {
		return err
	}

	htmlContent := message.GetHTML()
	if s.htmlScrubber != nil {
		htmlContent = s.htmlScrubber.Scrub(htmlContent)
	}

	emailStruct := email{
		From:          message.GetFrom(),
		To:            strings.Join(message.GetTo(), ","),
		CC:            strings.Join(message.GetCC(), ","),
		Subject:       message.GetSubject(),
		TextBody:      message.GetText(),
		HTMLBody:      htmlContent,
		ReplyTo:       message.GetReplyTo(),
		Bcc:           strings.Join(message.GetBCC(), ","),
		Headers:       make([]header, 0, len(message.Headers)),
		Metadata:      make(map[string]string, len(message.Tags)),
		MessageStream: s.stream,
		TrackOpens:    s.trackOpens,
		TrackLinks:    s.trackLinks,
	}

	for _, name := range slices.Sorted(maps.Keys(message.Headers)) {
		emailStruct.Headers = append(emailStruct.Headers, header{Name: name, Value: message.Headers[name]})
	}

	for _, tag := range message.Tags {
		if tag.Name == "" {
			continue
		}

		emailStruct.Metadata[tag.Name] = tag.Value
	}

	for _, a := range message.GetAttachments() {
		emailStruct.Attachments = append(emailStruct.Attachments, attachment{
			Name:        a.GetFilename(),
			Content:     a.GetBase64StringContent(),
			ContentType: cmp.Or(newman.GetMimeType(a.GetFilename()), httpsling.ContentTypeApplicationOctetStream),
		})
	}

	var response apiResponse

	resp, err := s.requester.ReceiveWithContext(ctx, &response, httpsling.Post(sendEndpoint), httpsling.Body(emailStruct))
	if resp == nil {
		return fmt.Errorf("%w: %w", ErrFailedToSendEmail, err)
	}

	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusTooManyRequests:
		return newman.NewRetryableError(responseError(ErrFailedToSendEmail, resp, response))
	case !httpsling.IsSuccess(resp), response.ErrorCode != 0:
		return responseError(ErrFailedToSendEmail, resp, response)
	}

	return nil
}

// Verify satisfies the EmailSender interface by fetching the server the token belongs to
func (s *postmarkEmailSender) Verify(ctx context.Context) error {
	var response apiResponse

	resp, err := s.requester.ReceiveWithContext(ctx, &response, httpsling.Get(serverEndpoint))
	if resp == nil {
		return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
	}

	defer resp.Body.Close()

	if !httpsling.IsSuccess(resp) || response.ErrorCode != 0 {
		return responseError(ErrVerifyFailed, resp, response)
	}

	return nil
}

// responseError wraps sentinel with the Postmark status, error code, and message, using the raw body when no message was decoded
func responseError(sentinel error, resp *http.Response, response apiResponse) error {
	message := response.Message
	if message == "" {
		if body, err := io.ReadAll(resp.Body); err == nil {
			message = strings.TrimSpace(string(body))
		}
	}

	return fmt.Errorf("%w: http %d: postmark %d: %s", sentinel, resp.StatusCode, response.ErrorCode, message)
}
