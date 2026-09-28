package postmark

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theopenlane/httpsling"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/scrubber"
	"github.com/theopenlane/newman/shared"
)

const (
	testServerToken = "test-server-token" // #nosec G101
	successBody     = `{"ErrorCode":0,"Message":"OK","MessageID":"b7bc2f4a-e38e-4336-af7d-e6c392c2f817","SubmittedAt":"2026-09-25T10:00:00.0000000-04:00","To":"jerry@seinfeld.com"}`
	shortTimeout    = 50 * time.Millisecond
	serverDelay     = 500 * time.Millisecond
	testStream      = "broadcasts"
	testTrackLinks  = "HtmlAndText"
	trackOpensKey   = "TrackOpens"
	trackLinksKey   = "TrackLinks"
)

// TestEmailSenderImplementation checks if postmarkEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*postmarkEmailSender)(nil)
}

func newTestSender(t *testing.T, handler http.HandlerFunc, opts ...Option) *postmarkEmailSender {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	emailSender, err := New(testServerToken, append([]Option{WithBaseURL(ts.URL)}, opts...)...)
	require.NoError(t, err)

	pm, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	return pm
}

func respond(t *testing.T, status int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		_, err := w.Write([]byte(body))
		assert.NoError(t, err)
	}
}

func basicMessage() *newman.EmailMessage {
	return newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry@seinfeld.com"}),
		newman.WithSubject("Hello Newman"),
		newman.WithText("The air is so dewy sweet you dont even have to lick the stamps"),
	)
}

func TestNew(t *testing.T) {
	_, err := New("")
	require.ErrorIs(t, err, ErrMissingServerToken)

	emailSender, err := New(testServerToken)
	require.NoError(t, err)

	pm, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	assert.Equal(t, testServerToken, pm.serverToken)
	assert.Equal(t, requestURL, pm.baseURL)
	assert.Equal(t, defaultTimeout, pm.client.Timeout)
	assert.NotNil(t, pm.requester)
}

func TestSendEmailWireContract(t *testing.T) {
	content := []byte("When you control the mail, you control... INFORMATION!")

	pm := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, sendEndpoint, r.URL.Path)
		assert.Equal(t, testServerToken, r.Header.Get(tokenHeader))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Contains(t, r.Header.Get("Content-Type"), "application/json")

		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))

		assert.Equal(t, "newman@usps.com", body["From"])
		assert.Equal(t, "jerry@seinfeld.com,elaine@benes.com", body["To"])
		assert.Equal(t, "Hello Newman", body["Subject"])
		assert.Equal(t, "<p>Hello, Jerry</p>", body["HtmlBody"])
		assert.Equal(t, "Hello, Jerry", body["TextBody"])
		assert.Equal(t, "replyto@usps.com", body["ReplyTo"])
		assert.Equal(t, "cc@usps.com", body["Cc"])
		assert.Equal(t, "bcc@usps.com", body["Bcc"])
		assert.Equal(t, defaultMessageStream, body["MessageStream"])
		assert.NotContains(t, body, trackOpensKey)
		assert.NotContains(t, body, trackLinksKey)
		assert.Equal(t, map[string]any{"campaign": "stamps", "route": "upper-west-side"}, body["Metadata"])
		assert.Equal(t, []any{
			map[string]any{"Name": "X-Alpha", "Value": "one"},
			map[string]any{"Name": "X-Beta", "Value": "two"},
		}, body["Headers"])
		assert.Equal(t, []any{
			map[string]any{
				"Name":        "test.txt",
				"Content":     base64.StdEncoding.EncodeToString(content),
				"ContentType": "text/plain; charset=utf-8",
			},
		}, body["Attachments"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(successBody))
		assert.NoError(t, err)
	})

	message := newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry@seinfeld.com", "elaine@benes.com"}),
		newman.WithSubject("Hello Newman"),
		newman.WithHTML("<p>Hello, Jerry</p>"),
		newman.WithText("Hello, Jerry"),
		newman.WithReplyTo("replyto@usps.com"),
		newman.WithCc([]string{"cc@usps.com"}),
		newman.WithBcc([]string{"bcc@usps.com"}),
		newman.WithTags([]newman.Tag{
			{Name: "campaign", Value: "stamps"},
			{Name: "route", Value: "upper-west-side"},
		}),
		newman.WithHeader("X-Beta", "two"),
		newman.WithHeader("X-Alpha", "one"),
		newman.WithAttachment(newman.NewAttachment("test.txt", content)),
	)

	require.NoError(t, pm.SendEmail(message))
}

func TestSendEmailStreamAndTracking(t *testing.T) {
	pm := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))

		assert.Equal(t, testStream, body["MessageStream"])
		assert.True(t, body[trackOpensKey] == true)
		assert.Equal(t, testTrackLinks, body[trackLinksKey])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(successBody))
		assert.NoError(t, err)
	}, WithMessageStream(testStream), WithTrackOpens(), WithTrackLinks(testTrackLinks))

	require.NoError(t, pm.SendEmail(basicMessage()))
}

func TestSendEmailErrorCodeOnSuccessStatus(t *testing.T) {
	pm := newTestSender(t, respond(t, http.StatusOK, `{"ErrorCode":406,"Message":"You tried to send to a recipient that has been marked as inactive."}`))

	err := pm.SendEmail(basicMessage())
	require.ErrorIs(t, err, ErrFailedToSendEmail)
	assert.Contains(t, err.Error(), "marked as inactive")
}

func TestSendEmailUnprocessableEntity(t *testing.T) {
	pm := newTestSender(t, respond(t, http.StatusUnprocessableEntity, `{"ErrorCode":300,"Message":"Zero recipients specified."}`))

	err := pm.SendEmail(basicMessage())
	require.ErrorIs(t, err, ErrFailedToSendEmail)
	assert.Contains(t, err.Error(), "Zero recipients specified")
}

func TestSendEmailRateLimited(t *testing.T) {
	pm := newTestSender(t, respond(t, http.StatusTooManyRequests, `{"ErrorCode":429,"Message":"Rate limit exceeded."}`))

	err := pm.SendEmail(basicMessage())
	require.Error(t, err)
	assert.True(t, newman.IsRetryableError(err))
}

func TestSendEmailServerError(t *testing.T) {
	pm := newTestSender(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "mail relay unavailable", http.StatusInternalServerError)
	})

	err := pm.SendEmail(basicMessage())
	require.ErrorIs(t, err, ErrFailedToSendEmail)
	assert.False(t, newman.IsRetryableError(err))
	assert.Contains(t, err.Error(), "mail relay unavailable")
}

func TestSendEmailCancelledContext(t *testing.T) {
	pm := newTestSender(t, respond(t, http.StatusOK, successBody))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := pm.SendEmailWithContext(ctx, basicMessage())
	require.ErrorIs(t, err, ErrFailedToSendEmail)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSendEmailTimeout(t *testing.T) {
	pm := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(serverDelay):
		case <-r.Context().Done():
		}

		w.WriteHeader(http.StatusOK)
	})

	pm.client.Timeout = shortTimeout

	err := pm.SendEmail(basicMessage())
	require.ErrorIs(t, err, ErrFailedToSendEmail)

	var netErr net.Error
	require.ErrorAs(t, err, &netErr)
	assert.True(t, netErr.Timeout())
}

// captureBody returns a handler that decodes the request body into body and responds with successBody
func captureBody(t *testing.T, body *map[string]any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, json.NewDecoder(r.Body).Decode(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(successBody))
		assert.NoError(t, err)
	}
}

func TestSendEmailSuccessEmptyBody(t *testing.T) {
	pm := newTestSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	require.NoError(t, pm.SendEmail(basicMessage()))
}

func TestWithMessageStreamEmptyKeepsDefault(t *testing.T) {
	emailSender, err := New(testServerToken, WithMessageStream(""))
	require.NoError(t, err)

	pm, ok := emailSender.(*postmarkEmailSender)
	require.True(t, ok)

	assert.Equal(t, defaultMessageStream, pm.stream)
}

func TestSendEmailSkipsEmptyTagName(t *testing.T) {
	var body map[string]any

	pm := newTestSender(t, captureBody(t, &body))

	message := basicMessage()
	message.Tags = []newman.Tag{
		{Name: "campaign", Value: "stamps"},
		{Name: "", Value: "dropped"},
	}

	require.NoError(t, pm.SendEmail(message))
	assert.Equal(t, map[string]any{"campaign": "stamps"}, body["Metadata"])
}

func TestSendEmailOmitsEmptyCollections(t *testing.T) {
	var body map[string]any

	pm := newTestSender(t, captureBody(t, &body))

	require.NoError(t, pm.SendEmail(basicMessage()))
	assert.NotContains(t, body, "Headers")
	assert.NotContains(t, body, "Metadata")
	assert.NotContains(t, body, "Attachments")
}

func TestSendEmailUnknownAttachmentType(t *testing.T) {
	var body map[string]any

	pm := newTestSender(t, captureBody(t, &body))

	message := basicMessage()
	message.Attachments = []*newman.Attachment{newman.NewAttachment("route.newmanzz", []byte("mail never stops"))}

	require.NoError(t, pm.SendEmail(message))

	attachments, ok := body["Attachments"].([]any)
	require.True(t, ok)
	require.Len(t, attachments, 1)

	first, ok := attachments[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, httpsling.ContentTypeApplicationOctetStream, first["ContentType"])
}

func TestSendEmailValidation(t *testing.T) {
	tests := []struct {
		name    string
		message *newman.EmailMessage
	}{
		{
			name: "empty from",
			message: newman.NewEmailMessageWithOptions(
				newman.WithTo([]string{"jerry@seinfeld.com"}),
				newman.WithSubject("Hello Newman"),
			),
		},
		{
			name: "empty to",
			message: newman.NewEmailMessageWithOptions(
				newman.WithFrom("newman@usps.com"),
				newman.WithSubject("Hello Newman"),
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pm := newTestSender(t, func(http.ResponseWriter, *http.Request) {
				t.Error("server should not be called for an invalid message")
			})

			err := pm.SendEmail(tc.message)

			var missing *shared.MissingRequiredFieldError
			assert.ErrorAs(t, err, &missing)
		})
	}
}

func TestVerify(t *testing.T) {
	pm := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, serverEndpoint, r.URL.Path)
		assert.Equal(t, testServerToken, r.Header.Get(tokenHeader))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`{"ID":1,"Name":"x"}`))
		assert.NoError(t, err)
	})

	require.NoError(t, pm.Verify(context.Background()))
}

func TestVerifyUnauthorized(t *testing.T) {
	pm := newTestSender(t, respond(t, http.StatusUnauthorized, `{"ErrorCode":10,"Message":"bad token"}`))

	err := pm.Verify(context.Background())
	require.ErrorIs(t, err, ErrVerifyFailed)
	assert.Contains(t, err.Error(), "bad token")
}

func TestVerifyCancelledContext(t *testing.T) {
	pm := newTestSender(t, respond(t, http.StatusOK, `{"ID":1,"Name":"x"}`))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := pm.Verify(ctx)
	require.ErrorIs(t, err, ErrVerifyFailed)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSendEmailHTMLScrubber(t *testing.T) {
	var htmlBody string

	pm := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))

		htmlBody, _ = body["HtmlBody"].(string)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(successBody))
		assert.NoError(t, err)
	}, WithHTMLScrubber(scrubber.NewPolicyScrubber(scrubber.WithEmailDefaults())))

	message := newman.NewEmailMessageWithOptions(
		newman.WithFrom("newman@usps.com"),
		newman.WithTo([]string{"jerry@seinfeld.com"}),
		newman.WithSubject("Hello Newman"),
		newman.WithHTML(`<p>Hello, Jerry</p><script>alert("newman")</script>`),
	)

	require.NoError(t, pm.SendEmail(message))
	assert.Contains(t, htmlBody, "<p>Hello, Jerry</p>")
	assert.NotContains(t, htmlBody, "<script>")
	assert.NotContains(t, htmlBody, "alert")
}
