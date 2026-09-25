package gmail

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/theopenlane/newman"
)

// MockTokenManager is a mock implementation of the GmailTokenManager interface
type MockTokenManager struct {
	token []byte
	err   error
}

// TestEmailSenderImplementation checks if gmailEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*gmailEmailSender)(nil)
}

func (m *MockTokenManager) GetToken() ([]byte, error) {
	return m.token, m.err
}

// GmailMockRoundTripper implements the http.RoundTripper interface
type GmailMockRoundTripper struct {
	Response *http.Response
	Err      error
}

func readRequestBody(req *http.Request) ([]byte, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	defer req.Body.Close()

	return body, nil
}

func (m *GmailMockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	reqBody, err := readRequestBody(req)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(reqBody)),
		Header:     make(http.Header),
	}, nil
}

// MockGmailService is a mock implementation of the Gmail UsersMessagesService
type MockGmailService struct{}

func (m *MockGmailService) Send(userID string, message *gmail.Message) *gmail.UsersMessagesSendCall {
	srv, _ := gmail.NewService(context.Background())
	return gmail.NewUsersMessagesService(srv).Send(userID, message)
}

// buildMockGmailMessageSenderWrapper builds a mock gmailMessageSenderWrapper for testing
func buildMockGmailMessageSenderWrapper(err error) *gmailEmailSender {
	mockGmailService, _ := gmail.NewService(context.Background(), option.WithHTTPClient(&http.Client{
		Transport: &GmailMockRoundTripper{
			Err: err,
		},
	}))

	return &gmailEmailSender{
		users: mockGmailService.Users,
		user:  "me",
	}
}

type mockGmailTokenManager struct{}

func (m *mockGmailTokenManager) GetToken() ([]byte, error) {
	token := &oauth2.Token{
		AccessToken: "mockAccessToken",
		TokenType:   "Bearer",
	}

	return json.Marshal(token)
}

// Mock implementations for new edge cases
type mockInvalidTokenManager struct{}

func (m *mockInvalidTokenManager) GetToken() ([]byte, error) {
	return nil, ErrInvalidToken
}

func TestSendEmailWithMockService(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(nil)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailWithSendMessageError(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(ErrFailedToSendEmail)
	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailWithMockServiceError(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(ErrMockServiceError)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailWithBCC(t *testing.T) {
	emailSender := buildMockGmailMessageSenderWrapper(nil)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps").
		SetBCC([]string{"bcc@example.com"})

	err := emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestNewGmailEmailSenderOauth2(t *testing.T) {
	configJSON := []byte(`{
		"installed": {
			"client_id": "your-client-id",
			"project_id": "your-project-id",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "your-client-secret",
			"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
		}
	}`)

	tokenManager := &mockGmailTokenManager{}
	user := "me"

	emailSender, err := NewWithOauth2(context.Background(), configJSON, tokenManager, user)
	assert.NoError(t, err)

	gmailSender, ok := emailSender.(*gmailEmailSender)
	require.True(t, ok)

	assert.Equal(t, user, gmailSender.user)
	assert.NotNil(t, gmailSender.users)
}

func TestGmailEmailSenderOauth2GetGmailMessageSender(t *testing.T) {
	configJSON := []byte(`{
		"installed": {
			"client_id": "your-client-id",
			"project_id": "your-project-id",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "your-client-secret",
			"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
		}
	}`)

	tokenManager := &mockGmailTokenManager{}
	user := "me"

	emailSender, err := NewWithOauth2(context.Background(), configJSON, tokenManager, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderOauth2InvalidToken(t *testing.T) {
	configJSON := []byte(`{
		"installed": {
			"client_id": "your-client-id",
			"project_id": "your-project-id",
			"auth_uri": "https://accounts.google.com/o/oauth2/auth",
			"token_uri": "https://oauth2.googleapis.com/token",
			"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
			"client_secret": "your-client-secret",
			"redirect_uris": ["urn:ietf:wg:oauth:2.0:oob","http://localhost"]
		}
	}`)

	tokenManager := &mockInvalidTokenManager{}
	user := "me"

	_, err := NewWithOauth2(context.Background(), configJSON, tokenManager, user)
	assert.Error(t, err)
}

func TestNewGmailEmailSenderServiceAccount(t *testing.T) {
	jsonCredentials := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewWithServiceAccount(context.Background(), jsonCredentials, user)
	assert.NoError(t, err)

	gmailSender, ok := emailSender.(*gmailEmailSender)
	require.True(t, ok)

	assert.Equal(t, user, gmailSender.user)
	assert.NotNil(t, gmailSender.users)
}

func TestGmailEmailSenderServiceAccountGetGmailMessageSender(t *testing.T) {
	jsonCredentials := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewWithServiceAccount(context.Background(), jsonCredentials, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderAPIKey(t *testing.T) {
	apiKey := "mock_api_key" // nolint: gosec
	user := "me"

	emailSender, err := NewWithAPIKey(context.Background(), apiKey, user)
	assert.NoError(t, err)

	gmailSender, ok := emailSender.(*gmailEmailSender)
	require.True(t, ok)

	assert.Equal(t, user, gmailSender.user)
	assert.NotNil(t, gmailSender.users)
}

func TestGmailEmailSenderAPIKeyGetGmailMessageSender(t *testing.T) {
	apiKey := "mock_api_key" // nolint: gosec
	user := "me"

	emailSender, err := NewWithAPIKey(context.Background(), apiKey, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderJWT(t *testing.T) {
	configJSON := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewWithJWTConfig(context.Background(), configJSON, user)
	assert.NoError(t, err)

	gmailSender, ok := emailSender.(*gmailEmailSender)
	require.True(t, ok)

	assert.Equal(t, user, gmailSender.user)
	assert.NotNil(t, gmailSender.users)
}

func TestGmailEmailSenderJWTGetGmailMessageSender(t *testing.T) {
	configJSON := []byte(`{
		"type": "service_account",
		"project_id": "mock_project_id",
		"private_key_id": "mock_private_key_id",
		"private_key": "-----BEGIN PRIVATE KEY-----\nmock_private_key\n-----END PRIVATE KEY-----\n",
		"client_email": "mock_client_email@mock_project_id.iam.gserviceaccount.com",
		"client_id": "mock_client_id",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/mock_client_email%40mock_project_id.iam.gserviceaccount.com"
	}`)
	user := "me"

	emailSender, err := NewWithJWTConfig(context.Background(), configJSON, user)
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestNewGmailEmailSenderJWTInvalidJson(t *testing.T) {
	configJSON := []byte(`{invalid_json}`)
	user := "me"

	_, err := NewWithJWTConfig(context.Background(), configJSON, user)
	assert.Error(t, err)
}

func TestNewGmailEmailSenderJWTAccessInvalidJson(t *testing.T) {
	jsonCredentials := []byte(`{invalid_json}`)
	user := "me"

	_, err := NewWithJWTAccess(context.Background(), jsonCredentials, user)
	assert.Error(t, err)
}

func TestSendEmailWithNilGmailService(t *testing.T) {
	emailSender := &gmailEmailSender{}

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err := emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}

// buildHTTPTestGmailSender builds a gmailEmailSender whose service points at an httptest server
func buildHTTPTestGmailSender(t *testing.T, handler http.HandlerFunc) *gmailEmailSender {
	t.Helper()

	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(ts.Client()), option.WithEndpoint(ts.URL+"/"))
	require.NoError(t, err)

	return &gmailEmailSender{users: srv.Users, user: "me"}
}

func TestVerify(t *testing.T) {
	emailSender := buildHTTPTestGmailSender(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/gmail/v1/users/me/profile", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(`{"emailAddress": "newman@usps.com", "messagesTotal": 1}`))
		assert.NoError(t, err)
	})

	require.NoError(t, emailSender.Verify(context.Background()))
}

func TestVerifyUnauthorized(t *testing.T) {
	emailSender := buildHTTPTestGmailSender(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		_, err := w.Write([]byte(`{"error": {"code": 401, "message": "Request had invalid authentication credentials."}}`))
		assert.NoError(t, err)
	})

	err := emailSender.Verify(context.Background())
	require.ErrorIs(t, err, ErrVerifyFailed)
	assert.Contains(t, err.Error(), "invalid authentication credentials")
}

func TestVerifyInsufficientScope(t *testing.T) {
	emailSender := buildHTTPTestGmailSender(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/gmail/v1/users/newman@usps.com/profile", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)

		_, err := w.Write([]byte(`{"error": {"code": 403, "message": "Request had insufficient authentication scopes.", "errors": [{"message": "Insufficient Permission", "domain": "global", "reason": "insufficientPermissions"}], "status": "PERMISSION_DENIED"}}`))
		assert.NoError(t, err)
	})
	emailSender.user = "newman@usps.com"

	require.NoError(t, emailSender.Verify(context.Background()))
}

func TestNewGmailEmailSenderServiceAccountInvalidJson(t *testing.T) {
	jsonCredentials := []byte(`{invalid_json}`)
	user := "me"

	_, err := NewWithServiceAccount(context.Background(), jsonCredentials, user)
	assert.Error(t, err)
}
