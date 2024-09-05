package postmark

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/newman"
)

// TestEmailSenderImplementation checks if postmarkEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*postmarkEmailSender)(nil)
}

func TestNew(t *testing.T) {
	serverToken := "test-server-token"

	emailSender, err := New(serverToken)
	require.NoError(t, err)

	postmarkSender, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	assert.NotNil(t, emailSender)
	assert.Equal(t, serverToken, postmarkSender.serverToken)
	assert.Equal(t, endpoint, postmarkSender.endpoint)
	assert.Equal(t, requestURL, postmarkSender.url)
}

func TestSendEmail(t *testing.T) {
	emailSender, err := New("test-server-token")
	assert.NoError(t, err)

	postmarkSender, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
		SetCC([]string{"cc@example.com"}).
		SetBCC([]string{"bcc@example.com"}).
		SetReplyTo("replyto@example.com").
		SetHTML("<p>The air is so dewy sweet you dont even have to lick the stamps</p>").
		SetBCC([]string{"bcc@example.com"}).
		AddAttachment(newman.NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!")))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`{"status": "sent"}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))

	defer ts.Close()

	postmarkSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailWithMarshalError(t *testing.T) {
	emailSender, err := New("test-server-token")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		string(make([]byte, 1<<20)), // Intentionally large string to cause marshal error
		[]string{"jerry@seinfeld.com"},
		"Look sister, go get yourself a cup of coffee or something",
		"Test Body",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailWithRequestCreationError(t *testing.T) {
	emailSender, err := New("test-server-token")
	assert.NoError(t, err)

	postmarkSender, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	postmarkSender.url = "not a url"

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailWithSendError(t *testing.T) {
	emailSender, err := New("test-server-token")
	assert.NoError(t, err)

	postmarkSender, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * clientTimeout)

		http.Error(w, "server error", http.StatusInternalServerError)
	}))

	defer ts.Close()

	postmarkSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func TestSendEmailWithNon200StatusCode(t *testing.T) {
	emailSender, err := New("test-server-token")
	assert.NoError(t, err)

	postmarkSender, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer ts.Close()

	postmarkSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

func TestSendEmailWithEmptyFields(t *testing.T) {
	emailSender, err := New("test-server-token")
	assert.NoError(t, err)

	postmarkSender, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{},
		"",
		"",
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`{"status": "sent"}`))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))

	defer ts.Close()

	postmarkSender.url = ts.URL

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}
