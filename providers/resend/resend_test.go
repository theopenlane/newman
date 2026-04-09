package resend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/resend/resend-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/newman"
)

// TestEmailSenderImplementation checks if resendEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*resendEmailSender)(nil)
}

func testServerSuccess(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`{"id": "sent"}`))
		require.NoError(t, err)
	}))
}

func testServerError(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)

		_, err := w.Write([]byte(`{"id": ""}`))
		require.NoError(t, err)
	}))
}

func mockClient(t *testing.T, apiKey string, success bool) (*resend.Client, *httptest.Server) {
	var ts *httptest.Server
	if success {
		ts = testServerSuccess(t)
	} else {
		ts = testServerError(t)
	}

	mockClient := resend.NewClient(apiKey)
	baseURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	mockClient.BaseURL = baseURL

	return mockClient, ts
}

func TestNew(t *testing.T) {
	// Test missing API key
	_, err := New("")
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrMissingAPIKey)

	// Test "valid" API key
	apiKey := "re_send_api_key" // #nosec G101

	emailSender, err := New(apiKey)
	require.NoError(t, err)

	resendSender, ok := emailSender.(*resendEmailSender)
	require.True(t, ok)

	require.NotNil(t, emailSender)
	require.NotNil(t, resendSender.client)
	assert.Equal(t, apiKey, resendSender.client.ApiKey)

	// Test WithBaseURL option
	baseURL, err := url.Parse("https://example.com")
	require.NoError(t, err)

	emailSender, err = New(apiKey, WithBaseURL(*baseURL))
	require.NoError(t, err)

	resendSender, ok = emailSender.(*resendEmailSender)
	require.True(t, ok)

	require.NotNil(t, emailSender)
	require.NotNil(t, resendSender.client)
	assert.Equal(t, baseURL, resendSender.client.BaseURL)

	// Test WithUserAgent option
	userAgent := "newman/1.0.0"

	emailSender, err = New(apiKey, WithUserAgent(userAgent))
	require.NoError(t, err)

	resendSender, ok = emailSender.(*resendEmailSender)
	require.True(t, ok)

	require.NotNil(t, emailSender)
	require.NotNil(t, resendSender.client)
	assert.Equal(t, userAgent, resendSender.client.UserAgent)

	// Test WithAPIKey option
	newAPIKey := "new_api_key"

	emailSender, err = New("", WithAPIKey(newAPIKey))
	require.NoError(t, err)

	resendSender, ok = emailSender.(*resendEmailSender)
	require.True(t, ok)

	require.NotNil(t, emailSender)
	require.NotNil(t, resendSender.client)
	assert.Equal(t, newAPIKey, resendSender.client.ApiKey)
}

func TestSendEmail(t *testing.T) {
	apiKey := "re_send_api_key" // #nosec G101

	mc, ts := mockClient(t, apiKey, true)
	defer ts.Close()

	emailSender, err := New(apiKey, WithClient(mc))
	assert.NoError(t, err)

	message := newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry@seinfed.com"}),
		newman.WithSubject("Test Email"),
		newman.WithText("The air is so dewy sweet you dont even have to lick the stamps"),
		newman.WithCc([]string{"cc@example.com"}),
		newman.WithBcc([]string{"bcc@example.com"}),
		newman.WithReplyTo("replyto@example.com"),
		newman.WithHTML("<p>The air is so dewy sweet you dont even have to lick the stamps</p>"),
		newman.WithAttachment(newman.NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!"))),
	)

	err = emailSender.SendEmailWithContext(context.Background(), message)
	assert.NoError(t, err)

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailError(t *testing.T) {
	apiKey := "re_send_api_key" // #nosec G101

	mc, ts := mockClient(t, apiKey, false)
	defer ts.Close()

	emailSender, err := New(apiKey, WithClient(mc))
	assert.NoError(t, err)

	message := newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry@seinfed.com"}),
		newman.WithSubject("Test Email"),
		newman.WithText("The air is so dewy sweet you dont even have to lick the stamps"),
		newman.WithCc([]string{"cc@example.com"}),
		newman.WithBcc([]string{"bcc@example.com"}),
		newman.WithReplyTo("replyto@example.com"),
		newman.WithHTML("<p>The air is so dewy sweet you dont even have to lick the stamps</p>"),
		newman.WithAttachment(newman.NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!"))),
	)

	err = emailSender.SendEmailWithContext(context.Background(), message)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrFailedToSendEmail)
}

func TestSendEmailValidateFail(t *testing.T) {
	apiKey := "re_send_api_key" // #nosec G101

	mc, ts := mockClient(t, apiKey, false)
	defer ts.Close()

	emailSender, err := New(apiKey, WithClient(mc))
	assert.NoError(t, err)

	message := newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry"}), // invalid email
		newman.WithSubject("Test Email"),
		newman.WithText("The air is so dewy sweet you dont even have to lick the stamps"),
		newman.WithCc([]string{"cc@example.com"}),
		newman.WithBcc([]string{"bcc@example.com"}),
		newman.WithReplyTo("replyto@example.com"),
		newman.WithHTML("<p>The air is so dewy sweet you dont even have to lick the stamps</p>"),
		newman.WithAttachment(newman.NewAttachment("test.txt", []byte("When you control the mail, you control… INFORMATION!"))),
	)

	err = emailSender.SendEmailWithContext(context.Background(), message)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "to is required")
}

func testServerTooManyRequests(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)

		_, err := w.Write([]byte(`{"message": "too many requests - rate limit exceeded"}`))
		require.NoError(t, err)
	}))
}

// TestSendEmailWithContext_RequestPayload verifies that SendEmailWithContext builds the outbound
// resend.SendEmailRequest using the getter methods on EmailMessage rather than raw field access.
// The critical difference: getters validate and filter addresses, so invalid entries are dropped
// rather than forwarded verbatim to the provider.
func TestSendEmailWithContext_RequestPayload(t *testing.T) {
	apiKey := "re_send_api_key" // #nosec G101

	var captured resend.SendEmailRequest

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&captured))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "sent"}`))
	}))
	defer ts.Close()

	mc := resend.NewClient(apiKey)
	baseURL, err := url.Parse(ts.URL)
	require.NoError(t, err)
	mc.BaseURL = baseURL

	sender := &resendEmailSender{client: mc}

	msg := newman.NewEmailMessageWithOptions(
		newman.WithFrom("sender@example.com"),
		newman.WithTo([]string{"to@example.com"}),
		newman.WithSubject("Hello"),
		newman.WithHTML("<p>Body</p>"),
		newman.WithText("Body"),
		newman.WithCc([]string{"cc@example.com"}),
		newman.WithBcc([]string{"bcc@example.com"}),
		newman.WithReplyTo("reply@example.com"),
	)

	require.NoError(t, sender.SendEmailWithContext(context.Background(), msg))

	assert.Equal(t, "sender@example.com", captured.From)
	assert.Equal(t, []string{"to@example.com"}, captured.To)
	assert.Equal(t, []string{"cc@example.com"}, captured.Cc)
	assert.Equal(t, []string{"bcc@example.com"}, captured.Bcc)
	assert.Equal(t, "reply@example.com", captured.ReplyTo)
	assert.Equal(t, "Hello", captured.Subject)
	assert.Contains(t, captured.Html, "<p>Body</p>")
	assert.Equal(t, "Body", captured.Text)
}

// TestSendEmailWithContext_InvalidAddressesFiltered verifies that invalid addresses in CC and BCC
// are dropped by the getter validation rather than forwarded to the provider.
func TestSendEmailWithContext_InvalidAddressesFiltered(t *testing.T) {
	apiKey := "re_send_api_key" // #nosec G101

	var captured resend.SendEmailRequest

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&captured))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": "sent"}`))
	}))
	defer ts.Close()

	mc := resend.NewClient(apiKey)
	baseURL, err := url.Parse(ts.URL)
	require.NoError(t, err)
	mc.BaseURL = baseURL

	sender := &resendEmailSender{client: mc}

	msg := newman.NewEmailMessageWithOptions(
		newman.WithFrom("sender@example.com"),
		newman.WithTo([]string{"to@example.com"}),
		newman.WithSubject("Hello"),
		newman.WithCc([]string{"valid@example.com", "not-an-email"}),
		newman.WithBcc([]string{"not-an-email"}),
	)

	require.NoError(t, sender.SendEmailWithContext(context.Background(), msg))

	assert.Equal(t, []string{"valid@example.com"}, captured.Cc)
	assert.Empty(t, captured.Bcc)
}

func TestSendEmailRetryable(t *testing.T) {
	apiKey := "re_send_api_key" // #nosec G101

	ts := testServerTooManyRequests(t)
	defer ts.Close()

	mockClient := resend.NewClient(apiKey)
	baseURL, err := url.Parse(ts.URL)
	require.NoError(t, err)

	mockClient.BaseURL = baseURL

	emailSender, err := New(apiKey, WithClient(mockClient))
	assert.NoError(t, err)

	message := newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry@seinfeld.com"}),
		newman.WithSubject("Test Email"),
		newman.WithText("The air is so dewy sweet you dont even have to lick the stamps"),
	)

	err = emailSender.SendEmailWithContext(context.Background(), message)
	assert.Error(t, err)
	assert.True(t, newman.IsRetryableError(err))
}
