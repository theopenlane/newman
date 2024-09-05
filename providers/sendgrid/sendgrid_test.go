package sendgrid

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/newman"
)

// TestEmailSenderImplementation checks if sendGridEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*sendGridEmailSender)(nil)
}

// Mocking the SendGrid API response
func mockSendGridServer(t *testing.T, statusCode int, responseBody string) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(statusCode)

		_, err := w.Write([]byte(responseBody))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	return httptest.NewServer(handler)
}

func NewMockSendGridEmailSender(apiKey, url string) *sendGridEmailSender {
	requestHeaders := map[string]string{
		"Authorization": "Bearer " + apiKey,
		"User-Agent":    "sendgrid/" + "3.14.0" + ";go",
		"Accept":        "application/json",
	}

	request := rest.Request{
		Method:  "POST",
		BaseURL: url,
		Headers: requestHeaders,
	}

	return &sendGridEmailSender{
		client: &sendgrid.Client{Request: request},
	}
}

func TestNewSendGridEmailSender(t *testing.T) {
	apiKey := "test-api-key"
	emailSender, err := New(apiKey)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestSendGridEmailSender_SendEmail(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusOK, `{"message": "success"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
		SetCC([]string{"cc@example.com"}).
		SetBCC([]string{"bcc@example.com"}).
		SetReplyTo("replyto@example.com").
		SetHTML("<p>The air is so dewy sweet you dont even have to lick the stamps</p>").
		SetBCC([]string{"bcc@example.com"}).
		AddAttachment(newman.NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!")))

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendGridEmailSender_SendEmailWithError(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusInternalServerError, `{"message": "error"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func TestSendGridEmailSender_SendEmailWithNon200StatusCode(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusBadRequest, `{"message": "Bad Request"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func TestSendGridEmailSender_SendEmailWithEmptyFields(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusOK, `{"message": "success"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{},
		"",
		"",
	)

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendGridEmailSender_SendEmailWithAttachments(t *testing.T) {
	ts := mockSendGridServer(t, http.StatusOK, `{"message": "success"}`)
	defer ts.Close()

	emailSender := NewMockSendGridEmailSender("test-api-key", ts.URL)

	attachmentContent := "When you control the mail, you control… INFORMATION!"
	attachmentContentBase64 := base64.StdEncoding.EncodeToString([]byte(attachmentContent))

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"jerry@seinfeld.com"},
		"Look sister, go get yourself a cup of coffee or something",
		"Test Body",
	).AddAttachment(newman.NewAttachment("test.txt", []byte(attachmentContent)))

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)

	// Verify the attachment content
	v3Mail := mail.NewV3Mail()
	attachment := mail.NewAttachment()
	attachment.SetContent(attachmentContentBase64)
	attachment.SetType("text/plain")
	attachment.SetFilename("test.txt")
	attachment.SetDisposition("attachment")
	v3Mail.AddAttachment(attachment)

	assert.Equal(t, v3Mail.Attachments[0].Content, message.GetAttachments()[0].GetBase64StringContent())
}
