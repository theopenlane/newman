package newman

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEmailSender struct{}

func (m *MockEmailSender) SendEmail(message *EmailMessage) error {
	return nil
}

func TestNewEmailMessage(t *testing.T) {
	from := "newman@usps.com"
	to := []string{"jerry@seinfeld.com"}
	subject := "Look sister, go get yourself a cup of coffee or something"
	body := "And three times a week I shall require a cannoli"

	emailMessage := NewEmailMessage(from, to, subject, body)

	assert.Equal(t, from, emailMessage.GetFrom())
	assert.Equal(t, to, emailMessage.GetTo())
	assert.Equal(t, subject, emailMessage.GetSubject())
	assert.Equal(t, body, emailMessage.GetText())
	assert.Empty(t, emailMessage.GetHTML())
}

func TestNewFullEmailMessage(t *testing.T) {
	from := "newman@usps.com"
	to := []string{"jerry@seinfeld.com"}
	cc := []string{"cc@example.com"}
	bcc := []string{"bcc@example.com"}
	replyTo := "replyto@example.com"
	subject := "Look sister, go get yourself a cup of coffee or something"
	textBody := "Ill tell you a little secret about ZIP codes - they are meaningless"
	htmlBody := "<p>This is a test HTML body.</p>"
	attachments := []*Attachment{
		NewAttachment("test.txt", []byte("test content")),
	}

	emailMessage := NewFullEmailMessage(from, to, subject, cc, bcc, replyTo, textBody, htmlBody, attachments)

	assert.Equal(t, from, emailMessage.GetFrom())
	assert.Equal(t, to, emailMessage.GetTo())
	assert.Equal(t, cc, emailMessage.GetCC())
	assert.Equal(t, bcc, emailMessage.GetBCC())
	assert.Equal(t, replyTo, emailMessage.GetReplyTo())
	assert.Equal(t, subject, emailMessage.GetSubject())
	assert.Equal(t, textBody, emailMessage.GetText())
	assert.Equal(t, htmlBody, emailMessage.GetHTML())
	assert.Equal(t, attachments, emailMessage.GetAttachments())
}

func TestNewAttachment(t *testing.T) {
	filename := "test.txt"
	content := []byte("test content")

	attachment := NewAttachment(filename, content)

	assert.Equal(t, filename, attachment.GetFilename())
	assert.Equal(t, content, attachment.GetRawContent())
}

func TestNewAttachmentFromFile(t *testing.T) {
	filePath := "testdata/testfile.txt"
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	attachment, err := NewAttachmentFromFile(filePath)

	assert.NoError(t, err)
	assert.Equal(t, "testfile.txt", attachment.GetFilename())
	assert.Equal(t, content, attachment.GetRawContent())
}

func TestBuildMimeMessage(t *testing.T) {
	from := "newman@usps.com"
	to := []string{"jerry@seinfeld.com"}
	subject := "Look sister, go get yourself a cup of coffee or something"
	body := "And three times a week I shall require a cannoli"

	emailMessage := NewEmailMessage(from, to, subject, body)
	mimeMessage, err := BuildMimeMessage(emailMessage)

	assert.NoError(t, err)
	assert.NotEmpty(t, mimeMessage)
}

func TestValidateEmail(t *testing.T) {
	validEmail := "kramer@cosmo.com"
	invalidEmail := "invalid-email"

	assert.Equal(t, validEmail, ValidateEmail(validEmail))
	assert.Empty(t, ValidateEmail(invalidEmail))
}

func TestValidateEmailSlice(t *testing.T) {
	emails := []string{"kramer@cosmo.com", "invalid-email"}
	validEmails := ValidateEmailSlice(emails)

	assert.Len(t, validEmails, 1)
	assert.Equal(t, "kramer@cosmo.com", validEmails[0])
}

func TestGetMimeType(t *testing.T) {
	assert.Contains(t, GetMimeType("test.txt"), "text/plain")
	assert.Contains(t, GetMimeType("tarFile.tar"), "application/x-tar")
}

func TestSendEmail(t *testing.T) {
	sender := &MockEmailSender{}
	message := NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Look sister, go get yourself a cup of coffee or something", "And three times a week I shall require a cannoli")

	err := sender.SendEmail(message)
	assert.NoError(t, err)
}
