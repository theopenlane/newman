package mailgun

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/newman"
)

const (
	testDomain = "usps.com"
	testAPIKey = "test-api-key" // #nosec G101
)

// TestEmailSenderImplementation checks if mailgunEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*mailgunEmailSender)(nil)
}

func newTestSender(t *testing.T, handler http.HandlerFunc) *mailgunEmailSender {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	emailSender, err := New(testDomain, testAPIKey)
	require.NoError(t, err)

	mg, ok := emailSender.(*mailgunEmailSender)
	require.True(t, ok)

	mg.client.SetAPIBase(ts.URL)

	return mg
}

func TestVerify(t *testing.T) {
	mg := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/domains/"+testDomain, r.URL.Path)

		user, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "api", user)
		assert.Equal(t, testAPIKey, password)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`{"domain": {"name": "usps.com", "state": "active"}}`))
		assert.NoError(t, err)
	})

	require.NoError(t, mg.Verify(context.Background()))
}

func TestVerifyUnauthorized(t *testing.T) {
	mg := newTestSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)

		_, err := w.Write([]byte("Forbidden"))
		assert.NoError(t, err)
	})

	err := mg.Verify(context.Background())
	require.ErrorIs(t, err, ErrVerifyFailed)
	assert.Contains(t, err.Error(), "Forbidden")
}

func TestVerifyDomainMissing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Error("server should not be called without a domain")
	}))
	t.Cleanup(ts.Close)

	emailSender, err := New("", testAPIKey)
	require.NoError(t, err)

	mg, ok := emailSender.(*mailgunEmailSender)
	require.True(t, ok)

	mg.client.SetAPIBase(ts.URL)

	require.ErrorIs(t, mg.Verify(context.Background()), ErrDomainMissing)
}
