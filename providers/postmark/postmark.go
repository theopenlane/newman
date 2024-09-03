package postmark

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/theopenlane/newman"
)

const (
	postMarkRequestMethod = "POST"
	postMarkRequestURL    = "https://api.postmarkapp.com/email"
	clientTimeOut         = time.Millisecond * 100
)

// postmarkEmailSender defines a struct for sending emails using the Postmark API
type postmarkEmailSender struct {
	serverToken   string
	requestMethod string
	url           string
}

// NewPostmarkEmailSender creates a new instance of postmarkEmailSender
func NewPostmarkEmailSender(serverToken string) (*postmarkEmailSender, error) {
	return &postmarkEmailSender{serverToken: serverToken, requestMethod: postMarkRequestMethod, url: postMarkRequestURL}, nil
}

// SendEmail sends an email using the Postmark API
func (s *postmarkEmailSender) SendEmail(message *newman.EmailMessage) error {
	emailStruct := struct {
		From        string               `json:"From"`
		To          string               `json:"To"`
		CC          string               `json:"Cc,omitempty"`
		Subject     string               `json:"Subject"`
		TextBody    string               `json:"TextBody,omitempty"`
		HTMLBody    string               `json:"HTMLBody,omitempty"`
		ReplyTo     string               `json:"ReplyTo,omitempty"`
		Bcc         string               `json:"Bcc,omitempty"`
		Attachments []postmarkAttachment `json:"Attachments,omitempty"`
	}{
		From:     message.GetFrom(),
		To:       strings.Join(message.GetTo(), ","),
		CC:       strings.Join(message.GetCC(), ","),
		Subject:  message.GetSubject(),
		TextBody: message.GetText(),
		HTMLBody: message.GetHTML(),
		ReplyTo:  message.GetReplyTo(),
		Bcc:      strings.Join(message.GetBCC(), ","),
	}

	// Add attachments
	for _, attachment := range message.GetAttachments() {
		emailStruct.Attachments = append(emailStruct.Attachments, postmarkAttachment{
			Name:        attachment.GetFilename(),
			Content:     attachment.GetBase64StringContent(),
			ContentType: newman.GetMimeType(attachment.GetFilename()),
		})
	}

	jsonData, err := json.Marshal(emailStruct)
	if err != nil {
		return ErrFailedToMarshallEmailData
	}

	req, err := http.NewRequestWithContext(context.Background(), s.requestMethod, s.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return ErrFailedToCreateHTTPRequest
	}

	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Postmark-Server-Token", s.serverToken)

	client := &http.Client{Timeout: clientTimeOut}

	resp, err := client.Do(req)
	if err != nil {
		return ErrFailedToSendEmail
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrFailedToSendEmail
	}

	return nil
}

// postmarkAttachment represents an attachment for a Postmark email.
type postmarkAttachment struct {
	Name        string `json:"Name"`
	Content     string `json:"Content"`
	ContentType string `json:"ContentType"`
}
