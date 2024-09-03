package shared

// Option is a type representing a function that modifies a ResendClient
type MessageOption func(*EmailMessage)

// NewResendClient is a function that creates a new ResendClient instance.
func NewEmailMessageWithOptions(options ...MessageOption) *EmailMessage {
	s := EmailMessage{}

	for _, option := range options {
		option(&s)
	}

	return &s
}

// WithFrom sets the from email address
func WithFrom(from string) MessageOption {
	return func(m *EmailMessage) {
		m.From = from
	}
}

// WithTo sets the to email address
func WithTo(to []string) MessageOption {
	return func(m *EmailMessage) {
		m.To = to
	}
}

// WithSubject sets the subject of the email
func WithSubject(subject string) MessageOption {
	return func(m *EmailMessage) {
		m.Subject = subject
	}
}

// WithBcc sets the bcc email address
func WithBcc(bcc []string) MessageOption {
	return func(m *EmailMessage) {
		m.Bcc = bcc
	}
}

// WithCc sets the cc email address
func WithCc(cc []string) MessageOption {
	return func(m *EmailMessage) {
		m.Cc = cc
	}
}

// WithReplyTo sets the reply to email address
func WithReplyTo(replyTo string) MessageOption {
	return func(m *EmailMessage) {
		m.ReplyTo = replyTo
	}
}

// WithHTML sets the html content of the email
func WithHTML(html string) MessageOption {
	return func(m *EmailMessage) {
		m.HTML = html
	}
}

// WithText sets the text content of the email
func WithText(text string) MessageOption {
	return func(m *EmailMessage) {
		m.Text = text
	}
}

// WithTag adds a tag to the email
func WithTag(tag Tag) MessageOption {
	return func(m *EmailMessage) {
		m.Tags = append(m.Tags, tag)
	}
}

// WithTags sets the tags of the email
func WithTags(tags []Tag) MessageOption {
	return func(m *EmailMessage) {
		m.Tags = tags
	}
}

// WithAttachment adds an attachment to the email
func WithAttachment(attachment *Attachment) MessageOption {
	return func(m *EmailMessage) {
		m.Attachments = append(m.Attachments, attachment)
	}
}

// WithAttachments sets the attachments of the email
func WithAttachments(attachments []*Attachment) MessageOption {
	return func(m *EmailMessage) {
		m.Attachments = attachments
	}
}

// WithHeader adds a header to the email
func WithHeader(key, value string) MessageOption {
	return func(m *EmailMessage) {
		m.Headers[key] = value
	}
}

// WithHeaders sets the headers of the email
func WithHeaders(headers map[string]string) MessageOption {
	return func(m *EmailMessage) {
		m.Headers = headers
	}
}

// WithHeaderMap adds a map of headers to the email
func WithHeaderMap(headers map[string]string) MessageOption {
	return func(m *EmailMessage) {
		for key, value := range headers {
			m.Headers[key] = value
		}
	}
}
