package shared

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/newman/scrubber"
)

func TestEmailMessageGetters(t *testing.T) {
	message := EmailMessage{
		From:    "newman@usps.com",
		To:      []string{"sfunk@funkytown.com", "kwaters@quiddich.com"},
		Cc:      []string{"elaine@seinfeld.com", "kramer@seinfeld.com"},
		Bcc:     []string{"belaine@seinfeld.com", "bkramer@seinfeld.com"},
		ReplyTo: "newman@usps.com",
		Subject: "Look sister, go get yourself a cup of coffee or something",
		Text:    "Test Text",
		HTML:    "<h1>Test HTML</h1>",
		Attachments: []*Attachment{
			{Filename: "test.txt", Content: []byte("test content")},
		},
		maxAttachmentSize: DefaultMaxAttachmentSize,
	}

	t.Run("GetFrom", func(t *testing.T) {
		expected := "newman@usps.com"
		result := message.GetFrom()

		assert.Equal(t, expected, result)
	})

	t.Run("GetReplyTo", func(t *testing.T) {
		expected := "newman@usps.com"
		result := message.GetReplyTo()

		assert.Equal(t, expected, result)
	})

	t.Run("GetTo", func(t *testing.T) {
		expected := []string{"sfunk@funkytown.com", "kwaters@quiddich.com"}
		result := message.GetTo()

		assert.Equal(t, expected, result)
	})

	t.Run("GetCC", func(t *testing.T) {
		expected := []string{"elaine@seinfeld.com", "kramer@seinfeld.com"}
		result := message.GetCC()

		assert.Equal(t, expected, result)
	})

	t.Run("GetBCC", func(t *testing.T) {
		expected := []string{"belaine@seinfeld.com", "bkramer@seinfeld.com"}
		result := message.GetBCC()

		assert.Equal(t, expected, result)
	})

	t.Run("GetSubject", func(t *testing.T) {
		expected := "Look sister, go get yourself a cup of coffee or something"
		result := message.GetSubject()

		assert.Equal(t, expected, result)
	})

	t.Run("GetText", func(t *testing.T) {
		expected := "Test Text"
		result := message.GetText()

		assert.Equal(t, expected, result)
	})

	t.Run("GetHTML", func(t *testing.T) {
		expected := "<h1>Test HTML</h1>"
		result := message.GetHTML()

		assert.Equal(t, expected, result)
	})

	t.Run("GetAttachments", func(t *testing.T) {
		expected := []*Attachment{
			{Filename: "test.txt", Content: []byte("test content")},
		}
		result := message.GetAttachments()

		assert.Equal(t, expected, result)
	})
}

func TestNilEmailMessageGetters(t *testing.T) {
	var message *EmailMessage

	t.Run("GetFrom", func(t *testing.T) {
		result := message.GetFrom()
		assert.Equal(t, "", result)
	})

	t.Run("GetReplyTo", func(t *testing.T) {
		result := message.GetReplyTo()
		assert.Equal(t, "", result)
	})

	t.Run("GetTo", func(t *testing.T) {
		result := message.GetTo()
		assert.Equal(t, []string{}, result)
	})

	t.Run("GetCC", func(t *testing.T) {
		result := message.GetCC()
		assert.Equal(t, []string{}, result)
	})

	t.Run("GetBCC", func(t *testing.T) {
		result := message.GetBCC()
		assert.Equal(t, []string{}, result)
	})

	t.Run("GetSubject", func(t *testing.T) {
		result := message.GetSubject()
		assert.Equal(t, "", result)
	})

	t.Run("GetText", func(t *testing.T) {
		result := message.GetText()
		assert.Equal(t, "", result)
	})

	t.Run("GetHTML", func(t *testing.T) {
		result := message.GetHTML()
		assert.Equal(t, "", result)
	})

	t.Run("GetAttachments", func(t *testing.T) {
		result := message.GetAttachments()
		assert.Equal(t, []*Attachment{}, result)
	})
}

func TestNewEmailMessage(t *testing.T) {
	t.Run("create plain text email", func(t *testing.T) {
		from := "newman@usps.com"
		to := []string{"jerry@seinfeld.com"}
		subject := "Subject"
		body := "Email body"
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, body, email.Text)
		assert.Equal(t, "", email.HTML)
	})

	t.Run("create HTML email", func(t *testing.T) {
		from := "newman@usps.com"
		to := []string{"jerry@seinfeld.com"}
		subject := "Subject"
		body := "<p>Email body</p>"
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, "", email.Text)
		assert.Equal(t, body, email.HTML)
	})
}

func TestNewFullEmailMessage(t *testing.T) {
	t.Run("create full email message", func(t *testing.T) {
		from := "newman@usps.com"
		to := []string{"jerry@seinfeld.com"}
		cc := []string{"cc@example.com"}
		bcc := []string{"bcc@example.com"}
		replyTo := "newman@usps.com"
		subject := "Subject"
		text := "Text body"
		html := "<p>HTML body</p>"
		attachments := []*Attachment{
			{Filename: "test.txt", Content: []byte("test content")},
		}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, cc, email.Cc)
		assert.Equal(t, bcc, email.Bcc)
		assert.Equal(t, replyTo, email.ReplyTo)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, text, email.Text)
		assert.Equal(t, html, email.HTML)
		assert.Equal(t, attachments, email.Attachments)
	})
}

func TestEmailMessageSetters(t *testing.T) {
	email := &EmailMessage{}

	t.Run("SetFrom", func(t *testing.T) {
		expected := "newman@usps.com"
		email.SetFrom(expected)

		assert.Equal(t, expected, email.From)
	})

	t.Run("SetSubject", func(t *testing.T) {
		expected := "Subject"
		email.SetSubject(expected)

		assert.Equal(t, expected, email.Subject)
	})

	t.Run("SetTo", func(t *testing.T) {
		expected := []string{"jerry@seinfeld.com"}
		email.SetTo(expected)

		assert.Equal(t, expected, email.To)
	})

	t.Run("SetCC", func(t *testing.T) {
		expected := []string{"cc@example.com"}
		email.SetCC(expected)

		assert.Equal(t, expected, email.Cc)
	})

	t.Run("SetBCC", func(t *testing.T) {
		expected := []string{"bcc@example.com"}
		email.SetBCC(expected)

		assert.Equal(t, expected, email.Bcc)
	})

	t.Run("SetReplyTo", func(t *testing.T) {
		expected := "newman@usps.com"
		email.SetReplyTo(expected)

		assert.Equal(t, expected, email.ReplyTo)
	})

	t.Run("SetText", func(t *testing.T) {
		expected := "Text body"
		email.SetText(expected)

		assert.Equal(t, expected, email.Text)
	})

	t.Run("SetHTML", func(t *testing.T) {
		expected := "<p>HTML body</p>"
		email.SetHTML(expected)

		assert.Equal(t, expected, email.HTML)
	})

	t.Run("SetAttachments", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("test content")}
		email.SetAttachments([]*Attachment{&attachment})

		assert.Contains(t, email.Attachments, &attachment)
		assert.EqualValues(t, email.Attachments, []*Attachment{&attachment})
	})

	t.Run("AddAttachment", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("test content")}
		email.AddAttachment(&attachment)

		assert.Contains(t, email.Attachments, &attachment)
	})

	t.Run("AddToRecipient", func(t *testing.T) {
		recipient := "newjerry@seinfeld.com"
		email.AddToRecipient(recipient)

		assert.Contains(t, email.To, recipient)
	})

	t.Run("AddCCRecipient", func(t *testing.T) {
		recipient := "newcc@example.com"
		email.AddCCRecipient(recipient)

		assert.Contains(t, email.Cc, recipient)
	})

	t.Run("AddBCCRecipient", func(t *testing.T) {
		recipient := "newbcc@example.com"
		email.AddBCCRecipient(recipient)

		assert.Contains(t, email.Bcc, recipient)
	})
}

func TestAddsEmailMessageToNils(t *testing.T) {
	t.Run("create full email message", func(t *testing.T) {
		from := "newman@usps.com"
		to := "jerry@seinfeld.com"
		cc := "cc@example.com"
		bcc := "bcc@example.com"
		replyTo := "newman@usps.com"
		subject := "Subject"
		text := "Text body"
		html := "<p>HTML body</p>"
		attachment := Attachment{Filename: "test.txt", Content: []byte("test content")}
		email := NewFullEmailMessage(from, nil, subject, nil, nil, replyTo, text, html, nil)

		email.AddToRecipient(to)
		email.AddCCRecipient(cc)
		email.AddBCCRecipient(bcc)
		email.AddAttachment(&attachment)

		assert.Equal(t, from, email.From)
		assert.Equal(t, []string{to}, email.To)
		assert.Equal(t, []string{cc}, email.Cc)
		assert.Equal(t, []string{bcc}, email.Bcc)
		assert.Equal(t, replyTo, email.ReplyTo)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, text, email.Text)
		assert.Equal(t, html, email.HTML)
		assert.Equal(t, []*Attachment{&attachment}, email.Attachments)
	})
}

func TestIsHTMLEdgeCases(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		input := ""
		result := IsHTML(input)

		assert.False(t, result)
	})

	t.Run("string without HTML tags", func(t *testing.T) {
		input := "Just a plain text"
		result := IsHTML(input)

		assert.False(t, result)
	})

	t.Run("string with incomplete HTML tag", func(t *testing.T) {
		input := "<div>Test"
		result := IsHTML(input)

		assert.True(t, result)
	})
}

func TestNewEmailMessageEdgeCases(t *testing.T) {
	t.Run("create email with empty body", func(t *testing.T) {
		from := "newman@usps.com"
		to := []string{"jerry@seinfeld.com"}
		subject := "Subject"
		body := ""
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, body, email.Text)
		assert.Equal(t, "", email.HTML)
	})

	t.Run("create email with only spaces in body", func(t *testing.T) {
		from := "newman@usps.com"
		to := []string{"jerry@seinfeld.com"}
		subject := "Subject"
		body := "     "
		email := NewEmailMessage(from, to, subject, body)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, "     ", email.Text)
		assert.Equal(t, "", email.HTML)
	})
}

func TestNewFullEmailMessageEdgeCases(t *testing.T) {
	t.Run("create full email message with no attachments", func(t *testing.T) {
		from := "newman@usps.com"
		to := []string{"jerry@seinfeld.com"}
		cc := []string{"cc@example.com"}
		bcc := []string{"bcc@example.com"}
		replyTo := "newman@usps.com"
		subject := "Subject"
		text := "Text body"
		html := "<p>HTML body</p>"
		attachments := []*Attachment{}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, cc, email.Cc)
		assert.Equal(t, bcc, email.Bcc)
		assert.Equal(t, replyTo, email.ReplyTo)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, text, email.Text)
		assert.Equal(t, html, email.HTML)
		assert.Equal(t, attachments, email.Attachments)
	})

	t.Run("create full email message with empty fields", func(t *testing.T) {
		from := ""
		to := []string{}
		cc := []string{}
		bcc := []string{}
		replyTo := ""
		subject := ""
		text := ""
		html := ""
		attachments := []*Attachment{}
		email := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, text, html, attachments)

		assert.Equal(t, from, email.From)
		assert.Equal(t, to, email.To)
		assert.Equal(t, cc, email.Cc)
		assert.Equal(t, bcc, email.Bcc)
		assert.Equal(t, replyTo, email.ReplyTo)
		assert.Equal(t, subject, email.Subject)
		assert.Equal(t, text, email.Text)
		assert.Equal(t, html, email.HTML)
		assert.Equal(t, attachments, email.Attachments)
	})
}

func TestAddToRecipientEdgeCases(t *testing.T) {
	t.Run("Add multiple recipients", func(t *testing.T) {
		email := &EmailMessage{}
		recipients := []string{"sfunk@funkytown.com", "kwaters@quiddich.com", "recipient3@example.com"}

		for _, recipient := range recipients {
			email.AddToRecipient(recipient)
		}

		assert.Equal(t, recipients, email.To)
	})

	t.Run("Add recipient to nil EmailMessage", func(t *testing.T) {
		var email *EmailMessage

		assert.Panics(t, func() { email.AddToRecipient("jerry@seinfeld.com") })
	})
}

func TestSetCCEdgeCases(t *testing.T) {
	t.Run("SetCC with empty slice", func(t *testing.T) {
		email := &EmailMessage{}
		expected := []string{}
		email.SetCC(expected)

		assert.Equal(t, expected, email.Cc)
	})
}

func TestSetBCCEdgeCases(t *testing.T) {
	t.Run("SetBCC with empty slice", func(t *testing.T) {
		email := &EmailMessage{}
		expected := []string{}
		email.SetBCC(expected)

		assert.Equal(t, expected, email.Bcc)
	})
}

func TestSetMaxAttachmentSize(t *testing.T) {
	email := &EmailMessage{}

	t.Run("SetMaxAttachmentSize", func(t *testing.T) {
		expected := 10 * 1024 * 1024 // 10 MB
		email.SetMaxAttachmentSize(expected)

		assert.Equal(t, expected, email.maxAttachmentSize)
	})
}

func TestGetAttachmentsWithMaxSize(t *testing.T) {
	email := &EmailMessage{
		Attachments: []*Attachment{
			{Filename: "small.txt", Content: []byte("small content")},
			{Filename: "large.txt", Content: make([]byte, 30*1024*1024)}, // 30 MB
		},
		maxAttachmentSize: 25 * 1024 * 1024, // 25 MB
	}

	t.Run("GetAttachments with size limit", func(t *testing.T) {
		expected := []*Attachment{
			{Filename: "small.txt", Content: []byte("small content")},
		}
		result := email.GetAttachments()

		assert.Equal(t, expected, result)
	})

	t.Run("GetAttachments with no size limit", func(t *testing.T) {
		email.SetMaxAttachmentSize(-1)

		expected := email.Attachments
		result := email.GetAttachments()

		assert.Equal(t, expected, result)
	})
}

func TestBuildMimeMessage(t *testing.T) {
	tests := []struct {
		message  *EmailMessage
		contains []string
	}{
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps"),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Subject: Test Email", "The air is so dewy sweet you dont even have to lick the stamps"},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "<p>The air is so dewy sweet you dont even have to lick the stamps</p>"),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Subject: Test Email", "Content-Type: text/html", "<p>The air is so dewy sweet you dont even have to lick the stamps</p>"},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
				SetCC([]string{"cc@example.com"}).
				SetBCC([]string{"bcc@example.com"}).
				SetAttachments([]*Attachment{NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!"))}),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Cc: cc@example.com", "Subject: Test Email", "The air is so dewy sweet you dont even have to lick the stamps", "Content-Disposition: attachment; filename=\"test.txt\"", base64.StdEncoding.EncodeToString([]byte("When you control the mail, you control… INFORMATION!"))},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
				SetCC([]string{"cc@example.com"}).
				SetBCC([]string{"bcc@example.com"}).
				SetReplyTo("reply-to@example.com"),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Cc: cc@example.com", "Subject: Test Email", "The air is so dewy sweet you dont even have to lick the stamps", "Reply-To: reply-to@example.com"},
		},
	}

	for _, test := range tests {
		t.Run(test.message.GetSubject(), func(t *testing.T) {
			result, err := BuildMimeMessage(test.message)
			require.NoError(t, err)

			for _, substring := range test.contains {
				assert.Contains(t, string(result), substring)
			}
		})
	}
}

func TestEmailMessageDefaultScrubbersEdgeCases(t *testing.T) {
	email := &EmailMessage{}

	t.Run("SetSubject with default scrubber", func(t *testing.T) {
		subjectInjected := `<Subject> & "attack"`
		expected := `&lt;Subject&gt; &amp; &#34;attack&#34;`

		email.SetSubject(subjectInjected)

		assert.Equal(t, subjectInjected, email.Subject)

		result := email.GetSubject()
		assert.Equal(t, expected, result)
	})

	t.Run("SetText with default scrubber", func(t *testing.T) {
		testInjected := `Hello <world> & "everyone"`

		expected := `Hello &lt;world&gt; &amp; &#34;everyone&#34;`

		email.SetText(testInjected)
		assert.Equal(t, testInjected, email.Text)

		result := email.GetText()
		assert.Equal(t, expected, result)
	})

	t.Run("SetHTML with default scrubber", func(t *testing.T) {
		htmlInjected := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`

		expected := `<div>XSS</div>`

		email.SetHTML(htmlInjected)
		assert.Equal(t, htmlInjected, email.HTML)

		result := email.GetHTML()
		assert.Equal(t, expected, result)
	})
}

func TestEmailMessageCustomScrubbers(t *testing.T) {
	message := &EmailMessage{
		Subject: `<Subject> & "attack"`,
		Text:    `Hello <world> & "everyone"`,
		HTML:    `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`,
	}

	customScrubber := scrubber.NonScrubber()

	t.Run("SetCustomTextScrubber", func(t *testing.T) {
		message.SetCustomTextScrubber(customScrubber)

		expected := `<Subject> & "attack"`
		result := message.GetSubject()
		assert.Equal(t, expected, result)

		expected = `Hello <world> & "everyone"`

		result = message.GetText()
		assert.Equal(t, expected, result)
	})

	t.Run("SetCustomHTMLScrubber", func(t *testing.T) {
		message.SetCustomHTMLScrubber(customScrubber)

		expected := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		result := message.GetHTML()
		assert.Equal(t, expected, result)
	})
}

func TestEmailMessageSettersAndScrubbersEdgeCases(t *testing.T) {
	email := &EmailMessage{}

	t.Run("SetSubject with custom scrubber", func(t *testing.T) {
		customScrubber := scrubber.NonScrubber()
		email.SetCustomTextScrubber(customScrubber)

		expected := `<Subject> & "attack"`
		email.SetSubject(expected)
		assert.Equal(t, expected, email.Subject)

		result := email.GetSubject()
		assert.Equal(t, expected, result)
	})

	t.Run("SetText with custom scrubber", func(t *testing.T) {
		customScrubber := scrubber.NonScrubber()

		email.SetCustomTextScrubber(customScrubber)

		expected := `Hello <world> & "everyone"`
		email.SetText(expected)
		assert.Equal(t, expected, email.Text)

		result := email.GetText()
		assert.Equal(t, expected, result)
	})

	t.Run("SetHTML with custom scrubber", func(t *testing.T) {
		customScrubber := scrubber.NonScrubber()

		email.SetCustomHTMLScrubber(customScrubber)

		expected := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		email.SetHTML(expected)

		assert.Equal(t, expected, email.HTML)

		result := email.GetHTML()
		assert.Equal(t, expected, result)
	})
}

func TestGetAttachmentsWithEdgeCases(t *testing.T) {
	t.Run("GetAttachments with mixed size attachments", func(t *testing.T) {
		email := &EmailMessage{
			Attachments: []*Attachment{
				{Filename: "small.txt", Content: []byte("small content")},
				{Filename: "large.txt", Content: make([]byte, 30*1024*1024)}, // 30 MB
			},
			maxAttachmentSize: 25 * 1024 * 1024, // 25 MB
		}

		expected := []*Attachment{
			{Filename: "small.txt", Content: []byte("small content")},
		}
		result := email.GetAttachments()

		assert.Equal(t, expected, result)
	})

	t.Run("GetAttachments with no size limit", func(t *testing.T) {
		email := &EmailMessage{
			Attachments: []*Attachment{
				{Filename: "small.txt", Content: []byte("small content")},
				{Filename: "large.txt", Content: make([]byte, 30*1024*1024)}, // 30 MB
			},
			maxAttachmentSize: -1, // No size limit
		}

		result := email.GetAttachments()

		assert.Equal(t, email.Attachments, result)
	})
}

func TestBuildMimeMessageWithScrubbers(t *testing.T) {
	tests := []struct {
		message  *EmailMessage
		contains []string
	}{
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps"),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Subject: Test Email", "The air is so dewy sweet you dont even have to lick the stamps"},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "<p>The air is so dewy sweet you dont even have to lick the stamps</p>"),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Subject: Test Email", "Content-Type: text/html", "<p>The air is so dewy sweet you dont even have to lick the stamps</p>"},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
				SetCC([]string{"cc@example.com"}).
				SetBCC([]string{"bcc@example.com"}).
				SetAttachments([]*Attachment{NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!"))}),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Cc: cc@example.com", "Subject: Test Email", "The air is so dewy sweet you dont even have to lick the stamps", "Content-Disposition: attachment; filename=\"test.txt\"", base64.StdEncoding.EncodeToString([]byte("When you control the mail, you control… INFORMATION!"))},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
				SetCC([]string{"cc@example.com"}).
				SetBCC([]string{"bcc@example.com"}).
				SetReplyTo("reply-to@example.com"),
			[]string{"From: newman@usps.com", "To: jerry@seinfeld.com", "Cc: cc@example.com", "Subject: Test Email", "The air is so dewy sweet you dont even have to lick the stamps", "Reply-To: reply-to@example.com"},
		},
		{
			NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "<script>alert('xss')</script>", "The air is so dewy sweet you dont even have to lick the stamps"),
			[]string{"Subject: &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		},
	}

	for _, test := range tests {
		t.Run(test.message.GetSubject(), func(t *testing.T) {
			result, err := BuildMimeMessage(test.message)
			require.NoError(t, err)

			for _, substring := range test.contains {
				assert.Contains(t, string(result), substring)
			}
		})
	}
}
