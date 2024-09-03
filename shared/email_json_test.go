package shared_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/newman/shared"
)

func TestMarshalJSONCustom(t *testing.T) {
	t.Run("Marshal EmailMessage with attachments", func(t *testing.T) {
		email := shared.NewFullEmailMessage(
			"newman@usps.com",
			[]string{"jerry@seinfeld.com"},
			"Subject",
			[]string{"cc@example.com"},
			[]string{"bcc@example.com"},
			"replyto@example.com",
			"This is the email content.",
			"<p>This is the email content.</p>",
			[]*shared.Attachment{shared.NewAttachment("attachment1.txt", []byte("file content"))},
		)
		jsonData, err := json.Marshal(email)
		assert.Nil(t, err)

		expected := `{"from":"newman@usps.com","to":["jerry@seinfeld.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"<p>This is the email content.</p>","attachments":[{"filename":"attachment1.txt","content":"ZmlsZSBjb250ZW50"}]}`

		assert.JSONEq(t, expected, string(jsonData))
	})

	t.Run("Marshal EmailMessage without attachments", func(t *testing.T) {
		email := shared.NewFullEmailMessage(
			"newman@usps.com",
			[]string{"jerry@seinfeld.com"},
			"Subject",
			[]string{"cc@example.com"},
			[]string{"bcc@example.com"},
			"replyto@example.com",
			"This is the email content.",
			"<p>This is the email content.</p>",
			nil,
		)
		jsonData, err := json.Marshal(email)
		assert.Nil(t, err)

		expected := `{"from":"newman@usps.com","to":["jerry@seinfeld.com"],"cc":["cc@example.com"],"bcc":["bcc@example.com"],"replyTo":"replyto@example.com","subject":"Subject","text":"This is the email content.","html":"<p>This is the email content.</p>"}`

		assert.JSONEq(t, expected, string(jsonData))
	})
}

func TestUnmarshalJSONCustom(t *testing.T) {
	t.Run("Unmarshal EmailMessage with attachments", func(t *testing.T) {
		jsonData := `{
			"from": "newman@usps.com",
			"to": ["jerry@seinfeld.com"],
			"cc": ["cc@example.com"],
			"bcc": ["bcc@example.com"],
			"replyTo": "replyto@example.com",
			"subject": "Subject",
			"text": "This is the email content.",
			"html": "<p>This is the email content.</p>",
			"attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}]
		}`

		var email shared.EmailMessage

		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "newman@usps.com", email.GetFrom())
		assert.Equal(t, []string{"jerry@seinfeld.com"}, email.GetTo())
		assert.Equal(t, []string{"cc@example.com"}, email.GetCC())
		assert.Equal(t, []string{"bcc@example.com"}, email.GetBCC())
		assert.Equal(t, "replyto@example.com", email.GetReplyTo())
		assert.Equal(t, "Subject", email.GetSubject())
		assert.Equal(t, "This is the email content.", email.GetText())
		assert.Equal(t, "<p>This is the email content.</p>", email.GetHTML())

		expectedAttachment := shared.NewAttachment("attachment1.txt", []byte("file content"))

		assert.Equal(t, []*shared.Attachment{expectedAttachment}, email.GetAttachments())
	})

	t.Run("Unmarshal EmailMessage without attachments", func(t *testing.T) {
		jsonData := `{
			"from": "newman@usps.com",
			"to": ["jerry@seinfeld.com"],
			"cc": ["cc@example.com"],
			"bcc": ["bcc@example.com"],
			"replyTo": "replyto@example.com",
			"subject": "Subject",
			"text": "This is the email content.",
			"html": "<p>This is the email content.</p>"
		}`

		var email shared.EmailMessage

		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "newman@usps.com", email.GetFrom())
		assert.Equal(t, []string{"jerry@seinfeld.com"}, email.GetTo())
		assert.Equal(t, []string{"cc@example.com"}, email.GetCC())
		assert.Equal(t, []string{"bcc@example.com"}, email.GetBCC())
		assert.Equal(t, "replyto@example.com", email.GetReplyTo())
		assert.Equal(t, "Subject", email.GetSubject())
		assert.Equal(t, "This is the email content.", email.GetText())
		assert.Equal(t, "<p>This is the email content.</p>", email.GetHTML())
		assert.Nil(t, email.GetAttachments())
	})
}

func TestMarshalJSONEdgeCases(t *testing.T) {
	t.Run("nil EmailMessage", func(t *testing.T) {
		var email *shared.EmailMessage
		result, err := json.Marshal(email)
		assert.Nil(t, err)
		assert.Equal(t, "null", string(result))
	})

	t.Run("nil Attachment", func(t *testing.T) {
		var attachment *shared.Attachment
		result, err := json.Marshal(attachment)
		assert.Nil(t, err)
		assert.Equal(t, "null", string(result))
	})
}

func TestUnmarshalJSONEdgeCases(t *testing.T) {
	t.Run("empty JSON EmailMessage", func(t *testing.T) {
		jsonData := `{}`

		var email shared.EmailMessage

		err := json.Unmarshal([]byte(jsonData), &email)
		assert.Nil(t, err)
		assert.Equal(t, "", email.GetFrom())
		assert.Equal(t, "", email.GetSubject())
		assert.Equal(t, "", email.GetText())
		assert.Equal(t, "", email.GetHTML())
	})

	t.Run("invalid JSON EmailMessage", func(t *testing.T) {
		invalidJSONData := `{
        "from": "newman@usps.com",
        "to": "invalid_jerry@seinfeld.com",
        "cc": ["cc@example.com"],
        "bcc": ["bcc@example.com"],
        "replyTo": "replyto@example.com",
        "subject": "Subject",
        "text": "This is the email content.",
        "html": "<p>This is the email content.</p>",
        "attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}]
    }`

		var email shared.EmailMessage
		err := json.Unmarshal([]byte(invalidJSONData), &email)
		assert.Error(t, err)
	})

	t.Run("invalid JSON Attachment", func(t *testing.T) {
		jsonData := `{"filename": 123456789, "content": "invalid_base64"}`

		var attachment shared.Attachment

		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.Error(t, err)
	})

	t.Run("empty JSON Attachment", func(t *testing.T) {
		jsonData := `{}`

		var attachment shared.Attachment

		err := json.Unmarshal([]byte(jsonData), &attachment)

		assert.Nil(t, err)
		assert.Equal(t, "", attachment.GetFilename())
		assert.Equal(t, []byte{}, attachment.GetRawContent())
	})

	t.Run("invalid base64 content Attachment", func(t *testing.T) {
		jsonData := `{"filename": "file.txt", "content": "invalid_base64"}`

		var attachment shared.Attachment

		err := json.Unmarshal([]byte(jsonData), &attachment)
		assert.NotNil(t, err)
	})
}

func ExampleEmailMessage_MarshalJSON() {
	email := shared.NewFullEmailMessage(
		"newman@usps.com",
		[]string{"jerry@seinfeld.com"},
		"Subject",
		[]string{"cc@example.com"},
		[]string{"bcc@example.com"},
		"replyto@example.com",
		"This is the email content.",
		"<p>This is the email content.</p>",
		[]*shared.Attachment{
			shared.NewAttachment("attachment1.txt", []byte("file content")),
		},
	)

	jsonData, err := json.Marshal(email)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	fmt.Println("JSON output:", string(jsonData))
}

func ExampleEmailMessage_UnmarshalJSON() {
	jsonInput := `{
	    "from": "newman@usps.com",
	    "to": ["jerry@seinfeld.com"],
	    "cc": ["cc@example.com"],
	    "bcc": ["bcc@example.com"],
	    "replyTo": "replyto@example.com",
	    "subject": "Subject",
	    "text": "This is the email content.",
	    "html": "<p>This is the email content.</p>",
	    "attachments": [{"filename": "attachment1.txt", "content": "ZmlsZSBjb250ZW50"}]
	}`

	var email shared.EmailMessage

	err := json.Unmarshal([]byte(jsonInput), &email)
	if err != nil {
		fmt.Println("Error unmarshalling from JSON:", err)
		return
	}

	jsonData, err := json.Marshal(&email)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	fmt.Println("JSON output:", string(jsonData))
}
