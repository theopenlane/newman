package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/theopenlane/newman/scrubber"
)

const DefaultMaxAttachmentSize = 25 * 1024 * 1024 // 25 MB

// EmailMessage contains the fields for sending an email
type EmailMessage struct {
	// From is the email address of the sender
	From string `json:"from"`
	// To is the email address of the recipient
	To []string `json:"to"`
	// Subject is the subject of the email
	Subject string `json:"subject"`
	// Bcc is the email address of the blind carbon copy recipient
	Bcc []string `json:"bcc,omitempty"`
	// Cc is the email address of the carbon copy recipient
	Cc []string `json:"cc,omitempty"`
	// ReplyTo is the email address to reply to
	ReplyTo string `json:"reply_to,omitempty"`
	// HTML is the HTML content of the email
	HTML string `json:"html,omitempty"`
	// Text is the text content of the email
	Text string `json:"text,omitempty"`
	// Tags is the list of tags associated with the email
	Tags []Tag `json:"tags,omitempty"`
	// Attachments is the list of attachments associated with the email
	Attachments []*Attachment `json:"attachments,omitempty"`
	// Headers is the list of headers associated with the email
	Headers map[string]string `json:"headers,omitempty"`
	// Maximum size for attachments
	maxAttachmentSize int
	// textScrubber is used to scrub text content
	textScrubber scrubber.Scrubber
	// htmlScrubber is used to scrub HTML content
	htmlScrubber scrubber.Scrubber
}

// Tag is used to define custom metadata for message
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewEmailMessage creates a new EmailMessage with the required fields
func NewEmailMessage(from string, to []string, subject string, body string) *EmailMessage {
	email := &EmailMessage{
		From:              from,
		To:                to,
		Subject:           subject,
		maxAttachmentSize: DefaultMaxAttachmentSize,
	}

	if IsHTML(body) {
		email.HTML = body
	} else {
		email.Text = body
	}

	return email
}

// NewFullEmailMessage creates a new EmailMessage with all fields
func NewFullEmailMessage(from string, to []string, subject string, cc []string, bcc []string, replyTo string, textBody string, htmlBody string, attachments []*Attachment) *EmailMessage {
	return &EmailMessage{
		From:              from,
		To:                to,
		Cc:                cc,
		Bcc:               bcc,
		ReplyTo:           replyTo,
		Subject:           subject,
		Text:              textBody,
		HTML:              htmlBody,
		Attachments:       attachments,
		maxAttachmentSize: DefaultMaxAttachmentSize,
	}
}

// SetFrom sets the sender email address
func (e *EmailMessage) SetFrom(from string) *EmailMessage {
	e.From = from
	return e
}

// SetSubject sets the email subject
func (e *EmailMessage) SetSubject(subject string) *EmailMessage {
	e.Subject = subject
	return e
}

// SetTo sets the recipient email addresses
func (e *EmailMessage) SetTo(to []string) *EmailMessage {
	e.To = to
	return e
}

// SetCC sets the CC recipients email addresses
func (e *EmailMessage) SetCC(cc []string) *EmailMessage {
	e.Cc = cc
	return e
}

// SetBCC sets the BCC recipients email addresses
func (e *EmailMessage) SetBCC(bcc []string) *EmailMessage {
	e.Bcc = bcc
	return e
}

// SetReplyTo sets the reply-to email address
func (e *EmailMessage) SetReplyTo(replyTo string) *EmailMessage {
	e.ReplyTo = replyTo
	return e
}

// SetText sets the plain text content of the email
func (e *EmailMessage) SetText(text string) *EmailMessage {
	e.Text = text
	return e
}

// SetHTML sets the HTML content of the email
func (e *EmailMessage) SetHTML(html string) *EmailMessage {
	e.HTML = html
	return e
}

// SetAttachments sets the attachments for the email
func (e *EmailMessage) SetAttachments(attachments []*Attachment) *EmailMessage {
	e.Attachments = attachments
	return e
}

// AddToRecipient adds a recipient email address to the To field
func (e *EmailMessage) AddToRecipient(recipient string) *EmailMessage {
	e.To = append(e.To, recipient)
	return e
}

// AddCCRecipient adds a recipient email address to the CC field
func (e *EmailMessage) AddCCRecipient(recipient string) *EmailMessage {
	e.Cc = append(e.Cc, recipient)
	return e
}

// AddBCCRecipient adds a recipient email address to the BCC field
func (e *EmailMessage) AddBCCRecipient(recipient string) *EmailMessage {
	e.Bcc = append(e.Bcc, recipient)
	return e
}

// AddAttachment adds an attachment to the email
func (e *EmailMessage) AddAttachment(attachment *Attachment) *EmailMessage {
	e.Attachments = append(e.Attachments, attachment)
	return e
}

// GetFrom returns the trimmed and validated sender email address
func (e *EmailMessage) GetFrom() string {
	if e == nil {
		return ""
	}

	return ValidateEmail(e.From)
}

// GetTo returns a slice of trimmed and validated recipient email addresses
func (e *EmailMessage) GetTo() []string {
	if e == nil {
		return []string{}
	}

	return ValidateEmailSlice(e.To)
}

// GetCC returns a slice of trimmed and validated CC recipient email addresses
func (e *EmailMessage) GetCC() []string {
	if e == nil {
		return []string{}
	}

	return ValidateEmailSlice(e.Cc)
}

// GetBCC returns a slice of trimmed and validated BCC recipient email addresses
func (e *EmailMessage) GetBCC() []string {
	if e == nil {
		return []string{}
	}

	return ValidateEmailSlice(e.Bcc)
}

// GetReplyTo returns the trimmed and validated reply-to email address
func (e *EmailMessage) GetReplyTo() string {
	if e == nil {
		return ""
	}

	return ValidateEmail(e.ReplyTo)
}

// GetSubject returns the scrubd email subject
func (e *EmailMessage) GetSubject() string {
	if e == nil {
		return ""
	}

	if e.textScrubber != nil {
		return e.textScrubber.Scrub(e.Subject)
	}

	return scrubber.DefaultTextScrubber().Scrub(e.Subject)
}

// GetText returns the scrubd plain text content of the email
func (e *EmailMessage) GetText() string {
	if e == nil {
		return ""
	}

	if e.textScrubber != nil {
		return e.textScrubber.Scrub(e.Text)
	}

	return scrubber.DefaultTextScrubber().Scrub(e.Text)
}

// GetHTML returns the scrubd HTML content of the email
func (e *EmailMessage) GetHTML() string {
	if e == nil {
		return ""
	}

	if e.htmlScrubber != nil {
		return e.htmlScrubber.Scrub(e.HTML)
	}

	return scrubber.DefaultHTMLScrubber().Scrub(e.HTML)
}

// SetMaxAttachmentSize sets the maximum attachment size
func (e *EmailMessage) SetMaxAttachmentSize(size int) *EmailMessage {
	e.maxAttachmentSize = size
	return e
}

// GetAttachments returns the attachments to be included in the email, filtering out those that exceed the maximum size
func (e *EmailMessage) GetAttachments() []*Attachment {
	if e == nil {
		return []*Attachment{}
	}

	if e.maxAttachmentSize < 0 {
		return e.Attachments
	}

	var validAttachments []*Attachment

	for _, attachment := range e.Attachments {
		if len(attachment.Content) <= e.maxAttachmentSize {
			validAttachments = append(validAttachments, attachment)
		}
	}

	return validAttachments
}

// SetCustomTextScrubber sets a custom scrubber for text content
func (e *EmailMessage) SetCustomTextScrubber(s scrubber.Scrubber) *EmailMessage {
	e.textScrubber = s

	return e
}

// SetCustomHTMLScrubber sets a custom scrubber for HTML content
func (e *EmailMessage) SetCustomHTMLScrubber(s scrubber.Scrubber) *EmailMessage {
	e.htmlScrubber = s

	return e
}

// jsonEmailMessage represents the JSON structure for an email message.
type jsonEmailMessage struct {
	From        string        `json:"from"`
	To          []string      `json:"to"`
	CC          []string      `json:"cc,omitempty"`
	BCC         []string      `json:"bcc,omitempty"`
	ReplyTo     string        `json:"replyTo,omitempty"`
	Subject     string        `json:"subject"`
	Text        string        `json:"text"`
	HTML        string        `json:"html,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
}

// MarshalJSON is a custom marshaler for EmailMessage
func (e *EmailMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonEmailMessage{
		From:        e.From,
		To:          e.To,
		CC:          e.Cc,
		BCC:         e.Bcc,
		ReplyTo:     e.ReplyTo,
		Subject:     e.Subject,
		Text:        e.Text,
		HTML:        e.HTML,
		Attachments: e.Attachments,
	})
}

// UnmarshalJSON is a custom unmarshaler for EmailMessage
func (e *EmailMessage) UnmarshalJSON(data []byte) error {
	aux := &jsonEmailMessage{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	e.maxAttachmentSize = DefaultMaxAttachmentSize

	e.From = aux.From
	e.To = aux.To
	e.Cc = aux.CC
	e.Bcc = aux.BCC
	e.ReplyTo = aux.ReplyTo
	e.Subject = aux.Subject
	e.Text = aux.Text
	e.HTML = aux.HTML
	e.Attachments = aux.Attachments

	return nil
}

// BuildMimeMessage constructs the MIME message for the email, including text, HTML, and attachments
func BuildMimeMessage(message *EmailMessage) ([]byte, error) {
	var msg bytes.Buffer

	// Determine boundaries
	mixedBoundary := fmt.Sprintf("mixed-boundary-%d", time.Now().UnixNano())
	altBoundary := fmt.Sprintf("alt-boundary-%d", time.Now().UnixNano())

	// Basic headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", message.GetFrom()))

	// Add To recipients
	toRecipients := message.GetTo()
	if len(toRecipients) > 0 {
		msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(toRecipients, ",")))
	}

	ccRecipients := message.GetCC()

	if len(ccRecipients) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(ccRecipients, ",")))
	}

	if message.GetReplyTo() != "" {
		msg.WriteString(fmt.Sprintf("Reply-To: %s\r\n", message.GetReplyTo()))
	}

	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", message.GetSubject()))

	msg.WriteString("MIME-Version: 1.0\r\n")

	// Use multipart/mixed if there are attachments, otherwise multipart/alternative
	attachments := message.GetAttachments()
	if len(attachments) > 0 {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", mixedBoundary))
		msg.WriteString("\r\n")
		// Start multipart/alternative
		msg.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", altBoundary))
		msg.WriteString("\r\n")
	} else {
		msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", altBoundary))
		msg.WriteString("\r\n")
	}

	// Plain text part
	textMessage := message.GetText()
	if textMessage != "" {
		msg.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(textMessage)
		msg.WriteString("\r\n")
	}

	// HTML part
	htmlMessage := message.GetHTML()
	if htmlMessage != "" {
		msg.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		msg.WriteString("\r\n")
		msg.WriteString(htmlMessage)
		msg.WriteString("\r\n")
	}

	// End multipart/alternative
	msg.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary))

	// Attachments
	if len(attachments) > 0 {
		for _, attachment := range attachments {
			fileName := attachment.GetFilename()

			mimeType := GetMimeType(fileName)

			msg.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
			msg.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
			msg.WriteString("Content-Transfer-Encoding: base64\r\n")
			msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", fileName))
			msg.WriteString("\r\n")
			msg.Write(attachment.GetBase64Content())
			msg.WriteString("\r\n")
		}

		// End multipart/mixed
		msg.WriteString(fmt.Sprintf("--%s--\r\n", mixedBoundary))
	}

	return msg.Bytes(), nil
}
